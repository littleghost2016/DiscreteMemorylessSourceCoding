package main

import (
	"DiscreteMemorylessSourceCoding/arithmetic"
	"DiscreteMemorylessSourceCoding/huffman"
	"DiscreteMemorylessSourceCoding/lempelziv"
	"DiscreteMemorylessSourceCoding/util"
	"os"
	"strings"
)

func main() {

	var filePath string
	// os.Args[0]是二进制文件自身的路径
	// os.Args[1]是拖拽到二进制文件上 文件的路径
	filePath = os.Args[1]

	// 以下几行仅做测试用
	// filePath = "files/test4.txt.2"
	// filePath = "files/test4.txt.2.bin"
	// filePath = "files/辰东-完美世界.txt.2"
	// filePath = "files/辰东-完美世界.txt.2.bin"
	// filePath = "files/共产党宣言.txt.2"
	// filePath = "files/共产党宣言.txt.2.bin"

	fileContent := util.ReadFromFile(filePath)

	// 检测是否0x19 0x15开头
	flagExist := util.ReadFlag(fileContent[0:2])
	// fmt.Println(flagExist)

	if flagExist {
		CodeType := util.ReadCodeType(fileContent[2])
		filePath = filePath[0 : len(filePath)-4]
		switch CodeType {
		// 霍夫曼编码
		case 0:
			huffman.DecodeHandler(filePath, fileContent[3:])
		// 算术编码
		// case 1:
		// 	arithmetic.DecodeHandler(filePath, fileContent[3:])
		// LZ编码
		case 2:
			lempelziv.DecodeHandler(filePath, fileContent[3:])
		}
	} else {
		filePathStringSlice := strings.Split(filePath, ".")
		encodeType := filePathStringSlice[len(filePathStringSlice)-1]
		switch encodeType {
		// 霍夫曼编码
		case "0":
			huffman.EncodeHandler(filePath, fileContent)
		// 算术编码
		case "1":
			arithmetic.EncodeHandler(filePath, fileContent)
		// LZ编码
		case "2":
			lempelziv.EncodeHandler(filePath, fileContent)
		}
	}

	// var tempInt int
	// fmt.Scanf("%d", &tempInt)
}

// TODO:
// 1. 增加压缩比的计算
// 2. 测试多种文件格式
