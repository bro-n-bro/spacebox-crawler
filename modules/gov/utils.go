package bank

import (
	"context"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bro-n-bro/spacebox/broker/model"
)

var (
	errInvalidProposalContent = errors.New("invalid proposal content type")
)

func (m *Module) getAndPublishProposal(ctx context.Context, proposalID uint64, proposer string) error {
	proposalResp, err := m.client.GovQueryClient.Proposal(ctx, &govtypes.QueryProposalRequest{ProposalId: proposalID})
	if err != nil {
		status, ok := status.FromError(err)
		if ok && status.Code() == codes.NotFound {
			m.log.Warn().Uint64("proposal_id", proposalID).Msg("proposal not found")
			return nil
		}
		return err
	}

	proposal := proposalResp.Proposal

	// Unpack the content
	var content govtypes.Content
	if err = m.cdc.UnpackAny(proposal.Content, &content); err != nil {
		return err
	}

	contentBytes, err := getProposalContentBytes(content, m.cdc)
	if err != nil {
		return err
	}

	return m.broker.PublishProposal(ctx, model.Proposal{
		ID:              proposal.ProposalId,
		Title:           content.GetTitle(),
		Description:     content.GetDescription(),
		ProposalRoute:   proposal.ProposalRoute(),
		ProposalType:    proposal.ProposalType(),
		ProposerAddress: proposer,
		Status:          proposal.Status.String(),
		Content:         contentBytes,
		SubmitTime:      proposal.SubmitTime,
		DepositEndTime:  proposal.DepositEndTime,
		VotingStartTime: proposal.VotingStartTime,
		VotingEndTime:   proposal.VotingEndTime,
	})
}

func getProposalContentBytes(content govtypes.Content, cdc codec.Codec) ([]byte, error) {
	// Encode the content properly
	protoContent, ok := content.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("%w: %T", errInvalidProposalContent, content)
	}

	anyContent, err := codectypes.NewAnyWithValue(protoContent)
	if err != nil {
		return nil, err
	}

	return cdc.MarshalJSON(anyContent)
}

func (m *Module) getAndPublishTallyResult(ctx context.Context, proposalID uint64, height int64) error {
	// publish only newest heights
	if m.tallyCache != nil && !m.tallyCache.UpdateCacheValue(proposalID, height) {
		return nil
	}

	respPb, err := m.client.GovQueryClient.TallyResult(
		ctx,
		&govtypes.QueryTallyResultRequest{ProposalId: proposalID},
	)

	if err != nil {
		status, ok := status.FromError(err)
		if ok && status.Code() == codes.NotFound {
			m.log.Warn().Uint64("proposal_id", proposalID).Msg("tally result not found")
			return nil
		}
		m.log.Error().
			Str("handler", "HandleEndBlocker").
			Err(err).
			Msg("failed to get proposal tally result")

		return err
	}

	if err := m.broker.PublishProposalTallyResult(ctx, model.ProposalTallyResult{
		ProposalID: proposalID,
		Yes:        respPb.Tally.Yes.Int64(),
		No:         respPb.Tally.No.Int64(),
		Abstain:    respPb.Tally.Abstain.Int64(),
		NoWithVeto: respPb.Tally.NoWithVeto.Int64(),
		Height:     height,
	}); err != nil {
		m.log.Error().Err(err).Msg("error while publishing proposal tally result")
		return err
	}

	return nil
}
