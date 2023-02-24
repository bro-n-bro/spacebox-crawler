package bank

import (
	"context"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	jsoniter "github.com/json-iterator/go"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
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

	// TODO: UpdateProposal

	// TODO: test it
	// TODO: maybe check diff from mongo in my side?

	dp := model.DepositParams{
		MinDeposit:       m.tbM.MapCoins(types.NewCoinsFromCdk(respDeposit.DepositParams.MinDeposit)),
		MaxDepositPeriod: respDeposit.DepositParams.MaxDepositPeriod.Nanoseconds(),
	}
	vp := model.VotingParams{
		VotingPeriod: respVoting.VotingParams.VotingPeriod.Nanoseconds(),
	}
	tp := model.TallyParams{
		Quorum:        respTally.TallyParams.Quorum.MustFloat64(),
		Threshold:     respTally.TallyParams.Threshold.MustFloat64(),
		VetoThreshold: respTally.TallyParams.VetoThreshold.MustFloat64(),
	}

	dbStr, _ := jsoniter.MarshalToString(dp)
	vpStr, _ := jsoniter.MarshalToString(vp)
	tpStr, _ := jsoniter.MarshalToString(tp)

	return m.broker.PublishGovParams(ctx, model.GovParams{
		DepositParams: dbStr,
		VotingParams:  vpStr,
		TallyParams:   tpStr,
		Height:        block.Height,
	})
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
