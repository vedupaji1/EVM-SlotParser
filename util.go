package slot_parser

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/umbracle/ethgo/abi"
)

func encodeUint256(data *big.Int) []byte {
	uint256Type := abi.MustNewType("uint256")
	uint256TypeBytes, err := uint256Type.Encode(data)
	if err != nil {
		log.Panic(err)
	}
	return uint256TypeBytes
}

func encodeTwoTypes(type1 string, type2 string, type1Data interface{}, type2Data interface{}) []byte {
	abiType := abi.MustNewType(fmt.Sprintf("tuple(%v type1, %v type2)", type1, type2))
	abiTypeBytes, err := abiType.Encode(map[string]interface{}{
		"type1": type1Data,
		"type2": type2Data,
	})
	if err != nil {
		log.Panic(err)
	}
	return abiTypeBytes
}

func strToCommonHash(str string) common.Hash {
	tempBigInt, ok := new(big.Int).SetString(str, 10)
	if !ok {
		log.Panic("Failed To Convert String To BigInt")
	}
	return common.BigToHash(tempBigInt)
}

func IsNonPrimitiveStorageType(storageType string) bool {
	return len(storageType) <= MinStorageTypeLen ||
		storageType[:ArrayAndBytesStorageTypeLen] != ArrayTypeStr && storageType[:ArrayAndBytesStorageTypeLen] != BytesTypeStr &&
			storageType[:StructAndStringStorageTypeLen] != StructTypeStr && storageType[:StructAndStringStorageTypeLen] != StringTypeStr &&
			storageType[:MapStorageTypeLen] != MapTypeStr
}

func BytesTo(data []byte, resultingType StorageHighLevelTypes) interface{} {
	if resultingType == IntegerStorageType {
		return new(big.Int).SetBytes(data)
	} else if resultingType == StringStorageType {
		return string(data)
	} else if resultingType == BytesStorageType {
		return hexutil.Encode(data)
	} else {
		log.Panic("Something Went Wrong Invalid Type Is Passed")
		return data
	}
}
