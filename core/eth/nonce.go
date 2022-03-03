package eth

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

func (e *EthChain) Nonce(spenderAddressHex string) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	nonce, err := e.RemoteRpcClient.PendingNonceAt(ctx, common.HexToAddress(spenderAddressHex))
	if err != nil {
		return uint64(0), err
	}

	return nonce, nil
}
