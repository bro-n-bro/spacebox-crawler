package liquidity

import (
	"context"
	"encoding/base64"
	"strconv"

	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox-crawler/types/liquidity"
	"github.com/bro-n-bro/spacebox/broker/model"
)

//nolint:lll
var (
	base64KeyPoolID                    = base64.StdEncoding.EncodeToString([]byte(liquidity.AttributeValuePoolID))
	base64KeyBatchIndex                = base64.StdEncoding.EncodeToString([]byte(liquidity.AttributeValueBatchIndex))
	base64KeyMsgIndex                  = base64.StdEncoding.EncodeToString([]byte(liquidity.AttributeValueMsgIndex))
	base64KeySwapRequester             = base64.StdEncoding.EncodeToString([]byte(liquidity.AttributeValueSwapRequester))
	base64KeyOfferCoinDenom            = base64.StdEncoding.EncodeToString([]byte(liquidity.AttributeValueOfferCoinDenom))
	base64KeyOfferCoinAmount           = base64.StdEncoding.EncodeToString([]byte(liquidity.AttributeValueOfferCoinAmount))
	base64KeyExchangedOfferCoinAmount  = base64.StdEncoding.EncodeToString([]byte(liquidity.AttributeValueExchangedOfferCoinAmount))
	base64KeyDemandCoinDenom           = base64.StdEncoding.EncodeToString([]byte(liquidity.AttributeValueDemandCoinDenom))
	base64KeyOrderPrice                = base64.StdEncoding.EncodeToString([]byte(liquidity.AttributeValueOrderPrice))
	base64KeySwapPrice                 = base64.StdEncoding.EncodeToString([]byte(liquidity.AttributeValueSwapPrice))
	base64KeyTransactedCoinAmount      = base64.StdEncoding.EncodeToString([]byte(liquidity.AttributeValueTransactedCoinAmount))
	base64KeyRemainingOfferCoinAmount  = base64.StdEncoding.EncodeToString([]byte(liquidity.AttributeValueRemainingOfferCoinAmount))
	base64KeyExchangedDemandCoinAmount = base64.StdEncoding.EncodeToString([]byte(liquidity.AttributeValueExchangedDemandCoinAmount))
	base64KeyOfferCoinFeeAmount        = base64.StdEncoding.EncodeToString([]byte(liquidity.AttributeValueOfferCoinFeeAmount))
	base64KeyExchangedCoinFeeAmount    = base64.StdEncoding.EncodeToString([]byte(liquidity.AttributeValueExchangedCoinFeeAmount))
	base64KeyOrderExpiryHeight         = base64.StdEncoding.EncodeToString([]byte(liquidity.AttributeValueOrderExpiryHeight))
	base64KeySuccess                   = base64.StdEncoding.EncodeToString([]byte(liquidity.AttributeValueSuccess))
)

