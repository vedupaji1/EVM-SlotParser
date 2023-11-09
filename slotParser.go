package slot_parser

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	viperPKG "github.com/spf13/viper"
)

func NewSlotParer(contractName string, contractPath string, contractAddress string, rpcURL string) (*SlotParser, error) {
	if _, err := os.Stat(contractPath); err != nil {
		return nil, fmt.Errorf("invalid path of contract have been passed")
	}
	solcCheckCommand := exec.Command("solc", "--help")
	_, err := solcCheckCommand.Output()
	if err != nil {
		return nil, fmt.Errorf("solc compiler not found, make sure your system have solc compiler, err: %v", err)
	}
	currentWorkingDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory, err: %v", err)
	}
	storageLayoutGenCmd := exec.Command("solc", "--storage-layout", "--pretty-json", "-o", currentWorkingDir+StorageLayoutDirName, "--overwrite", contractPath)
	_, err = storageLayoutGenCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to generate storage layout file using solc compiler, err: %v", err)
	}
	storageLayoutDestDirPath := currentWorkingDir + StorageLayoutDirName
	storageLayoutDirFiles, err := os.ReadDir(storageLayoutDestDirPath)
	if err != nil {
		return nil, fmt.Errorf("filed to open storage layout files dir, err: %v", err)
	}
	isContractStorageLayoutExists := false
	for _, file := range storageLayoutDirFiles {
		if file.Name()[:len(contractName)] == contractName {
			isContractStorageLayoutExists = true
			break
		}
	}
	if !isContractStorageLayoutExists {
		return nil, fmt.Errorf("something went wrong, failed to get storage layout file")
	}

	viper := viperPKG.New()
	viper.SetConfigFile(storageLayoutDestDirPath + "/" + contractName + StorageLayoutFilePaddingStr)
	err = viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read storage layout file, err: %v", err)
	}
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create ETH client, err: %v", err)
	}
	slotParser := &SlotParser{
		ethClient:       client,
		contractAddress: contractAddress,
		viper:           viper,
	}
	return slotParser, nil
}

func (slotParser *SlotParser) GetSlotNum(inputKey []interface{}) (*SlotParserResponse, error) {
	if len(inputKey) < 1 {
		return nil, fmt.Errorf("zero input key passed")
	}
	storageInfo, storageTypeData, err := slotParser.GetStorageInfo(inputKey)
	fmt.Println(inputKey)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Storage Info: %#v\n", storageInfo)
	return slotParser.GetStorageSlotNum(strToCommonHash(storageInfo.Slot), storageInfo.Type, storageInfo.Offset, storageTypeData, inputKey, 1)
}

func (slotParser *SlotParser) GetStorageInfo(inputKey []interface{}) (*StorageLayoutInfo, map[string]interface{}, error) {
	storageLayout := cast.ToSlice(slotParser.viper.Get("storage"))
	for _, data := range storageLayout {
		storageLabel, ok := data.(map[string]interface{})
		if !ok {
			return nil, nil, fmt.Errorf("failed to convert storage layout any type to map type")
		}
		storageLabelStr, ok := storageLabel["label"].((string))
		if !ok {
			return nil, nil, fmt.Errorf("failed to convert any type to string type")
		}
		if storageLabelStr == cast.ToString(inputKey[0]) {
			var tempStorageInfo StorageLayoutInfo
			err := mapstructure.Decode(data, &tempStorageInfo)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to decode storage info: %v", err)
			}
			return &tempStorageInfo, slotParser.viper.GetStringMap("types"), nil
		}
	}
	return nil, nil, fmt.Errorf("storage variable not found")
}

func (slotParser *SlotParser) GetParsedStorageTypeData(storageTypeData map[string]interface{}, typeToGet string) (*StorageLayoutTypeData, error) {
	typeToGetInLowerCase := strings.ToLower(typeToGet)
	for key, data := range storageTypeData {
		if key == typeToGetInLowerCase {
			var tempStorageTypeData StorageLayoutTypeData
			err := mapstructure.Decode(data, &tempStorageTypeData)
			if err != nil {
				return nil, fmt.Errorf("failed to decode storage type data, err: %v", err)
			}
			return &tempStorageTypeData, nil
		}
	}
	return nil, nil
}

