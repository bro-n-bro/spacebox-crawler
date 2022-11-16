package bank

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/types"
)

func (m *Module) HandleMessage(ctx context.Context, index int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch msg := cosmosMsg.(type) {
	case *govtypes.MsgSubmitProposal:
		return handleMsgSubmitProposal(ctx, tx, index, msg, m.client.GovQueryClient, m.cdc)

	case *govtypes.MsgDeposit:
		return handleMsgDeposit(ctx, tx, msg, m.client.GovQueryClient)

	case *govtypes.MsgVote:
		return handleMsgVote(tx, msg)
	}

	return nil
}

// handleMsgSubmitProposal allows to properly handle a handleMsgSubmitProposal
func handleMsgSubmitProposal(
	ctx context.Context, tx *types.Tx, index int, msg *govtypes.MsgSubmitProposal,
	govClient govtypes.QueryClient, cdc codec.Codec,
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
		&govtypes.QueryProposalRequest{ProposalId: proposalID},
	)
	if err != nil {
		return err
	}

	proposal := res.Proposal

	// Unpack the content
	var content govtypes.Content
	err = cdc.UnpackAny(proposal.Content, &content)
	if err != nil {
		return err
	}

	proposalObj := types.NewProposal(
		proposal.ProposalId,
		proposal.ProposalRoute(),
		proposal.ProposalType(),
		proposal.GetContent(),
		proposal.Status.String(),
		proposal.SubmitTime,
		proposal.DepositEndTime,
		proposal.VotingStartTime,
		proposal.VotingEndTime,
		msg.Proposer,
	)

	// Store the deposit
	deposit := types.NewDeposit(proposal.ProposalId, msg.Proposer, msg.InitialDeposit, tx.Height)

	// TODO:
	_, _ = proposalObj, deposit

	return nil
}

// handleMsgDeposit allows to properly handle a handleMsgDeposit
func handleMsgDeposit(ctx context.Context, tx *types.Tx, msg *govtypes.MsgDeposit, govClient govtypes.QueryClient) error {
	res, err := govClient.Deposit(
		ctx,
		&govtypes.QueryDepositRequest{ProposalId: msg.ProposalId, Depositor: msg.Depositor},
		grpcClient.GetHeightRequestHeader(tx.Height),
	)
	if err != nil {
		return fmt.Errorf("error while getting proposal deposit: %s", err)
	}

	deposit := types.NewDeposit(msg.ProposalId, msg.Depositor, res.Deposit.Amount, tx.Height)

	_ = deposit
	// TODO:
	return nil
}

// handleMsgVote allows to properly handle a handleMsgVote
func handleMsgVote(tx *types.Tx, msg *govtypes.MsgVote) error {
	vote := types.NewVote(msg.ProposalId, msg.Voter, msg.Option, tx.Height)
	_ = vote
	// TODO:
	return nil
}
