package slot_parser

const (
	ArrayAndBytesStorageTypeLen   = 7
	StructAndStringStorageTypeLen = 8
	MapStorageTypeLen             = 9
	MinStorageTypeLen             = 9
	StringTypeStr                 = "t_string"
	BytesTypeStr                  = "t_bytes"
	ArrayTypeStr                  = "t_array"
	MapTypeStr                    = "t_mapping"
	StructTypeStr                 = "t_struct"
	InPlaceEncoding               = "inplace"
	DynamicArrayEncoding          = "dynamic_array"
	StorageLayoutDirName          = "/DirForSolc"
	StorageLayoutFilePaddingStr   = "_storage.json"
)
const (
	IntegerStorageType StorageHighLevelTypes = "Integer"
	StringStorageType  StorageHighLevelTypes = "String"
	BytesStorageType   StorageHighLevelTypes = "Bytes"
)
