// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.8;

// Note: This Contract Is Not Well Written And Not Follows Standard, Its Only Written For Testing Of This Package
// I Have Tried To Test All Type Of Storage Slot Using This Contract.

contract RequestReceiver {
    struct Message {
        uint256 id;
        bytes data;
        bytes metaData;
        Status status;
    }

    enum Status {
        Requested,
        Accepted
    }

    struct Message1 {
        uint256 id;
        bytes data;
        mapping(uint256 => bytes) metaData;
    }
    event ReceivedData(uint256 data);
    uint8 public temp8_1 = 10;
    address public tempAddress_1 = 0x90a23757BabC1a2823D1f085A7f3f3d426E751d6;
    uint128 temp128_1 = 10;
    uint128 temp128_2 = 11;
    Status public status_1 = Status.Accepted;
    uint256 temp256_1 = 11;
    string tempStr_1 = "111111111111111111111111111";
    string tempStr_2 = "111111111111111111111111111111111111111111111111111111";
    bytes public tempBytes1 = abi.encode(40, "node1");
    uint256 public receivedData = 1;
    uint256 public proxyReceivedData = 2;
    bytes public tempBytes_1 = hex"0a";
    bytes public tempBytes_2 =
        hex"010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101";
    bytes public tempBytes_3 =
        hex"000000000000000000000000000000000000000000000000000000000001c9";
    Message public tempMessage =
        Message({
            id: 1,
            data: "0x01",
            metaData: "0x02",
            status: Status.Accepted
        });
    uint256[2] public tempData = [1, 2];
    string[2] public tempData_1 = [
        "1",
        "7777777777777777777777777777777777777777777777777777777"
    ];
    string[] public tempData_2 = [
        "1",
        "7777777777777777777777777777777777777777777777777777777"
    ];
    uint256[] public tempArr_1 = [11, 22];
    uint128[2] public tempArr_2 = [11, 22];
    uint8[] public tempArr_3 = [11, 22];
    Message[2] public mess_1;
    Message[] public mess_2;
    bytes16 public tempByte_1;
    bytes16 public tempByte_2;
    mapping(uint256 => uint256) private temp;
    mapping(uint256 => uint8) private temp_1;
    mapping(uint256 => Message) private tempMess;
    mapping(uint256 => bytes) private tempMess_1;
    mapping(uint256 => mapping(uint256 => uint256)) private tempMess_2;
    mapping(uint256 => mapping(uint256 => Message1)) private tempMess_3;
    mapping(string => mapping(bytes => uint256)) public tempMess_4;
    mapping(bytes => uint256) public tempMess_5;

    constructor() {
        temp[1] = 1;
        temp_1[1] = 11;
        temp_1[2] = 12;
        mess_1[0] = Message({
            id: 1,
            data: "0x01",
            metaData: "0x02",
            status: Status.Accepted
        });
        mess_1[1] = Message({
            id: 1,
            data: "0x09",
            metaData: "0x05",
            status: Status.Accepted
        });
        mess_2.push(
            Message({
                id: 1,
                data: "0x01",
                metaData: "0x02",
                status: Status.Accepted
            })
        );
        mess_2.push(
            Message({
                id: 1,
                data: "0x06",
                metaData: "0x07",
                status: Status.Accepted
            })
        );
        tempMess[1] = Message({
            id: 1,
            data: "0x01",
            metaData: "0x02",
            status: Status.Accepted
        });
        tempMess_1[1] = hex"0a";
        tempMess_2[1][1] = 10;
        tempMess_3[1][1].id = 100;
        tempMess_3[1][1].data = hex"0a0a";
        tempMess_3[1][1].metaData[1] = hex"0a0a";
        tempMess_4["node1"][hex"0a"] = 10;
        tempMess_5[hex"0a"]=10101010;
    }

    function getSlotNumForMap(
        uint256 key,
        uint256 slotIndex
    ) external pure returns (bytes32) {
        return keccak256(abi.encode(key, slotIndex));
    }

    function getSlotNumForMapToStruct(
        uint256 key,
        uint256 slotIndex,
        uint256 structItemindex
    ) external pure returns (bytes32) {
        return
            bytes32(
                uint256(keccak256(abi.encode(key, slotIndex))) + structItemindex
            );
    }

    // Using This Method We Can Get Slot Num For Dynamic Array Type Variables,
    // Example: string, bytes, []int256
    // Note: This Method Is Only Useful For Getting Slot Num For String And Bytes Type When Size Of Data Will Be More Than 32 Bytes,
    // To Check Size Of Data We Can Just Query Variable Slot Num.
    function getSlotNumForDynamicArrayType(
        uint256 variableSlotNum,
        uint256 index
    ) public pure returns (bytes32) {
        return bytes32(uint256(keccak256(abi.encode(variableSlotNum))) + index);
    }

    receive() external payable {}
}
