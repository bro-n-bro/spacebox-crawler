package model

type (
	// DepositParams contains the data of the deposit parameters of the x/gov module
	DepositParams struct {
		MinDeposit       Coins `json:"min_deposit,omitempty" yaml:"min_deposit"`
		MaxDepositPeriod int64 `json:"max_deposit_period,omitempty" yaml:"max_deposit_period"`
	}

	// VotingParams contains the voting parameters of the x/gov module
	VotingParams struct {
		VotingPeriod int64 `json:"voting_period,omitempty" yaml:"voting_period"`
	}

	// TallyParams contains the tally parameters of the x/gov module
	TallyParams struct {
		Quorum        float64 `json:"quorum,omitempty"`
		Threshold     float64 `json:"threshold,omitempty"`
		VetoThreshold float64 `json:"veto_threshold,omitempty" yaml:"veto_threshold"`
	}

	GovParams struct {
		DepositParams DepositParams `json:"deposit_params" yaml:"deposit_params"`
		VotingParams  VotingParams  `json:"voting_params" yaml:"voting_params"`
		TallyParams   TallyParams   `json:"tally_params" yaml:"tally_params"`
		Height        int64         `json:"height" ymal:"height"`
	}
)
