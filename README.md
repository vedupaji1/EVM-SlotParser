
# EVM-SlotParser

**TL;DR:** This package can help you to read any type of data of smart contract even private data.

  
This pakage is developed to parse and get data of smart contract storage slot.  
  
In smart contract data is stored in form of slots, every slot have unique id or number using which we can access contant of any slot.

By using this we can read data of smart contract at low level and we can access data of variables whose visibility is private or internal.

Devs can use this package to deeply analyse their contract storage layout.

## Note 

- You will need smart contract and address of contract whose data you want to read.
- Solc compiler must be installed in your system.
- RPC url will be required, you can use any public url or use RPC service of Custom RPC Provider such as Alchemy and Infura.
- To access data of any slot or variable you need to pass their name or required keys, such as if you want to access content of any ```mapping(uint256 => uint256) balance;``` so  here you need to first pass name of this map and then key of map, you have to follow this same structure to access any varibale data or slot data.
- You can get list of available public RPC from [here](https://ethereumnodes.com).

## Example

```
	slotParser, err := NewSlotParer("WETH9", "./contracts/WETH_TestContract.sol", "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2", "https://eth.llamarpc.com")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(slotParser.GetSlotNum([]interface{}{"name"}))
	fmt.Println(slotParser.GetSlotNum([]interface{}{"symbol"}))
	fmt.Println(slotParser.GetSlotNum([]interface{}{"decimals"}))
	fmt.Println(slotParser.GetSlotNum([]interface{}{"balanceOf", "0xF04a5cC80B1E94C69B48f5ee68a08CD2F09A7c3E"}))
	fmt.Println(slotParser.GetSlotNum([]interface{}{"allowance", "0xaf1f6b64346750f5AA006CEeE749Be8b9a595303", "0x216B4B4Ba9F3e719726886d34a177484278Bfcae"}))

```