func (slotParser *SlotParser) getSlotData(slotNum common.Hash, offset int, numberOfBytesOfType int) ([]byte, error) {
	lastestBlockNum, err := slotParser.ethClient.BlockNumber(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block number, rrr: %v", err)
	}
	slotData, err := slotParser.ethClient.StorageAt(context.Background(), common.HexToAddress(slotParser.contractAddress), slotNum, big.NewInt(int64(lastestBlockNum)))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch slot data, Err: %v", err)
	}
	return slotData[32-numberOfBytesOfType-offset : 32-offset], nil
}

func (slotParser *SlotParser) getStaticArrayTypeSlot(parentSlot common.Hash, storageTypeParsedData *StorageLayoutTypeData, storageTypeData map[string]interface{}, inputKey []interface{}, curKeyIndex int) (*SlotParserResponse, error) {
	if len(inputKey)-1 < curKeyIndex {
		return nil, fmt.Errorf("insufficient number of keys are passed")
	}
	var numberOfItems uint64
	itemToRetrieve, ok := inputKey[curKeyIndex].(int)
	if !ok {
		return nil, fmt.Errorf("array index must be int type")
	}

	storageTypeBaseParsedData, err := slotParser.GetParsedStorageTypeData(storageTypeData, storageTypeParsedData.Base)
	if err != nil {
		return nil, err
	}
	storageTypeBaseNumberOfBytesBigInt, ok := new(big.Int).SetString(storageTypeBaseParsedData.NumberOfBytes, 10)
	if !ok {
		return nil, fmt.Errorf("failed to convert string to bigInt")
	}
	numToAddForNextSlot := new(big.Int).Div(new(big.Int).Mul(big.NewInt(int64(itemToRetrieve)), storageTypeBaseNumberOfBytesBigInt), big.NewInt(32))
	if storageTypeParsedData.Encoding == InPlaceEncoding {
		numberOfBytesBigInt, ok := new(big.Int).SetString(storageTypeParsedData.NumberOfBytes, 10)
		if !ok {
			return nil, fmt.Errorf("failed to convert string to bigInt")
		}
		numberOfItems = new(big.Int).Div(numberOfBytesBigInt, storageTypeBaseNumberOfBytesBigInt).Uint64()
		parentSlot = common.BigToHash(new(big.Int).Add(parentSlot.Big(), numToAddForNextSlot))
	} else {
		slotData, err := slotParser.getSlotData(parentSlot, 0, 32)
		if err != nil {
			return nil, err
		}
		numberOfItems = new(big.Int).SetBytes(slotData).Uint64()
		parentSlot = common.BigToHash(new(big.Int).Add(crypto.Keccak256Hash(encodeUint256(parentSlot.Big())).Big(), numToAddForNextSlot))
	}
	if itemToRetrieve >= int(numberOfItems) {
		return nil, fmt.Errorf("index number for array must be lesser than array size")
	}
	if IsNonPrimitiveStorageType(storageTypeParsedData.Base) {
		tempStorageTypeParsedData, err := slotParser.GetParsedStorageTypeData(storageTypeData, storageTypeParsedData.Base)
		if err != nil {
			return nil, err
		}
		numberOfBytesOfType, err := strconv.Atoi(tempStorageTypeParsedData.NumberOfBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to convert string numberOfBytes to int, Err: %v", err)
		}
		dataOffset := itemToRetrieve - 1
		if dataOffset < 0 {
			dataOffset = 0
		} else if dataOffset <= 0 {
			dataOffset = numberOfBytesOfType
		} else {
			dataOffset *= itemToRetrieve
		}
		slotData, err := slotParser.getSlotData(parentSlot, dataOffset, numberOfBytesOfType)
		if err != nil {
			return nil, err
		}
		return &SlotParserResponse{
			Slots:       []common.Hash{parentSlot},
			Data:        BytesTo(slotData, IntegerStorageType),
			DataInBytes: slotData,
			DataType:    storageTypeParsedData.Base,
		}, nil
	}
	return slotParser.GetStorageSlotNum(parentSlot, storageTypeParsedData.Base, 0, storageTypeData, inputKey, curKeyIndex+1)
}

