package util

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func DecodeRawTx(rawTx string) (*types.Transaction, error) {
	rawTxBytes, err := hex.DecodeString(rawTx)
	if err != nil {
		return nil, err
	}
	tx := new(types.Transaction)
	rlp.DecodeBytes(rawTxBytes, &tx)
	return tx, err
}

func EncodeTx(txn *types.Transaction) ([]byte, error) {
	rawTxBytes, err := rlp.EncodeToBytes(txn)
	return rawTxBytes, err
}
