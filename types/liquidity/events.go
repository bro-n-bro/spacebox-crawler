package liquidity

// Event types for the liquidity module.
const (
	EventTypeSwapTransacted = "swap_transacted"

	AttributeValuePoolID     = "pool_id"
	AttributeValueBatchIndex = "batch_index"
	AttributeValueMsgIndex   = "msg_index"

	AttributeValueOfferCoinDenom         = "offer_coin_denom"
	AttributeValueOfferCoinAmount        = "offer_coin_amount"
	AttributeValueOfferCoinFeeAmount     = "offer_coin_fee_amount"
	AttributeValueExchangedCoinFeeAmount = "exchanged_coin_fee_amount"
	AttributeValueDemandCoinDenom        = "demand_coin_denom"
	AttributeValueOrderPrice             = "order_price"

	AttributeValueSuccess       = "success"
	AttributeValueSwapRequester = "swap_requester"
	AttributeValueSwapPrice     = "swap_price"

	AttributeValueTransactedCoinAmount      = "transacted_coin_amount"
	AttributeValueRemainingOfferCoinAmount  = "remaining_offer_coin_amount"
	AttributeValueExchangedOfferCoinAmount  = "exchanged_offer_coin_amount"
	AttributeValueExchangedDemandCoinAmount = "exchanged_demand_coin_amount"
	AttributeValueOrderExpiryHeight         = "order_expiry_height"

	Success = "success"
)
