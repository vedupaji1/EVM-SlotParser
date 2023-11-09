package slot_parser

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"strings"

// 	"math/big"

// 	"github.com/ethereum/go-ethereum/crypto"

// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/ethclient"
// 	"github.com/mitchellh/mapstructure"
// 	"github.com/spf13/cast"
// 	viperPKG "github.com/spf13/viper"
// )

// func GetStorageInfo(viper *viperPKG.Viper, inputKey []interface{}) (*StorageLayoutInfo, map[string]interface{}) {
// 	storageLayout := cast.ToSlice(viper.Get("storage"))
// 	for _, data := range storageLayout {
// 		var tempStorageInfo StorageLayoutInfo
// 		err := mapstructure.Decode(data, &tempStorageInfo)
// 		if err != nil {
// 			log.Panic("Failed To Decode Storage Info: ", err)
// 		}
// 		if tempStorageInfo.Label == cast.ToString(inputKey[0]) {
// 			return &tempStorageInfo, viper.GetStringMap("types")
// 		}
// 	}
// 	return nil, nil
// }

// func GetParsedStorageTypeData(storageTypeData map[string]interface{}, typeToGet string) *StorageLayoutTypeData {
// 	typeToGetInLowerCase := strings.ToLower(typeToGet)
// 	for key, data := range storageTypeData {
// 		if key == typeToGetInLowerCase {
// 			var tempStorageTypeData StorageLayoutTypeData
// 			err := mapstructure.Decode(data, &tempStorageTypeData)
// 			if err != nil {
// 				log.Panic("Failed To Decode Storage Type Data: ", err)
// 			}
// 			return &tempStorageTypeData
// 		}
// 	}
// 	return nil
// }

// func getBytesOrStringTypeSlot(parentSlotNum common.Hash) []common.Hash {
// 	client, err := ethclient.Dial("https://eth-goerli.g.alchemy.com/v2/barNHxwKcvdxJuDoKlbor5qx6mhT2C_O")
// 	if err != nil {
// 		log.Panic("Failed To Create ETH Client, Err:", err)
// 	}
// 	lastestBlockNum, err := client.BlockNumber(context.Background())
// 	if err != nil {
// 		log.Panic("Failed To Get Latest Block Number, Err:", err)
// 	}
// 	slotData, err := client.StorageAt(context.Background(), common.HexToAddress(ContractAddressStr), parentSlotNum, big.NewInt(int64(lastestBlockNum)))
// 	if err != nil {
// 		log.Panic("Failed To Get Slot Data, Err:", err)
// 	}
// 	slotDataBigInt := new(big.Int).SetBytes(slotData)
// 	if new(big.Int).Mod(slotDataBigInt, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
// 		return []common.Hash{parentSlotNum}
// 	}
// 	expectedConsumedSlotsWithDecimals := new(big.Float).Quo(new(big.Float).SetInt(new(big.Int).Sub(slotDataBigInt, big.NewInt(1))), big.NewFloat(64))
// 	expectedConsumedSlotsInt, _ := expectedConsumedSlotsWithDecimals.Int(nil)
// 	expectedConsumedSlots := new(big.Float).SetInt(expectedConsumedSlotsInt)
// 	if expectedConsumedSlots.Cmp(expectedConsumedSlotsWithDecimals) != 0 {
// 		expectedConsumedSlots.Add(expectedConsumedSlots, big.NewFloat(1))
// 	}
// 	fmt.Println("Expected Consumed Slots For Bytes Data:", expectedConsumedSlots)
// 	expectedConsumedSlotsInt64, _ := expectedConsumedSlots.Int64()
// 	childSlots := []common.Hash{}
// 	var i int64 = 0
// 	for ; i < expectedConsumedSlotsInt64; i++ {
// 		childSlots = append(childSlots, common.BigToHash(new(big.Int).Add(crypto.Keccak256Hash(encodeUint256(parentSlotNum.Big())).Big(), big.NewInt(i))))
// 	}
// 	return childSlots
// }

// func getMapTypeSlot(parentSlot common.Hash, storageTypeParsedData *StorageLayoutTypeData, storageTypeData map[string]interface{}, inputKey []interface{}, curKeyIndex int) []common.Hash {
// 	if len(inputKey)-1 < curKeyIndex {
// 		log.Panic("Insufficient Number Of Keys Are Passed")
// 	}
// 	fmt.Println("Parent Slot:", parentSlot)
// 	parentSlot = crypto.Keccak256Hash(encodeTwoTypes(storageTypeParsedData.Key[2:], "uint256", inputKey[curKeyIndex], parentSlot.Big()))
// 	return GetStorageSlotNum(parentSlot, storageTypeParsedData.Value, storageTypeData, inputKey, curKeyIndex+1)
// }

