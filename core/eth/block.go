package eth

import (
	"context"
	"math/big"
)

func (e *EthChain) LatestBlockNumber() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	number, err := e.RemoteRpcClient.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	return new(big.Int).SetUint64(number), nil
}
