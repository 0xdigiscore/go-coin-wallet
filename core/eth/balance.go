package eth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func (e *EthChain) Balance(address string) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	result, err := e.RemoteRpcClient.BalanceAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}
