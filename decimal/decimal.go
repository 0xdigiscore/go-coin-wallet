package decimal

import (
	"fmt"
	"math/big"
	"strings"

	go_reflect "github.com/pefish/go-reflect"
	"github.com/pkg/errors"
	decimal "github.com/qinghuan-chain/go-coin-wallet/util"
)

type DecimalClass struct {
	result decimal.Decimal
}

var Decimal = DecimalClass{}

// =
func (decimalInstance *DecimalClass) Eq(a interface{}) bool {
	return decimalInstance.result.Equal(decimalInstance.interfaceToDecimal(a))
}

// !=
func (decimalInstance *DecimalClass) Neq(a interface{}) bool {
	return !decimalInstance.result.Equal(decimalInstance.interfaceToDecimal(a))
}

// <
func (decimalInstance *DecimalClass) Lt(a interface{}) bool {
	return decimalInstance.result.LessThan(decimalInstance.interfaceToDecimal(a))
}

// <=
func (decimalInstance *DecimalClass) Lte(a interface{}) bool {
	return decimalInstance.result.LessThanOrEqual(decimalInstance.interfaceToDecimal(a))
}

// >
func (decimalInstance *DecimalClass) Gt(a interface{}) bool {
	return decimalInstance.result.GreaterThan(decimalInstance.interfaceToDecimal(a))
}

// >=
func (decimalInstance *DecimalClass) Gte(a interface{}) bool {
	return decimalInstance.result.GreaterThanOrEqual(decimalInstance.interfaceToDecimal(a))
}

// 开始计算。小数后面有0的话会自动去除
func (decimalInstance *DecimalClass) Start(a interface{}) *DecimalClass {
	decimalInstanceNew := DecimalClass{}
	decimalInstanceNew.result = decimalInstanceNew.interfaceToDecimal(a)
	return &decimalInstanceNew
}

// \-1\ = 1
func (decimalInstance *DecimalClass) AbsForString() string {
	return decimalInstance.Abs().result.String()
}

// \-1\ = 1
func (decimalInstance *DecimalClass) Abs() *DecimalClass {
	decimalInstance.result = decimalInstance.result.Abs()
	return decimalInstance
}

// +
func (decimalInstance *DecimalClass) AddForString(a interface{}) string {
	return decimalInstance.Add(decimalInstance.interfaceToDecimal(a)).result.String()
}

// +
func (decimalInstance *DecimalClass) Add(a interface{}) *DecimalClass {
	decimalInstance.result = decimalInstance.result.Add(decimalInstance.interfaceToDecimal(a))
	return decimalInstance
}

// -
func (decimalInstance *DecimalClass) SubForString(a interface{}) string {
	return decimalInstance.Sub(decimalInstance.interfaceToDecimal(a)).result.String()
}

// -
func (decimalInstance *DecimalClass) Sub(a interface{}) *DecimalClass {
	decimalInstance.result = decimalInstance.result.Sub(decimalInstance.interfaceToDecimal(a))
	return decimalInstance
}

// /
func (decimalInstance *DecimalClass) DivForString(a interface{}) string {
	return decimalInstance.Div(decimalInstance.interfaceToDecimal(a)).result.String()
}

// /
func (decimalInstance *DecimalClass) Div(a interface{}) *DecimalClass {
	decimalInstance.result = decimalInstance.result.Div(decimalInstance.interfaceToDecimal(a))
	return decimalInstance
}

func (decimalInstance *DecimalClass) MustShiftedBy(a interface{}) *DecimalClass {
	result, err := decimalInstance.ShiftedBy(a)
	if err != nil {
		panic(err)
	}
	return result
}

func (decimalInstance *DecimalClass) ShiftedBy(a interface{}) (*DecimalClass, error) {
	int32_, err := go_reflect.Reflect.ToInt32(a)
	if err != nil {
		return nil, err
	}
	decimalInstance.result = decimalInstance.result.Shift(int32_)
	return decimalInstance, nil
}

func (decimalInstance *DecimalClass) MustUnShiftedBy(a interface{}) *DecimalClass {
	result, err := decimalInstance.UnShiftedBy(a)
	if err != nil {
		panic(err)
	}
	return result
}

