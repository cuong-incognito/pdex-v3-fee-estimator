package main

import (
	"math/big"

	"github.com/cuong-incognito/pdex-v3-fee-estimator/jsonresult"
	"github.com/incognitochain/incognito-chain/common"
)

func getPriceAgainstPRV(sellToken string, pdexState jsonresult.PdexState) [3]*big.Int {
	result := [3]*big.Int{new(big.Int).SetInt64(0), new(big.Int).SetInt64(0), new(big.Int).SetInt64(0)}
	for _, pair := range pdexState.PoolPairs {
		var tokenID string
		virtualTokenReserve := big.NewInt(0)
		virtualPRVReserve := big.NewInt(0)
		if pair.State.Token0ID == common.PRVIDStr {
			tokenID = pair.State.Token1ID
			virtualTokenReserve.Set(pair.State.Token1VirtualAmount)
			virtualPRVReserve.Set(pair.State.Token0VirtualAmount)
		} else if pair.State.Token1ID == common.PRVIDStr {
			tokenID = pair.State.Token0ID
			virtualTokenReserve.Set(pair.State.Token0VirtualAmount)
			virtualPRVReserve.Set(pair.State.Token1VirtualAmount)
		}

		if tokenID != sellToken {
			continue
		}

		normalizedLiquidity := big.NewInt(0).Mul(virtualTokenReserve, virtualPRVReserve)
		normalizedLiquidity.Mul(normalizedLiquidity, big.NewInt(BaseAmplifier))
		normalizedLiquidity.Div(normalizedLiquidity, big.NewInt(0).SetUint64(uint64(pair.State.Amplifier)))
		normalizedLiquidity.Mul(normalizedLiquidity, big.NewInt(BaseAmplifier))
		normalizedLiquidity.Div(normalizedLiquidity, big.NewInt(0).SetUint64(uint64(pair.State.Amplifier)))

		if normalizedLiquidity.Cmp(result[2]) == 1 {
			result = [3]*big.Int{virtualTokenReserve, virtualPRVReserve, normalizedLiquidity}
		}
	}
	return result
}
