package eth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
)

func (e *EthChain) SuggestGasPrice() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	gasPrice, err := e.RemoteRpcClient.SuggestGasPrice(ctx)

	if err != nil {
		return big.NewInt(0), nil
	}
	return gasPrice, err
}

func (e *EthChain) EstimateGas(msg ethereum.CallMsg) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	gasCount, err := e.RemoteRpcClient.EstimateGas(ctx, msg)
	if err != nil {
		return 0, err
	}
	return gasCount, nil
}
