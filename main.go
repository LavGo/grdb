package main

import (
	"github.com/LavGo/grdb/dumpfile"
	"fmt"
)

func main()  {
	reader:=dumpfile.ReadFile("dump.rdb")
	dump:=dumpfile.ParseDump(reader)
	fmt.Println(dump)
}
