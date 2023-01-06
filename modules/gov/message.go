package bank

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hexy-dev/spacebox/broker/model"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	"github.com/hexy-dev/spacebox-crawler/types"
)

func (m *Module) HandleMessage(ctx context.Context, index int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch msg := cosmosMsg.(type) {
	case *govtypesv1beta1.MsgSubmitProposal:
		return m.handleMsgSubmitProposal(ctx, tx, index, msg)

	case *govtypesv1beta1.MsgDeposit:
		return m.handleMsgDeposit(ctx, tx, msg)

	case *govtypesv1beta1.MsgVote:
		pvm := model.NewProposalVoteMessage(
			msg.ProposalId,
			tx.Height,
			msg.Voter,
			msg.Option.String(),
		)

		// TODO: TEST IT
		return m.broker.PublishProposalVoteMessage(ctx, pvm)
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
	resPb, err := m.client.GovQueryClient.Proposal(
		ctx,
		&govtypesv1beta1.QueryProposalRequest{ProposalId: proposalID},
	)
	if err != nil {
		return err
	}

	proposal := resPb.Proposal

	// Unpack the content
	var content govtypesv1beta1.Content
	err = m.cdc.UnpackAny(proposal.Content, &content)
	if err != nil {
		return err
	}

	// publish the deposit
	// TODO: test it
	if err = m.broker.PublishProposalDeposit(ctx,
		model.NewProposalDeposit(proposal.ProposalId, tx.Height, msg.Proposer,
			m.tbM.MapCoins(types.NewCoinsFromCdk(msg.InitialDeposit)))); err != nil {

		return err
	}

	// TODO: test it
	if err = m.broker.PublishProposalDepositMessage(ctx,
		model.NewProposalDepositMessage(proposal.ProposalId, tx.Height, msg.Proposer, tx.TxHash,
			m.tbM.MapCoins(types.NewCoinsFromCdk(msg.InitialDeposit)))); err != nil {

		return err
	}

	contentBytes, err := types.GetProposalContentBytes(content, m.cdc)
	if err != nil {
		return err
	}

	// TODO: test it
	if err = m.broker.PublishProposal(ctx,
		model.NewProposal(
			proposal.ProposalId, content.GetTitle(), content.GetDescription(),
			proposal.ProposalRoute(), proposal.ProposalType(), msg.Proposer, proposal.Status.String(), contentBytes,
			proposal.SubmitTime, proposal.DepositEndTime, proposal.VotingStartTime, proposal.VotingEndTime),
	); err != nil {
		return err
	}

	return nil
}

// handleMsgDeposit handles a handleMsgDeposit.
// publishes proposalDeposit and proposalDepositMessage to the broker.
func (m *Module) handleMsgDeposit(ctx context.Context, tx *types.Tx, msg *govtypesv1beta1.MsgDeposit) error {
	res, err := m.client.GovQueryClient.Deposit(
		ctx,
		&govtypesv1beta1.QueryDepositRequest{ProposalId: msg.ProposalId, Depositor: msg.Depositor},
		grpcClient.GetHeightRequestHeader(tx.Height),
	)
	if err != nil {
		return fmt.Errorf("error while getting proposal deposit: %s", err)
	}

	// TODO: test it
	if err = m.broker.PublishProposalDeposit(ctx, model.NewProposalDeposit(
		msg.ProposalId, tx.Height, msg.Depositor,
		m.tbM.MapCoins(types.NewCoinsFromCdk(res.Deposit.Amount)))); err != nil {

		return err
	}

	// TODO: test it
	if err = m.broker.PublishProposalDepositMessage(ctx, model.NewProposalDepositMessage(
		msg.ProposalId, tx.Height, msg.Depositor, tx.TxHash,
		m.tbM.MapCoins(types.NewCoinsFromCdk(res.Deposit.Amount)))); err != nil {

		return err
	}

	return nil
}
