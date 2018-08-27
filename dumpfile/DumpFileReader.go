package dumpfile

import (
	"os"
	"fmt"
	"encoding/binary"
	"strconv"
)

type FileReader struct {
	filename string
	file     *os.File
}

func ReadFile(filename string) *FileReader {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	return &FileReader{
		filename: filename,
		file:     file,
	}
}

func (self *FileReader) fileName() string {
	return self.filename
}
func (self *FileReader) ReadByte() byte {
	buf := make([]byte, 1)
	self.file.Read(buf)
	return buf[0]
}

func (self *FileReader) ReadBytes(size int32) []byte {
	buf := make([]byte, size)
	self.file.Read(buf)
	return buf;
}

func (self *FileReader) ReadString(size int32) string {
	buf := self.ReadBytes(size)
	return string(buf)
}

func (self *FileReader) ReadUInt32(size int32) uint32 {
	cache := make([]byte, 4)
	buf := self.ReadBytes(size)
	copy(cache[4-size:], buf)
	return binary.BigEndian.Uint32(cache)
}
func (self *FileReader) ReadStringToUInt64(size int32) uint64 {
	buf := self.ReadBytes(size)

	v, err := strconv.ParseUint(string(buf), 10, 64)
	if err != nil {
		panic(err)
	}
	return v
}

func (self *FileReader) ReadUInt64(size int32) uint64 {
	buf := self.ReadBytes(size)
	return binary.BigEndian.Uint64(buf)
}

/**
0x80 0x40 0x3F
如果以"00"开头，那么接下来的6个bit表示长度；
如果以“01”开头，那么接下来的14个bit表示长度；
如果以"10"开头，该byte的剩余6bit废弃，接着读入4个bytes表示长度(BigEndian)；
如果以"11"开头，那么接下来的6个bit表示特殊的编码格式，一般用来存储数字：
	0表示用接下来的1byte表示长度
	1表示用接下来的2bytes表示长度；
	2表示用接下来的4bytes表示长度
 */
func (self *FileReader) ReadLength() uint32 {
	buf := self.ReadByte()
	if (0x80&buf == 0) && (0x40&buf == 0) {
		return uint32(buf & 0X3F)
	}
	if (0x80&buf == 0) && (0x40&buf == 1) {
		buf2 := self.ReadByte()
		return uint32(buf&0X3F)<<8 | uint32(buf2)
	}

	if (0x80&buf == 1) && (0x40&buf == 0) {
		buf2 := self.ReadBytes(4)
		return binary.BigEndian.Uint32(buf2)
	}
	if (0x80&buf == 1) && (0x40&buf == 1) {
		if buf&0x3F == 0x00 {
			buf2 := self.ReadByte()
			return uint32(buf2)
		}
		if buf&0x3F == 0x01 {
			buf2 := self.ReadBytes(2)
			// Big Endian
			return uint32(buf2[0])<<8 | uint32(buf2[1])
		}
		if buf&0x3F == 0x02 {
			buf2 := self.ReadBytes(4)
			return binary.BigEndian.Uint32(buf2)
		}
	}
	return 0
}