func (decimalInstance *DecimalClass) UnShiftedBy(a interface{}) (*DecimalClass, error) {
	int32_, err := go_reflect.Reflect.ToInt32(a)
	if err != nil {
		return nil, err
	}
	decimalInstance.result = decimalInstance.result.Shift(-int32_)
	return decimalInstance, nil
}

// *
func (decimalInstance *DecimalClass) MultiForString(a interface{}) string {
	return decimalInstance.Multi(decimalInstance.interfaceToDecimal(a)).result.String()
}

// *
func (decimalInstance *DecimalClass) Multi(a interface{}) *DecimalClass {
	decimalInstance.result = decimalInstance.result.Mul(decimalInstance.interfaceToDecimal(a))
	return decimalInstance
}

func (decimalInstance *DecimalClass) End() decimal.Decimal {
	return decimalInstance.result
}

func (decimalInstance *DecimalClass) EndForString() string {
	return decimalInstance.result.String()
}

func (decimalInstance *DecimalClass) EndForBigInt() *big.Int {
	result, ok := new(big.Int).SetString(decimalInstance.result.String(), 10)
	if !ok {
		panic(errors.New(fmt.Sprintf("string %s to bigInt error", decimalInstance.result.String())))
	}
	return result
}

// 直接截取
func (decimalInstance *DecimalClass) TruncForString(precision int32) string {
	return decimalInstance.Trunc(precision).result.String()
}

// 直接截取
func (decimalInstance *DecimalClass) Trunc(precision int32) *DecimalClass {
	decimalInstance.result = decimalInstance.result.Truncate(precision)
	return decimalInstance
}

// 四舍五入
func (decimalInstance *DecimalClass) RoundForString(precision int32) string {
	return decimalInstance.Round(precision).result.String()
}

// 四舍五入（后面保留0）
func (decimalInstance *DecimalClass) RoundForRemainZeroString(precision int32) string {
	return decimalInstance.Round(precision).result.StringRemain()
}

// 四舍五入
func (decimalInstance *DecimalClass) Round(precision int32) *DecimalClass {
	decimalInstance.result = decimalInstance.result.Round(precision)
	return decimalInstance
}

// 向上截取
func (decimalInstance *DecimalClass) RoundUpForString(precision int32) string {
	return decimalInstance.RoundUp(precision).result.String()
}

// 向上截取
func (decimalInstance *DecimalClass) RoundUp(precision int32) *DecimalClass {
	if decimalInstance.result.Round(precision).Equal(decimalInstance.result) {
		return decimalInstance
	}

	halfPrecision := decimal.New(5, -precision-1)

	decimalInstance.result = decimalInstance.result.Add(halfPrecision).Round(precision)
	return decimalInstance
}

// 向下截取
func (decimalInstance *DecimalClass) RoundDownForString(precision int32) string {
	return decimalInstance.RoundDown(precision).result.String()
}

// 向下截取
func (decimalInstance *DecimalClass) RoundDown(precision int32) *DecimalClass {
	if decimalInstance.result.Round(precision).Equal(decimalInstance.result) {
		return decimalInstance
	}

	halfPrecision := decimal.New(5, -precision-1)

	decimalInstance.result = decimalInstance.result.Sub(halfPrecision).Round(precision)
	return decimalInstance
}

func (decimalInstance *DecimalClass) interfaceToDecimal(a interface{}) decimal.Decimal {
	if inst, ok := a.(decimal.Decimal); ok {
		return inst
	}
	if inst, ok := a.(*DecimalClass); ok {
		return inst.result
	}
	if inst, ok := a.(DecimalClass); ok {
		return inst.result
	}
	str := ""
	if inst, ok := a.(*big.Int); ok {
		str = inst.String()
	} else {
		str = go_reflect.Reflect.ToString(a)
	}

	decimal_, err := decimal.NewFromString(str)
	if err != nil {
		panic(err)
	}
	return decimal_
}

// 判断小数的精度是不是指定精度
func (decimalInstance *DecimalClass) IsPrecision(precision int32) bool {
	return decimalInstance.GetPrecision() == precision
}

func (decimalInstance *DecimalClass) GetPrecision() int32 {
	splitAmount := strings.Split(decimalInstance.result.String(), `.`)
	if len(splitAmount) <= 1 {
		return 0
	} else {
		return int32(len(splitAmount[1]))
	}
}
