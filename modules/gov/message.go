package bank

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
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
		// TODO: TEST IT
		return m.broker.PublishProposalVoteMessage(ctx, model.ProposalVoteMessage{
			ProposalID:   msg.ProposalId,
			Height:       tx.Height,
			VoterAddress: msg.Voter,
			Option:       msg.Option.String(),
			TxHash:       tx.TxHash,
			MsgIndex:     int64(index),
		})
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
		return err
	}

	// Get the proposal
	respPb, err := m.client.GovQueryClient.Proposal(
		ctx,
		&govtypesv1beta1.QueryProposalRequest{ProposalId: proposalID},
	)
	if err != nil {
		return err
	}

	proposal := respPb.Proposal

	// Unpack the content
	var content govtypesv1beta1.Content
	if err = m.cdc.UnpackAny(proposal.Content, &content); err != nil {
		return err
	}

	// publish the deposit
	// TODO: test it
	if err = m.broker.PublishProposalDeposit(ctx, model.ProposalDeposit{
		ProposalID:       proposal.ProposalId,
		Height:           tx.Height,
		DepositorAddress: msg.Proposer,
		Coins:            m.tbM.MapCoins(types.NewCoinsFromCdk(msg.InitialDeposit)),
	}); err != nil {
		return err
	}

	// TODO: test it
	if err = m.broker.PublishProposalDepositMessage(ctx, model.ProposalDepositMessage{
		ProposalDeposit: model.ProposalDeposit{
			ProposalID:       proposal.ProposalId,
			Height:           tx.Height,
			DepositorAddress: msg.Proposer,
			Coins:            m.tbM.MapCoins(types.NewCoinsFromCdk(msg.InitialDeposit)),
		},
		TxHash:   tx.TxHash,
		MsgIndex: int64(index),
	}); err != nil {
		return err
	}

	contentBytes, err := utils.GetProposalContentBytes(content, m.cdc)
	if err != nil {
		return err
	}

	// TODO: test it
	if err = m.broker.PublishProposal(ctx, model.Proposal{
		ID:              proposal.ProposalId,
		Title:           content.GetTitle(),
		Description:     content.GetDescription(),
		ProposalRoute:   proposal.ProposalRoute(),
		ProposalType:    proposal.ProposalType(),
		ProposerAddress: msg.Proposer,
		Status:          proposal.Status.String(),
		Content:         contentBytes,
		SubmitTime:      proposal.SubmitTime,
		DepositEndTime:  proposal.DepositEndTime,
		VotingStartTime: proposal.VotingStartTime,
		VotingEndTime:   proposal.VotingEndTime,
	}); err != nil {
		return err
	}

	return nil
}

// handleMsgDeposit handles a handleMsgDeposit.
// publishes proposalDeposit and proposalDepositMessage to the broker.
func (m *Module) handleMsgDeposit(ctx context.Context, tx *types.Tx, index int, msg *govtypesv1beta1.MsgDeposit) error {
	res, err := m.client.GovQueryClient.Deposit(
		ctx,
		&govtypesv1beta1.QueryDepositRequest{ProposalId: msg.ProposalId, Depositor: msg.Depositor},
		grpcClient.GetHeightRequestHeader(tx.Height),
	)
	if err != nil {
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

	// TODO: test it
	if err = m.broker.PublishProposalDepositMessage(ctx, model.ProposalDepositMessage{
		ProposalDeposit: model.ProposalDeposit{
			ProposalID:       msg.ProposalId,
			DepositorAddress: msg.Depositor,
			Height:           tx.Height,
			Coins:            m.tbM.MapCoins(types.NewCoinsFromCdk(res.Deposit.Amount)),
		},
		TxHash:   tx.TxHash,
		MsgIndex: int64(index),
	}); err != nil {
		return err
	}

	return nil
}