//nolint:gocyclo
func (m *Module) HandleEndBlocker(ctx context.Context, eventsMap types.BlockerEvents, height int64) error {
	events, ok := eventsMap[liquidity.EventTypeSwapTransacted]
	if !ok {
		return nil
	}

	var err error
	for _, event := range events {
		if len(event.Attributes) < 13 {
			m.log.Warn().
				Int64("height", height).
				Str("handler", "HandleEndBlocker").
				Msg("not enough attributes in event")
			continue
		}

		var (
			msgIndex, batchIndex, poolID                   uint32
			swapRequester, offerCoinDenom, demandCoinDenom string
			success                                        bool
			orderExpiryHeight                              int64
			exchangedCoinFeeAmount, orderPrice, swapPrice, offerCoinAmount, exchangedDemandCoinAmount,
			transactedCoinAmount, offerCoinFeeAmount, remainingOfferCoinAmount float64
		)

		for _, attr := range event.Attributes {
			// try to decode value if needed
			switch attr.Key {
			case base64KeyPoolID, base64KeyBatchIndex, base64KeyMsgIndex, base64KeySwapRequester,
				base64KeyOfferCoinDenom, base64KeyExchangedOfferCoinAmount, base64KeyDemandCoinDenom,
				base64KeyOrderPrice, base64KeySwapPrice, base64KeyTransactedCoinAmount,
				base64KeyRemainingOfferCoinAmount, base64KeyExchangedDemandCoinAmount,
				base64KeyOfferCoinFeeAmount, base64KeyExchangedCoinFeeAmount, base64KeyOrderExpiryHeight,
				base64KeySuccess, base64KeyOfferCoinAmount:

				attr.Value, err = utils.DecodeToString(attr.Value)
				if err != nil {
					return err
				}
			}

			switch attr.Key {
			case liquidity.AttributeValuePoolID, base64KeyPoolID:
				var id uint64
				id, err = strconv.ParseUint(attr.Value, 10, 32)
				poolID = uint32(id)
			case liquidity.AttributeValueBatchIndex, base64KeyBatchIndex:
				var index uint64
				index, err = strconv.ParseUint(attr.Value, 10, 32)
				batchIndex = uint32(index)
			case liquidity.AttributeValueMsgIndex, base64KeyMsgIndex:
				var index uint64
				index, err = strconv.ParseUint(attr.Value, 10, 32)
				msgIndex = uint32(index)
			case liquidity.AttributeValueSwapRequester, base64KeySwapRequester:
				swapRequester = attr.Value
			// case liquidity.AttributeValueSwapTypeId:
			case liquidity.AttributeValueOfferCoinDenom, base64KeyOfferCoinDenom:
				offerCoinDenom = attr.Value
			case liquidity.AttributeValueOfferCoinAmount, base64KeyOfferCoinAmount:
				offerCoinAmount, err = strconv.ParseFloat(attr.Value, 64)
			case liquidity.AttributeValueDemandCoinDenom, base64KeyDemandCoinDenom:
				demandCoinDenom = attr.Value
			case liquidity.AttributeValueOrderPrice, base64KeyOrderPrice:
				orderPrice, err = strconv.ParseFloat(attr.Value, 64)
			case liquidity.AttributeValueSwapPrice, base64KeySwapPrice:
				swapPrice, err = strconv.ParseFloat(attr.Value, 64)
			case liquidity.AttributeValueTransactedCoinAmount, base64KeyTransactedCoinAmount:
				transactedCoinAmount, err = strconv.ParseFloat(attr.Value, 64)
			case liquidity.AttributeValueRemainingOfferCoinAmount, base64KeyRemainingOfferCoinAmount:
				remainingOfferCoinAmount, err = strconv.ParseFloat(attr.Value, 64)
			case liquidity.AttributeValueExchangedDemandCoinAmount, base64KeyExchangedDemandCoinAmount:
				exchangedDemandCoinAmount, err = strconv.ParseFloat(attr.Value, 64)
			case liquidity.AttributeValueOfferCoinFeeAmount, base64KeyOfferCoinFeeAmount:
				offerCoinFeeAmount, err = strconv.ParseFloat(attr.Value, 64)
			case liquidity.AttributeValueExchangedCoinFeeAmount, base64KeyExchangedCoinFeeAmount:
				exchangedCoinFeeAmount, err = strconv.ParseFloat(attr.Value, 64)
			// case liquidity.AttributeValueReservedOfferCoinFeeAmount:
			case liquidity.AttributeValueOrderExpiryHeight, base64KeyOrderExpiryHeight:
				orderExpiryHeight, err = strconv.ParseInt(attr.Value, 10, 64)
			case liquidity.AttributeValueSuccess, base64KeySuccess:
				success = attr.Value == liquidity.Success
			}

			if err != nil {
				return errors.Wrap(err, "liquidity: failed to parse event attributes")
			}
		}

		if err = m.broker.PublishSwap(ctx, model.Swap{
			Height:                    height,
			MsgIndex:                  msgIndex,
			BatchIndex:                batchIndex,
			PoolID:                    poolID,
			SwapRequester:             swapRequester,
			OfferCoinDenom:            offerCoinDenom,
			OfferCoinAmount:           offerCoinAmount,
			DemandCoinDenom:           demandCoinDenom,
			ExchangedDemandCoinAmount: exchangedDemandCoinAmount,
			TransactedCoinAmount:      transactedCoinAmount,
			RemainingOfferCoinAmount:  remainingOfferCoinAmount,
			OfferCoinFeeAmount:        offerCoinFeeAmount,
			OrderExpiryHeight:         orderExpiryHeight,
			ExchangedCoinFeeAmount:    exchangedCoinFeeAmount,
			OrderPrice:                orderPrice,
			SwapPrice:                 swapPrice,
			Success:                   success,
		}); err != nil {
			return err
		}
	}

	return nil
}
