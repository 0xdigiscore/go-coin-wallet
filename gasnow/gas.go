package gasnow

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/0xhelloweb3/go-coin-wallet/constants"
)

type GasNowResult struct {
	Data map[string]float64 `json:"data"`
	Code string             `json:"origin"`
}

type EstimateGasPrice struct {
	Fast     *big.Int
	Rapid    *big.Int
	Slow     *big.Int
	Standard *big.Int
}

func GetEstimateGasPrice() (*EstimateGasPrice, error) {
	resp, err := http.Get(constants.GAS_NOW_URL)
	if err != nil {
		return &EstimateGasPrice{}, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var estimate GasNowResult

	_ = json.Unmarshal(body, &estimate)

	estimateGasPrice := &EstimateGasPrice{
		Fast:     big.NewInt(int64(estimate.Data["fast"])),
		Slow:     big.NewInt(int64(estimate.Data["slow"])),
		Rapid:    big.NewInt(int64(estimate.Data["rapid"])),
		Standard: big.NewInt(int64(estimate.Data["standard"])),
	}

	return estimateGasPrice, err
}