func (slotParser *SlotParser) getBytesOrStringTypeSlot(parentSlot common.Hash, storageType string) (*SlotParserResponse, error) {
	slotData, err := slotParser.getSlotData(parentSlot, 0, 32)
	if err != nil {
		return nil, err
	}
	var highLevelStorageType StorageHighLevelTypes
	if storageType[:ArrayAndBytesStorageTypeLen] == BytesTypeStr {
		highLevelStorageType = BytesStorageType
	} else {
		highLevelStorageType = StringStorageType
	}
	slotDataBigInt := new(big.Int).SetBytes(slotData)
	if new(big.Int).Mod(slotDataBigInt, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		return &SlotParserResponse{
			Slots:       []common.Hash{parentSlot},
			Data:        BytesTo(slotData[:31], highLevelStorageType),
			DataInBytes: slotData[:31],
			DataType:    storageType,
		}, nil
	}
	expectedConsumedSlotsWithDecimals := new(big.Float).Quo(new(big.Float).SetInt(new(big.Int).Sub(slotDataBigInt, big.NewInt(1))), big.NewFloat(64))
	expectedConsumedSlotsInt, _ := expectedConsumedSlotsWithDecimals.Int(nil)
	expectedConsumedSlots := new(big.Float).SetInt(expectedConsumedSlotsInt)
	if expectedConsumedSlots.Cmp(expectedConsumedSlotsWithDecimals) != 0 {
		expectedConsumedSlots.Add(expectedConsumedSlots, big.NewFloat(1))
	}
	fmt.Println("Expected Consumed Slots For Bytes Data:", expectedConsumedSlots)
	expectedConsumedSlotsInt64, _ := expectedConsumedSlots.Int64()
	childSlots := []common.Hash{}
	childSlotsData := []byte{}
	lastestBlockNum, err := slotParser.ethClient.BlockNumber(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block number, rrr: %v", err)
	}
	var i int64 = 0
	for ; i < expectedConsumedSlotsInt64; i++ {
		childSlotNum := common.BigToHash(new(big.Int).Add(crypto.Keccak256Hash(encodeUint256(parentSlot.Big())).Big(), big.NewInt(i)))
		childSlots = append(childSlots, childSlotNum)
		slotData, err := slotParser.ethClient.StorageAt(context.Background(), common.HexToAddress(slotParser.contractAddress), childSlotNum, big.NewInt(int64(lastestBlockNum)))
		if err != nil {
			return nil, err
		}
		childSlotsData = append(childSlotsData, slotData...)
	}
	return &SlotParserResponse{
		Slots:       childSlots,
		Data:        BytesTo(childSlotsData, highLevelStorageType),
		DataInBytes: childSlotsData,
		DataType:    storageType,
	}, nil
}

func (slotParser *SlotParser) getStructTypeSlot(parentSlot common.Hash, storageTypeParsedData *StorageLayoutTypeData, storageTypeData map[string]interface{}, inputKey []interface{}, curKeyIndex int) (*SlotParserResponse, error) {
	if len(inputKey)-1 < curKeyIndex {
		return nil, fmt.Errorf("insufficient number of keys are passed")
	}
	var structItemType string
	var offset int
	for _, data := range storageTypeParsedData.Members {
		if data.Label == inputKey[curKeyIndex] {
			childIndex, ok := new(big.Int).SetString(data.Slot, 10)
			if !ok {
				return nil, fmt.Errorf("failed to convert string to bigInt")
			}
			parentSlot = common.BigToHash(new(big.Int).Add(parentSlot.Big(), childIndex))
			structItemType = data.Type
			offset = data.Offset
			break
		}
	}
	if IsNonPrimitiveStorageType(structItemType) {
		tempStorageTypeParsedData, err := slotParser.GetParsedStorageTypeData(storageTypeData, structItemType)
		if err != nil {
			return nil, err
		}
		numberOfBytesOfType, err := strconv.Atoi(tempStorageTypeParsedData.NumberOfBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to convert string numberOfBytes to int, Err: %v", err)
		}
		slotData, err := slotParser.getSlotData(parentSlot, offset, numberOfBytesOfType)
		if err != nil {
			return nil, err
		}
		return &SlotParserResponse{
			Slots:       []common.Hash{parentSlot},
			Data:        BytesTo(slotData, IntegerStorageType),
			DataInBytes: slotData,
			DataType:    structItemType,
		}, nil
	}
	return slotParser.GetStorageSlotNum(parentSlot, structItemType, 0, storageTypeData, inputKey, curKeyIndex+1)
}

