package tobroker

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (tb ToBroker) MapTransaction(tx *types.Tx) (model.Transaction, error) {
	signatures := make([]string, len(tx.Signatures))
	for i, s := range tx.Signatures {
		signatures[i] = string(s)
	}

	var msgs = make([][]byte, len(tx.Body.Messages))
	for i, msg := range tx.Body.Messages {
		msgBytes, err := tb.cdc.MarshalJSON(msg)
		if err != nil {
			return model.Transaction{}, err
		}
		msgs[i] = msgBytes
	}

	logs, err := tb.amino.MarshalJSON(tx.Logs)
	if err != nil {
		return model.Transaction{}, err
	}

	t := model.Transaction{
		Hash:       tx.TxHash,
		Height:     tx.Height,
		Success:    tx.Successful(),
		Messages:   msgs,
		Memo:       tx.Body.Memo,
		Signatures: signatures,
		Signer:     tx.Signer,
		GasWanted:  tx.GasWanted,
		GasUsed:    tx.GasUsed,
		RawLog:     tx.RawLog,
		Logs:       logs,
	}

	if tx.AuthInfo != nil {
		if tx.AuthInfo.SignerInfos != nil {
			infos := make([]model.SignersInfo, len(tx.AuthInfo.SignerInfos))
			for i, info := range tx.AuthInfo.SignerInfos {
				// info.ModeInfo // TODO: add it
				infos[i] = model.SignersInfo{
					PublicKey: info.PublicKey.String(),
					Sequence:  info.Sequence,
				}
			}
			t.SignerInfos = infos
		}

		if tx.AuthInfo.Fee != nil {
			var payer string
			if tx.AuthInfo.Fee.Payer == "" && len(tx.Body.Messages) > 0 {
				// XXX
				// without this we will get a panic if transaction cannot contain a feePayer
				var stdMsg sdk.Msg
				if err = tb.cdc.UnpackAny(tx.Body.Messages[0], &stdMsg); err == nil {
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
		}
	}

	return t, nil
}
