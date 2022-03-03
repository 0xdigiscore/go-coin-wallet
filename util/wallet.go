package util

import (
	"math"
	"math/big"
)

func ToFloat(val *big.Int) *big.Float {

	fval := new(big.Float)
	fval.SetString(val.String())
	return new(big.Float).Quo(fval, big.NewFloat(math.Pow10(int(18))))
}

func ToGwei(wei *big.Int) *big.Int {

	return new(big.Int).Div(wei, big.NewInt(1000000000))
}
