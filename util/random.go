package util

import (
	"errors"
	"math/rand"
	"time"
)

type Random struct {
}

var RandomInstance = Random{}

func (randomInstance *Random) RandomFromStringSlice(slice []string) string {
	rand.Seed(time.Now().UnixNano())
	return slice[rand.Intn(len(slice))]
}

func (randomInstance *Random) MustRandomInt64(start int64, end int64) int64 {
	re, err := randomInstance.RandomInt64(start, end)
	if err != nil {
		panic(err)
	}
	return re
}

func (randomInstance *Random) RandomInt64(start int64, end int64) (int64, error) {
	rand.Seed(time.Now().UnixNano())
	if end <= start {
		return 0, errors.New(`end must gt start`)
	}
	return rand.Int63n(end-start) + start, nil
}

func (randomInstance *Random) MustRandomString(count int32) string {
	result, err := randomInstance.RandomString(count)
	if err != nil {
		panic(err)
	}
	return result
}

func (randomInstance *Random) RandomString(count int32) (string, error) {
	return randomInstance.RandomStringFromDic("_0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", count)
}

func (randomInstance *Random) RandomStringFromDic(dictionary string, count int32) (string, error) {
	b := make([]byte, count)
	l := len(dictionary)

	_, err := rand.New(rand.NewSource(time.Now().UnixNano())).Read(b)

	if err != nil {
		return "", err
	}
	for i, v := range b {
		b[i] = dictionary[v%byte(l)]
	}

	return string(b), nil
}

func (randomInstance *Random) MustRandomBytes(count int32) []byte {
	result, err := randomInstance.RandomBytes(count)
	if err != nil {
		panic(err)
	}
	return result
}

func (randomInstance *Random) RandomBytes(count int32) ([]byte, error) {
	b := make([]byte, count)

	_, err := rand.New(rand.NewSource(time.Now().UnixNano())).Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func (randomInstance *Random) RandomNumberStr(count int32) (string, error) {
	return randomInstance.RandomStringFromDic("0123456789", count)
}

func (randomInstance *Random) MustRandomNumberStr(count int32) string {
	result, err := randomInstance.RandomNumberStr(count)
	if err != nil {
		panic(err)
	}
	return result
}
