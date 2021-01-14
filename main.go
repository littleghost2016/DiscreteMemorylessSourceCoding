package main

import (
	"DiscreteMemorylessSourceCoding/arithmetic"
	"DiscreteMemorylessSourceCoding/huffman"
	"DiscreteMemorylessSourceCoding/lempelziv"
	"DiscreteMemorylessSourceCoding/util"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {

	var filePath string
	// os.Args[0]是二进制文件自身的路径
	// os.Args[1]是拖拽到二进制文件上 文件的路径
	filePath = os.Args[1]

	// 以下几行仅做测试用
	// filePath = "files/test7.txt.1"
	// filePath = "files/test7.txt.1.bin"
	// filePath = "files/辰东-完美世界.txt.1"
	// filePath = "files/辰东-完美世界.txt.2.bin"
	// filePath = "files/共产党宣言.txt.1"
	// filePath = "files/共产党宣言.txt.1.bin"
	// filePath = "files/1.mp3.2.bin"

	startTime := time.Now()

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
		case 1:
			arithmetic.DecodeHandler(filePath, fileContent[3:])
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

		originFile, _ := os.Stat(filePath)
		originFileSize := originFile.Size()

		compressedFileName := fmt.Sprintf("%s.bin", filePath)
		compressedFile, _ := os.Stat(compressedFileName)
		compressedFileSize := compressedFile.Size()

		// fmt.Printf("压缩前文件大小 %d\n", originFileSize)
		fmt.Println(originFileSize)
		// fmt.Printf("压缩后文件大小 %d\n", compressedFileSize)
		fmt.Println(compressedFileSize)
		compressRate := float32(compressedFileSize) / float32(originFileSize)
		// fmt.Printf("压缩后/压缩前=%.2f%%\n", compressRate*100)
		fmt.Printf("%.2f%%\n", compressRate*100)
	}

	stopTime := time.Now()
	duration := stopTime.Sub(startTime)
	// fmt.Printf("开始时间\t%v\n结束时间\t%v\n总耗时\t%v\n", startTime.Format("20060102-15:04:05"), stopTime.Format("20060102-15:04:05"), duration)
	fmt.Println(duration)

	var tempInt int
	fmt.Scanf("%d", &tempInt)
}
