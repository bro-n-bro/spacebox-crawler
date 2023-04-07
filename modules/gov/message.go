package bank

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	errDepositerNotFoundForProposal = `rpc error: code = %s desc = depositer: %s not found for proposal: %d`
)

func (m *Module) HandleMessage(ctx context.Context, index int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch msg := cosmosMsg.(type) {
	case *govtypesv1beta1.MsgSubmitProposal:
		return m.handleMsgSubmitProposal(ctx, tx, index, msg)

	case *govtypesv1beta1.MsgDeposit:
		return m.handleMsgDeposit(ctx, tx, index, msg)

	case *govtypesv1beta1.MsgVote:
		return m.handleMsgVote(ctx, tx, index, msg)
	case *govtypesv1beta1.MsgVoteWeighted:
		return m.handlerMsgVoteWeighted(ctx, tx, msg)
	}

	return nil
}

// handleMsgSubmitProposal handles a handleMsgSubmitProposal
// publishes proposal, proposalDeposit and proposalDepositMessage to the broker.
func (m *Module) handleMsgSubmitProposal(ctx context.Context, tx *types.Tx, index int,
	msg *govtypesv1beta1.MsgSubmitProposal) error {
	// Get the proposal id
	event, err := tx.FindEventByType(index, govtypes.EventTypeSubmitProposal)
	if err != nil {
		return err
	}

	id, err := tx.FindAttributeByKey(event, govtypes.AttributeKeyProposalID)
	if err != nil {
		return err
	}

	proposalID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		m.log.Error().Err(err).Str("handler", "handleMsgSubmitProposal").Msg("parse uint error")
		return err
	}

	if err = m.getAndPublishProposal(ctx, proposalID, msg.Proposer); err != nil {
		return err
	}

	// publish the deposit
	// TODO: test it
	if err = m.broker.PublishProposalDeposit(ctx, model.ProposalDeposit{
		ProposalID:       proposalID,
		Height:           tx.Height,
		DepositorAddress: msg.Proposer,
		Coins:            m.tbM.MapCoins(types.NewCoinsFromCdk(msg.InitialDeposit)),
	}); err != nil {
		return err
	}

	// TODO: test it
	if err = m.broker.PublishProposalDepositMessage(ctx, model.ProposalDepositMessage{
		ProposalDeposit: model.ProposalDeposit{
			ProposalID:       proposalID,
			Height:           tx.Height,
			DepositorAddress: msg.Proposer,
			Coins:            m.tbM.MapCoins(types.NewCoinsFromCdk(msg.InitialDeposit)),
		},
		TxHash:   tx.TxHash,
		MsgIndex: int64(index),
	}); err != nil {
		return err
	}

	return nil
}

// handleMsgDeposit handles a MsgDeposit message.
// Publishes proposalDeposit and proposalDepositMessage data to the broker.
func (m *Module) handleMsgDeposit(ctx context.Context, tx *types.Tx, index int, msg *govtypesv1beta1.MsgDeposit) error {
	if err := m.broker.PublishProposalDepositMessage(ctx, model.ProposalDepositMessage{
		ProposalDeposit: model.ProposalDeposit{
			ProposalID:       msg.ProposalId,
			DepositorAddress: msg.Depositor,
			Height:           tx.Height,
			Coins:            m.tbM.MapCoins(types.NewCoinsFromCdk(msg.Amount)),
		},
		TxHash:   tx.TxHash,
		MsgIndex: int64(index),
	}); err != nil {
		return err
	}

	res, err := m.client.GovQueryClient.Deposit(
		ctx,
		&govtypesv1beta1.QueryDepositRequest{ProposalId: msg.ProposalId, Depositor: msg.Depositor},
		grpcClient.GetHeightRequestHeader(tx.Height),
	)
	if err != nil {
		var code string
		if _, err = fmt.Sscanf(
			err.Error(),
			errDepositerNotFoundForProposal,
			&code, &msg.Depositor, &msg.ProposalId,
		); err != nil {
			return err
		}

		if code == codes.InvalidArgument.String() {
			return nil
		}

		return fmt.Errorf("error while getting proposal deposit: %w", err)
	}

	// TODO: test it
	if err = m.broker.PublishProposalDeposit(ctx, model.ProposalDeposit{
		ProposalID:       msg.ProposalId,
		DepositorAddress: msg.Depositor,
		Height:           tx.Height,
		Coins:            m.tbM.MapCoins(types.NewCoinsFromCdk(res.Deposit.Amount)),
	}); err != nil {
		return err
	}

	return nil
}

// handleMsgVote handles a MsgVote message.
// Publishes proposalVoteMessage and proposalTallyResult data to the broker.
func (m *Module) handleMsgVote(ctx context.Context, tx *types.Tx, index int, msg *govtypesv1beta1.MsgVote) error {
	// TODO: TEST IT
	if err := m.broker.PublishProposalVoteMessage(ctx, model.ProposalVoteMessage{
		ProposalID:   msg.ProposalId,
		Height:       tx.Height,
		VoterAddress: msg.Voter,
		Option:       msg.Option.String(),
		TxHash:       tx.TxHash,
		MsgIndex:     int64(index),
	}); err != nil {
		m.log.Error().Err(err).Msg("error while publishing proposal vote message")
		return err
	}

	return m.getAndPublishTallyResult(ctx, msg.ProposalId, tx.Height)
}

// handlerMsgVoteWeighted handles MsgVoteWeighted message.
// Gets tallyResult data from node and publishes it to the broker.
func (m *Module) handlerMsgVoteWeighted(ctx context.Context, tx *types.Tx, msg *govtypesv1beta1.MsgVoteWeighted) error {
	if m.tallyCache != nil && !m.tallyCache.UpdateCacheValue(msg.ProposalId, tx.Height) {
		return nil
	}

	respPb, err := m.client.GovQueryClient.TallyResult(
		ctx,
		&govtypesv1beta1.QueryTallyResultRequest{ProposalId: msg.ProposalId},
	)

	if err != nil {
		status, ok := status.FromError(err)
		if ok && status.Code() == codes.NotFound {
			return nil
		}
		return err
	}

	return m.broker.PublishProposalTallyResult(ctx, model.ProposalTallyResult{
		ProposalID: msg.ProposalId,
		Yes:        respPb.Tally.Yes.Int64(),
		No:         respPb.Tally.No.Int64(),
		Abstain:    respPb.Tally.Abstain.Int64(),
		NoWithVeto: respPb.Tally.NoWithVeto.Int64(),
		Height:     tx.Height,
	})
}
