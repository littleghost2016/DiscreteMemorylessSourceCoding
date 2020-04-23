package main

import (
	"DiscreteMemorylessSourceCoding/huffman"
	"DiscreteMemorylessSourceCoding/util"
	"fmt"
	"strings"
)

func main() {

	var filePath string
	// filePath = os.Args[1]
	filePath = "files/test2.txt.0"
	// filePath = "test4.bin"

	fileContent := util.ReadFromFile(filePath)
	flagExist := util.ReadFlag(fileContent[0:2])

	// 检测是否0x19 0x15开头
	// fmt.Println(flagExist)

	if flagExist {
		CodeType := util.ReadCodeType(fileContent[2])
		switch CodeType {
		case 0:
			huffman.DecodeHandler(fileContent[3:])
			// case 1:
			// 	///
		}
	} else {
		filePathStringSlice := strings.Split(filePath, ".")
		encodeType := filePathStringSlice[len(filePathStringSlice)-1]
		switch encodeType {
		case "0":
			huffman.EncodeHandler(filePath, fileContent)
		}
	}

	var tempInt int
	fmt.Scanf("%d", &tempInt)
}
