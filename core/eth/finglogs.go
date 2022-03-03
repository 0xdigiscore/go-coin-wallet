package eth

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

/*
查找历史的已经确认的事件，但不能实时接受后面的事件。取不到 pending 中的 logs
fromBlock；nil 就是 最新块号 - 4900，负数就是 pending ，正数就是 blockNumber
toBlock；nil 就是 latest ，负数就是 pending ，正数就是 blockNumber
*/
func (e *EthChain) FindLogs(contractAddress, abiStr, eventName string, fromBlock, toBlock *big.Int, query ...[]interface{}) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	parsedAbi, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, err
	}

	query = append([][]interface{}{{parsedAbi.Events[eventName].ID}}, query...)

	topics, err := abi.MakeTopics(query...)
	if err != nil {
		return nil, err
	}
	contractInstance := bind.NewBoundContract(common.HexToAddress(contractAddress), parsedAbi,
		e.RemoteRpcClient, e.RemoteRpcClient, e.RemoteRpcClient)

	if fromBlock == nil {
		number, err := e.LatestBlockNumber()
		if err != nil {
			return nil, err
		}
		fromBlock = number.Sub(number, new(big.Int).SetUint64(4900))
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	logs, err := e.RemoteRpcClient.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: []common.Address{
			common.HexToAddress(contractAddress),
		},
		Topics: topics,
	})
	if err != nil {
		return nil, err
	}
	for _, log := range logs {
		map_ := make(map[string]interface{})
		err := contractInstance.UnpackLogIntoMap(map_, eventName, log)
		if err != nil {
			return nil, err
		}
		result = append(result, map_)
	}
	return result, nil
}
