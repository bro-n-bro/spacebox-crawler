package bank

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	"github.com/hexy-dev/spacebox-crawler/types"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block) error {
	header := grpcClient.GetHeightRequestHeader(block.Height)

	respDeposit, err := m.client.GovQueryClient.Params(
		ctx,
		&govtypes.QueryParamsRequest{ParamsType: govtypes.ParamDeposit},
		header,
	)
	if err != nil {
		return err
	}

	respVoting, err := m.client.GovQueryClient.Params(
		ctx,
		&govtypes.QueryParamsRequest{ParamsType: govtypes.ParamVoting},
		header,
	)
	if err != nil {
		return err

	}

	respTally, err := m.client.GovQueryClient.Params(
		ctx,
		&govtypes.QueryParamsRequest{ParamsType: govtypes.ParamTallying},
		header,
	)
	if err != nil {
		return err
	}

	params := model.NewGowParams(
		block.Height,
		model.NewDepositParams(
			respDeposit.DepositParams.MaxDepositPeriod.Nanoseconds(),
			m.tbM.MapCoins(types.NewCoinsFromCdk(respDeposit.DepositParams.MinDeposit))),
		model.NewVotingParams(respVoting.VotingParams.VotingPeriod.Nanoseconds()),
		model.NewTallyParams(
			respTally.TallyParams.Quorum.MustFloat64(),
			respTally.TallyParams.Threshold.MustFloat64(),
			respTally.TallyParams.VetoThreshold.MustFloat64(),
		),
	)

	// TODO: test it
	// TODO: maybe check diff from mongo in my side?
	if err = m.broker.PublishGovParams(ctx, params); err != nil {
		return err
	}

	// TODO: UpdateProposal
	return err
}

// updateProposals updates the proposals
// TODO: how to update it?
// func (m *Module) updateProposals(ctx context.Context, height int64, blockVals *tmctypes.ResultValidators) error {
//	var ids []uint64
//	// ids, err := GetOpenProposalsIds()
//	// if err != nil {
//	//	return err
//	// }
//
//	if len(ids) > 0 {
//		clients := govutils.NewUpdateProposalClients(m.client.GovQueryClient, m.client.BankQueryClient,
//			m.client.StakingQueryClient)
//
//		for _, id := range ids {
//			err := govutils.UpdateProposal(ctx, height, blockVals, id, clients, m.cdc, m.tbM, m.broker)
//			if err != nil {
//				return err
//			}
//		}
//	}
//
//	return nil
// }
