package main

import (
	"DiscreteMemorylessSourceCoding/huffman"
	"DiscreteMemorylessSourceCoding/util"
	"os"
	"strings"
)

func main() {

	var filePath string
	// os.Args[0]是二进制文件自身的路径
	// os.Args[1]是拖拽到二进制文件上 文件的路径
	filePath = os.Args[1]

	// 以下两行做测试用
	// filePath = "files/test2.txt.0"
	// filePath = "files/test2.txt.0.bin"

	fileContent := util.ReadFromFile(filePath)

	// 检测是否0x19 0x15开头
	flagExist := util.ReadFlag(fileContent[0:2])
	// fmt.Println(flagExist)

	if flagExist {
		CodeType := util.ReadCodeType(fileContent[2])
		filePath = filePath[0 : len(filePath)-4]
		switch CodeType {
		case 0:
			huffman.DecodeHandler(filePath, fileContent[3:])
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

	// var tempInt int
	// fmt.Scanf("%d", &tempInt)
}