func (slotParser *SlotParser) getMapTypeSlot(parentSlot common.Hash, storageTypeParsedData *StorageLayoutTypeData, storageTypeData map[string]interface{}, inputKey []interface{}, curKeyIndex int) (*SlotParserResponse, error) {
	if len(inputKey)-1 < curKeyIndex {
		return nil, fmt.Errorf("insufficient number of keys are passed")
	}
	keyType := strings.Split(storageTypeParsedData.Key, "_")[1]
	var encodedValueForSlot []byte
	if keyType == "string" {
		keyValueStr, ok := inputKey[curKeyIndex].(string)
		if !ok {
			return nil, fmt.Errorf("invalid key type, key type must be string")
		}
		encodedValueForSlot = append([]byte(keyValueStr), encodeUint256(parentSlot.Big())...)
	} else if keyType[:4] == "byte" {
		keyValueHexStr, ok := inputKey[curKeyIndex].(string)
		if !ok {
			return nil, fmt.Errorf("invalid key type, key type must be hex string for representing bytes data")
		}
		keyValueBytes, err := hexutil.Decode(keyValueHexStr)
		encodedValueForSlot = append(keyValueBytes, encodeUint256(parentSlot.Big())...)
		if err != nil {
			return nil, fmt.Errorf("filed to convert hex string to bytes, Err: %v", err)
		}
	} else {
		encodedValueForSlot = encodeTwoTypes(keyType, "uint256", inputKey[curKeyIndex], parentSlot.Big())
	}
	parentSlot = crypto.Keccak256Hash(encodedValueForSlot)
	if IsNonPrimitiveStorageType(storageTypeParsedData.Value) {
		slotData, err := slotParser.getSlotData(parentSlot, 0, 32)
		if err != nil {
			return nil, err
		}
		return &SlotParserResponse{
			Slots:       []common.Hash{parentSlot},
			Data:        BytesTo(slotData, IntegerStorageType),
			DataInBytes: slotData,
			DataType:    storageTypeParsedData.Value,
		}, nil
	}
	return slotParser.GetStorageSlotNum(parentSlot, storageTypeParsedData.Value, 0, storageTypeData, inputKey, curKeyIndex+1)
}

func (slotParser *SlotParser) GetStorageSlotNum(parentSlot common.Hash, storageType string, offset int, storageTypeData map[string]interface{}, inputKey []interface{}, nextKeyIndex int) (*SlotParserResponse, error) {
	storageTypeParsedData, err := slotParser.GetParsedStorageTypeData(storageTypeData, storageType)
	if err != nil {
		return nil, err
	}
	if len(storageType) >= MinStorageTypeLen {
		if storageType[:ArrayAndBytesStorageTypeLen] == ArrayTypeStr {
			return slotParser.getStaticArrayTypeSlot(parentSlot, storageTypeParsedData, storageTypeData, inputKey, nextKeyIndex)
		} else if storageType[:ArrayAndBytesStorageTypeLen] == BytesTypeStr || storageType[:StructAndStringStorageTypeLen] == StringTypeStr {
			return slotParser.getBytesOrStringTypeSlot(parentSlot, storageType)
		} else if storageType[:StructAndStringStorageTypeLen] == StructTypeStr {
			return slotParser.getStructTypeSlot(parentSlot, storageTypeParsedData, storageTypeData, inputKey, nextKeyIndex)
		} else if storageType[:MapStorageTypeLen] == MapTypeStr {
			return slotParser.getMapTypeSlot(parentSlot, storageTypeParsedData, storageTypeData, inputKey, nextKeyIndex)
		}
	}
	numberOfBytes, err := strconv.Atoi(storageTypeParsedData.NumberOfBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert string numberOfBytes to int, Err: %v", err)
	}
	slotData, err := slotParser.getSlotData(parentSlot, offset, numberOfBytes)
	if err != nil {
		return nil, err
	}
	return &SlotParserResponse{
		Slots:       []common.Hash{parentSlot},
		Data:        BytesTo(slotData, IntegerStorageType),
		DataInBytes: slotData,
		DataType:    storageType,
	}, nil
}
