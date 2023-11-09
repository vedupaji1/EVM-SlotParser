package slot_parser

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	viperPKG "github.com/spf13/viper"
)

type SlotParser struct {
	ethClient       *ethclient.Client
	contractAddress string
	viper           *viperPKG.Viper
}

type StorageLayoutInfo struct {
	AstId    int    `mapstructure:"astId"`
	Contract string `mapstructure:"contract"`
	Label    string `mapstructure:"label"`
	Offset   int    `mapstructure:"offset"`
	Slot     string `mapstructure:"slot"`
	Type     string `mapstructure:"type"`
}

type StorageLayoutTypeData struct {
	Encoding      string              `mapstructure:"encoding"`
	Label         string              `mapstructure:"label"`
	NumberOfBytes string              `mapstructure:"numberOfBytes"`
	Key           string              `mapstructure:"key"`
	Value         string              `mapstructure:"value"`
	Base          string              `mapstructure:"base"`
	Members       []StorageLayoutInfo `mapstructure:"members"`
}

type SlotParserResponse struct {
	Slots       []common.Hash
	Data        interface{}
	DataInBytes []byte
	DataType    string
}

type StorageHighLevelTypes string
