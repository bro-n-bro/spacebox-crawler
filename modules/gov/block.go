package bank

import (
	govutils "bro-n-bro-osmosis/modules/gov/utils"
	"context"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/types"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block, _ *tmctypes.ResultValidators) error {
	params, err := m.getGovParams(ctx, block.Height)
	// TODO: test it
	// TODO: maybe check diff from mongo in my side?
	if err := m.broker.PublishGovParams(ctx, m.tbM.MapGovParams(params)); err != nil {
		return err
	}

	// TODO: UpdateProposal
	return err
}

func (m *Module) getGovParams(ctx context.Context, height int64) (*types.GovParams, error) {
	respDeposit, err := m.client.GovQueryClient.Params(
		ctx,
		&govtypes.QueryParamsRequest{ParamsType: govtypes.ParamDeposit},
		grpcClient.GetHeightRequestHeader(height),
	)
	if err != nil {
		return nil, err
	}

	respVoting, err := m.client.GovQueryClient.Params(
		ctx,
		&govtypes.QueryParamsRequest{ParamsType: govtypes.ParamVoting},
		grpcClient.GetHeightRequestHeader(height),
	)
	if err != nil {
		return nil, err

	}

	respTally, err := m.client.GovQueryClient.Params(
		ctx,
		&govtypes.QueryParamsRequest{ParamsType: govtypes.ParamTallying},
		grpcClient.GetHeightRequestHeader(height),
	)
	if err != nil {
		return nil, err
	}

	govParams := types.NewGovParams(
		types.NewVotingParams(respVoting.GetVotingParams()),
		types.NewDepositParam(respDeposit.GetDepositParams()),
		types.NewTallyParams(respTally.GetTallyParams()),
		height,
	)

	return govParams, nil
}

// updateProposals updates the proposals
func (m *Module) updateProposals(ctx context.Context, height int64, blockVals *tmctypes.ResultValidators) error {
	var ids []uint64
	//ids, err := db.GetOpenProposalsIds()
	//if err != nil {
	//	log.Error().Err(err).Str("module", "gov").Msg("error while getting open ids")
	//}

	if len(ids) > 0 {
		clients := govutils.NewUpdateProposalClients(m.client.GovQueryClient, m.client.BankQueryClient,
			m.client.StakingQueryClient)

		for _, id := range ids {
			err := govutils.UpdateProposal(ctx, height, blockVals, id, clients, m.cdc, m.broker, m.tbM)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
