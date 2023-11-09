package slot_parser

import (
	"fmt"
	"testing"
)

func TestTemp(t *testing.T) {
	// slotParser, err := NewSlotParer("WETH9", "./contracts/WETH_TestContract.sol", "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2", "https://ethereum.publicnode.com")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(slotParser.GetSlotNum([]interface{}{"name"}))
	// fmt.Println(slotParser.GetSlotNum([]interface{}{"symbol"}))
	// fmt.Println(slotParser.GetSlotNum([]interface{}{"decimals"}))
	// fmt.Println(slotParser.GetSlotNum([]interface{}{"balanceOf", "0xF04a5cC80B1E94C69B48f5ee68a08CD2F09A7c3E"}))
	// fmt.Println(slotParser.GetSlotNum([]interface{}{"allowance", "0xaf1f6b64346750f5AA006CEeE749Be8b9a595303", "0x216B4B4Ba9F3e719726886d34a177484278Bfcae"}))

	slotParser, err := NewSlotParer("Temp", "./contracts/temp.sol", "0xF414FC31909D97fF3239229fD0d239cd6883AEF0", "https://ethereum-goerli.publicnode.com")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(slotParser.GetSlotNum([]interface{}{"temp8_1"}))
	fmt.Println(slotParser.GetSlotNum([]interface{}{"tempMess_4", "node1", "0x0a"}))
}
