package to_broker

import (
	"bro-n-bro-osmosis/adapter/broker/model"
	"bro-n-bro-osmosis/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (tb ToBroker) MapTransaction(tx *types.Tx) model.Transaction {
	signatures := make([]string, len(tx.Signatures))
	for i, s := range tx.Signatures {
		signatures[i] = string(s)
	}

	t := model.Transaction{
		Hash:       tx.TxHash,
		Height:     tx.Height,
		Success:    tx.Successful(),
		Messages:   nil, // TODO
		Memo:       tx.Body.Memo,
		Signatures: signatures,
		Signer:     "", // TODO
		GasWanted:  tx.GasWanted,
		GasUsed:    tx.GasUsed,
		RawLog:     tx.RawLog,
		Logs:       nil, // TODO
	}

	if tx.AuthInfo != nil {
		if tx.AuthInfo.SignerInfos != nil {
			infos := make([]model.SignersInfo, len(tx.AuthInfo.SignerInfos))
			for i, info := range tx.AuthInfo.SignerInfos {
				// info.ModeInfo // TODO
				infos[i] = model.SignersInfo{
					PublicKey: info.PublicKey.String(),
					Sequence:  info.Sequence,
				}
			}
			t.SignerInfos = infos
		}

		if tx.AuthInfo.Fee != nil {
			var payer string
			if tx.AuthInfo.Fee.Payer == "" {
				// XXX
				// without this we will get a panic if transaction cannot contain a feePayer
				var stdMsg sdk.Msg
				if err := tb.cdc.UnpackAny(tx.Body.Messages[0], &stdMsg); err == nil {
					payer = stdMsg.GetSigners()[0].String()
				}
			} else {
				payer = tx.FeePayer().String()
			}

			t.Fee = &model.Fee{
				Coins:    tb.MapCoins(types.NewCoinsFromCdk(tx.GetFee())),
				GasLimit: tx.GetGas(),
				Granter:  tx.FeeGranter().String(),
				Payer:    payer,
			}

			if payer := tx.AuthInfo.Fee.Payer; payer != "" {
				t.Fee.Payer = payer
			}
		}
	}

	return t
}
