package to_broker

import "github.com/hexy-dev/spacebox/broker/model"

func (tb ToBroker) MapAnnualProvision(height int64, annualProvision float64) model.AnnualProvision {
	return model.AnnualProvision{
		Height:          height,
		AnnualProvision: annualProvision,
	}
}
