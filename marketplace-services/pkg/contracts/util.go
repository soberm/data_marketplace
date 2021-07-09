package contracts

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

func ContainsLocation(locations []Location, location Location) bool {
	for _, l := range locations {
		if l == location {
			return true
		}
	}
	return false
}

func AddressesToHex(addresses []common.Address) []string {
	var hex []string
	hex = make([]string, len(addresses))
	for i, address := range addresses {
		hex[i] = address.Hex()
	}
	return hex
}

func HexToAddresses(hex []string) []common.Address {
	var addresses []common.Address
	addresses = make([]common.Address, len(hex))
	for i, address := range hex {
		addresses[i] = common.HexToAddress(address)
	}
	return addresses
}

func BigIntToInt64(values []*big.Int) []int64 {
	data := make([]int64, len(values))
	for i := range data {
		data[i] = values[i].Int64()
	}
	return data
}

func Int64ToBigInt(values []int64) []*big.Int {
	data := make([]*big.Int, len(values))
	for i := range data {
		data[i] = big.NewInt(values[i])
	}
	return data
}

func UInt64ToBigInt(values []uint64) []*big.Int {
	data := make([]*big.Int, len(values))
	for i := range data {
		data[i] = big.NewInt(int64(values[i]))
	}
	return data
}
