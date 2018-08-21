package dumpfile

import "fmt"

type DumpFile struct {
	Magic   string
	Version uint64
	Dbs     []*RedisDb
}

type RedisDb struct {
	Num uint32
	KVPaires []*KVPair
}

/**
Value :Type
0 = “String Encoding”
1 = “List Encoding”
2 = “Set Encoding”
3 = “Sorted Set Encoding”
4 = “Hash Encoding”
9 = “Zipmap Encoding”
10 = “Ziplist Encoding”
11 = “Intset Encoding”
12 = “Sorted Set in Ziplist Encoding”
13 = “Hashmap in Ziplist Encoding” (Introduced in rdb version 4)
 */
type KVPair struct {
	// 1:Second, 2: ms 0:dont have expire time
	ExpireType int32
	Expire     uint32
	ValueType  uint32
	Key        string
	value      []byte
}

func ParseDump(reader *FileReader)*DumpFile  {
	//fmt.Print(reader.ReadBytes(5))
	magic:=reader.ReadString(5)
	if "REDIS" != magic{
		panic(reader.fileName()+" is not a standard dump file.")
	}
	dump:=&DumpFile{Magic:magic}
	dump.Version=reader.ReadStringToUInt64(4)
	dbs:=make([]*RedisDb,0)
	dbs=append(dbs,readDb(reader))

	return dump
}

func readDb(reader *FileReader)*RedisDb {
	flag1:=reader.ReadBytes(100)
	fmt.Print(flag1)

	flag:=reader.ReadByte()
	fmt.Print(flag)
	fmt.Print(flag)
	if 0xFE!=flag{
		panic(" it is not db flag")
	}
	db:=&RedisDb{}
	db.Num=reader.ReadUInt32(1)
	return db
}