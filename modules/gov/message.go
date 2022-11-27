package bank

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/types"
)

func (m *Module) HandleMessage(ctx context.Context, index int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch msg := cosmosMsg.(type) {
	case *govtypesv1beta1.MsgSubmitProposal:
		return handleMsgSubmitProposal(ctx, tx, index, msg, m.client.GovQueryClient, m.cdc)

	case *govtypesv1beta1.MsgDeposit:
		return handleMsgDeposit(ctx, tx, msg, m.client.GovQueryClient)

	case *govtypesv1beta1.MsgVote:
		pvm := m.tbM.MapProposalVoteMessage(types.NewProposalVoteMessage(msg.ProposalId, msg.Voter, msg.Option,
			tx.Height))
		// TODO: TEST IT
		return m.broker.PublishProposalVoteMessage(ctx, pvm)
	}

	return nil
}

// handleMsgSubmitProposal allows to properly handle a handleMsgSubmitProposal
func handleMsgSubmitProposal(
	ctx context.Context, tx *types.Tx, index int, msg *govtypesv1beta1.MsgSubmitProposal,
	govClient govtypesv1beta1.QueryClient, cdc codec.Codec,
) error {
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
	res, err := govClient.Proposal(
		ctx,
		&govtypesv1beta1.QueryProposalRequest{ProposalId: proposalID},
	)
	if err != nil {
		return err
	}

	proposal := res.Proposal

	// Unpack the content
	var content govtypesv1beta1.Content
	err = cdc.UnpackAny(proposal.Content, &content)
	if err != nil {
		return err
	}

	proposalObj := types.NewProposal(
		proposal.ProposalId,
		proposal.ProposalRoute(),
		proposal.ProposalType(),
		msg.Proposer,
		proposal.Status.String(),
		proposal.GetContent(),
		proposal.SubmitTime,
		proposal.DepositEndTime,
		proposal.VotingStartTime,
		proposal.VotingEndTime,
	)

	// Store the deposit
	deposit := types.NewProposalDeposit(proposal.ProposalId, msg.Proposer, msg.InitialDeposit, tx.Height)

	// TODO:
	_, _ = proposalObj, deposit

	return nil
}

// handleMsgDeposit allows to properly handle a handleMsgDeposit
func handleMsgDeposit(ctx context.Context, tx *types.Tx, msg *govtypesv1beta1.MsgDeposit, govClient govtypesv1beta1.QueryClient) error {
	res, err := govClient.Deposit(
		ctx,
		&govtypesv1beta1.QueryDepositRequest{ProposalId: msg.ProposalId, Depositor: msg.Depositor},
		grpcClient.GetHeightRequestHeader(tx.Height),
	)
	if err != nil {
		return fmt.Errorf("error while getting proposal deposit: %s", err)
	}

	deposit := types.NewProposalDeposit(msg.ProposalId, msg.Depositor, res.Deposit.Amount, tx.Height)

	_ = deposit
	// TODO:
	return nil
}
