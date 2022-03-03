package flashbots

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	flashbots "github.com/0xblocks/flashbots-bundle"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Bot struct {
	Client     *ethclient.Client
	SigningKey *ecdsa.PrivateKey
	PrivateKey *ecdsa.PrivateKey
	Address    *common.Address
}

func NewBot(signingKey string, privateKey string, providerURL string) (*Bot, error) {
	var err error
	bot := Bot{}

	bot.Client, err = ethclient.Dial(providerURL)
	if err != nil {
		return nil, err
	}

	bot.SigningKey, err = crypto.HexToECDSA(signingKey)
	if err != nil {
		return nil, err
	}

	bot.PrivateKey, err = crypto.HexToECDSA(privateKey)
	fromAddress := crypto.PubkeyToAddress(bot.PrivateKey.PublicKey)
	bot.Address = &fromAddress

	if err != nil {
		return nil, err
	}

	return &bot, nil
}

func (bot *Bot) BuildTransaction(to *common.Address, value *big.Int, data []byte, nonce uint64, zeroGas bool) (*types.Transaction, error) {

	var err error
	if nonce == 0 {
		nonce, err = bot.Client.PendingNonceAt(context.Background(), *bot.Address)
		if err != nil {
			return nil, err
		}
	}

	gasPrice := big.NewInt(0)
	if !zeroGas {
		gasPrice, err = bot.Client.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, err
		}
	}

	gasLimit, err := bot.Client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: *bot.Address,
		To:   to,
		Data: data,
	})
	if err != nil {
		log.Printf("Caught error estimating gas: %s\n", err)
		gasLimit = 100000
	}

	chainID, err := bot.Client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(nonce, *to, value, gasLimit*2, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), bot.PrivateKey)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func (bot *Bot) SendFlashbotsBundle(signedTxs []*types.Transaction, attempts int64) error {

	//start := time.Now()
	fb := flashbots.NewProvider(bot.SigningKey, bot.PrivateKey, flashbots.DefaultRelayURL)

	txs := []string{}
	for _, tx := range signedTxs {
		data, err := tx.MarshalBinary()
		if err != nil {
			return err
		}
		txs = append(txs, hexutil.Encode(data))
	}

	block, err := bot.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}

	// First simulate the transaction to make sure it doesn't revert
	// resp, err := fb.Simulate(txs, block.Number, "latest", 0)
	// if err != nil {
	// 	return err
	// }

	// err = resp.HasError()
	// if err != nil {
	// 	return err
	// }

	// cb, _ := new(big.Float).SetString(resp.Result.CoinbaseDiff)
	// eth := new(big.Float).Quo(cb, big.NewFloat(math.Pow10(18)))
	// wei, _ := resp.EffectiveGasPrice()
	// gwei := util.ToGwei(wei)

	// fmt.Printf("Simulation completed in %fs. Cost: %f eth, effective price: %d gwei\n", time.Since(start).Seconds(), eth, gwei)

	// Send the bundle to the FB relay for inclusion in a block
	// Unless your tx is only valid for a single block, you should target several blocks since flashbots
	// isn't available on 100% of the hashing power. On Goerli, it's only available on a single validator,
	// so target at least 9 blocks to ensure it gets picked up
	for i := int64(0); i < attempts; i++ {
		targetBlockNumber := new(big.Int).Add(block.Number, big.NewInt(int64(i)))
		_, err := fb.SendBundle(txs, targetBlockNumber, &flashbots.Options{})
		if err != nil {
			return err
		}
		fmt.Printf("submitted for block: %d\n", targetBlockNumber)
	}
	return nil
}

func (bot *Bot) WaitForTx(tx *types.Transaction, maxWaitSeconds uint) error {

	log.Println("Waiting for tx to complete...")
	loops := uint(0)
	for {
		_, isPending, err := bot.Client.TransactionByHash(context.Background(), tx.Hash())
		if err != nil && err != ethereum.NotFound {
			return err
		}

		if !isPending {
			// It's not pending, so check if it's been mined
			receipt, err := bot.Client.TransactionReceipt(context.Background(), tx.Hash())
			if err != nil && err != ethereum.NotFound {
				return err
			}
			if receipt != nil {
				log.Println("tx complete. it may take a few minutes to appear in etherscan.")
				log.Printf("https://etherscan.io/tx/%s\n", tx.Hash().Hex())
				break
			}
		}

		time.Sleep(1 * time.Second)

		loops = loops + 1
		if loops > maxWaitSeconds {
			log.Printf("timed out after %d seconds. check manually here:\n", maxWaitSeconds)
			log.Printf("https://etherscan.io/tx/%s\n", tx.Hash().Hex())
			break
		}
	}

	return nil
}