// func getStaticArrayTypeSlot(parentSlot common.Hash, storageTypeParsedData *StorageLayoutTypeData, storageTypeData map[string]interface{}, inputKey []interface{}, curKeyIndex int) []common.Hash {
// 	if len(inputKey)-1 < curKeyIndex {
// 		log.Panic("Insufficient Number Of Keys Are Passed")
// 	}
// 	fmt.Println("Parent Slot:", parentSlot)
// 	itemToRetrieve, ok := inputKey[curKeyIndex].(int)
// 	if !ok {
// 		log.Panic("Array Index Must Be Int Type")
// 	}
// 	var numberOfItems uint64
// 	if storageTypeParsedData.Encoding == InPlaceEncoding {
// 		numberOfBytesBigInt, ok := new(big.Int).SetString(storageTypeParsedData.NumberOfBytes, 10)
// 		if !ok {
// 			log.Panic("Failed To Convert String To BigInt")
// 		}
// 		numberOfItems = new(big.Int).Div(numberOfBytesBigInt, big.NewInt(32)).Uint64()
// 		parentSlot = common.BigToHash(new(big.Int).Add(parentSlot.Big(), big.NewInt(int64(itemToRetrieve))))
// 	} else {
// 		client, err := ethclient.Dial("https://eth-goerli.g.alchemy.com/v2/barNHxwKcvdxJuDoKlbor5qx6mhT2C_O")
// 		if err != nil {
// 			log.Panic("Failed To Create ETH Client, Err:", err)
// 		}
// 		lastestBlockNum, err := client.BlockNumber(context.Background())
// 		if err != nil {
// 			log.Panic("Failed To Get Latest Block Number, Err:", err)
// 		}
// 		slotData, err := client.StorageAt(context.Background(), common.HexToAddress(ContractAddressStr), parentSlot, big.NewInt(int64(lastestBlockNum)))
// 		if err != nil {
// 			log.Panic("Failed To Get Slot Data, Err:", err)
// 		}
// 		numberOfItems = new(big.Int).SetBytes(slotData).Uint64()
// 		parentSlot = common.BigToHash(new(big.Int).Add(crypto.Keccak256Hash(encodeUint256(parentSlot.Big())).Big(), big.NewInt(int64(itemToRetrieve))))
// 	}
// 	if itemToRetrieve >= int(numberOfItems) {
// 		log.Panic("Index Number For Array Must Be Lesser Than Array Size")
// 	}
// 	return GetStorageSlotNum(parentSlot, storageTypeParsedData.Base, storageTypeData, inputKey, curKeyIndex+1)
// }

// func getStructTypeSlot(parentSlot common.Hash, storageTypeParsedData *StorageLayoutTypeData, storageTypeData map[string]interface{}, inputKey []interface{}, curKeyIndex int) []common.Hash {
// 	if len(inputKey)-1 < curKeyIndex {
// 		log.Panic("Insufficient Number Of Keys Are Passed")
// 	}
// 	fmt.Println("Parent Slot:", parentSlot)
// 	for _, data := range storageTypeParsedData.Members {
// 		if data.Label == inputKey[curKeyIndex] {
// 			childIndex, ok := new(big.Int).SetString(data.Slot, 10)
// 			if !ok {
// 				log.Panic("Failed To Convert String To BigInt")
// 			}
// 			parentSlot = common.BigToHash(new(big.Int).Add(parentSlot.Big(), childIndex))
// 		}
// 	}
// 	return GetStorageSlotNum(parentSlot, storageTypeParsedData.Base, storageTypeData, inputKey, curKeyIndex+1)
// }

// func GetStorageSlotNum(parentSlot common.Hash, storageType string, storageTypeData map[string]interface{}, inputKey []interface{}, nextKeyIndex int) []common.Hash {
// 	storageTypeParsedData := GetParsedStorageTypeData(storageTypeData, storageType)
// 	if len(storageType) >= MinStorageTypeLen {
// 		if storageType[:ArrayAndBytesStorageTypeLen] == ArrayTypeStr {
// 			return getStaticArrayTypeSlot(parentSlot, storageTypeParsedData, storageTypeData, inputKey, nextKeyIndex)
// 		} else if storageType[:ArrayAndBytesStorageTypeLen] == BytesTypeStr || storageType[:StructAndStringStorageTypeLen] == StringTypeStr {
// 			return getBytesOrStringTypeSlot(parentSlot)
// 		} else if storageType[:StructAndStringStorageTypeLen] == StructTypeStr {
// 			return getStructTypeSlot(parentSlot, storageTypeParsedData, storageTypeData, inputKey, nextKeyIndex)
// 		} else if storageType[:MapStorageTypeLen] == MapTypeStr {
// 			return getMapTypeSlot(parentSlot, storageTypeParsedData, storageTypeData, inputKey, nextKeyIndex)
// 		} else {
// 			return []common.Hash{parentSlot}
// 		}
// 	} else {
// 		return []common.Hash{parentSlot}
// 	}
// }

// func main() {
// 	viper := viperPKG.New()
// 	viper.SetConfigFile("storageLayout.json")
// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		log.Panic("Failed To Read Storage Layout File: ", err)
// 	}
// 	inputKey := []interface{}{"tempBytes_2", 1, "id"}
// 	fmt.Println(inputKey)
// 	storageInfo, storageTypeData := GetStorageInfo(viper, inputKey)
// 	if storageInfo == nil {
// 		fmt.Println("StorageVariable Not Found")
// 		return
// 	}

// 	fmt.Printf("Storage Info: %#v\n", storageInfo)
// 	fmt.Println("Storage Slot Num:", GetStorageSlotNum(strToCommonHash(storageInfo.Slot), storageInfo.Type, storageTypeData, inputKey, 1))
// }
