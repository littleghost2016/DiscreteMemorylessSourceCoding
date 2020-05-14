package lempelziv

import (
	"DiscreteMemorylessSourceCoding/util"
	"bytes"
	"fmt"
	"math"
	"os"
)

func EncodeHandler(filePath string, textByteSlice []byte) {

	if len(textByteSlice) == 0 {
		fmt.Println("There is no character in text!")
		os.Exit(1)
	}

	// 注意：decodingDirectorySlice是译码表，所有从文中统计出来的不同段是从1开始索引的，即第0个虽然没有直接写在slice里，
	// 但默认单个字符的SegmentNumber为0，其余SegmentNumber从1开始（本身所在的索引+1）
	// singleCharacterDirectorySlice就是单个字符字典，其CharacterNubmer与本身的索引值相等

	// 注意：此处lastIsSpecialFlag若为true表明，decodingDirectorySlice最后一个元素无法再组成新的段，只能使用前面的segmentNumber，其lastCharacterNumber无效
	// 即segmentNumber起作用，lastCharacterNumber无作用
	decodingDirectorySlice, singleCharacterDirectorySlice, lastIsSpecialFlag := generateDecodingDirectorySliceAndSingleCharacterDirectorySlice(textByteSlice)

	decodingDirectoryLength := calculateBinaryDigitsNumber(len(decodingDirectorySlice))
	singleCharacterDirectoryLength := calculateBinaryDigitsNumber(len(singleCharacterDirectorySlice))

	// fmt.Println(decodingDirectoryLength, singleCharacterDirectoryLength)

	outputByteChannel := make(chan byte, 64)

	// // 准备二进制文件所需的数据
	go WriteBinaryToFile(decodingDirectorySlice, singleCharacterDirectorySlice, decodingDirectoryLength, singleCharacterDirectoryLength, lastIsSpecialFlag, outputByteChannel)

	// 构造输出文件名
	binaryFileName := fmt.Sprintf("%s.bin", filePath)
	fmt.Println("binaryFileName", binaryFileName)
	util.WriteByteToFile(binaryFileName, outputByteChannel)
}

func generateDecodingDirectorySliceAndSingleCharacterDirectorySlice(fileContent []byte) (decodingDirectorySlice []*DecodingDirectoryNode, singleCharacterDirectorySlice []*SingleCharacterDirectoryNode, lastIsSpecialFlag bool) {

	tempDecodingDirectoryCharacterByteSlice := []byte{}
	// tempSingleCharacterDirectoryNode

	var tempSegmentNumber int = 0
	var singleCharacterDirectorySliceIndex uint8 = 0

	lastIsSpecialFlag = false

	for _, eachTextByte := range fileContent {

		// 处理singleCharacterDirectorySlice的生成
		singleCharacterExistFlag := false
		for _, each := range singleCharacterDirectorySlice {
			if eachTextByte == each.Character {
				singleCharacterExistFlag = true
				break
			}
		}
		// 在singleCharacterDirectorySlice里，如果不存在则加入
		if !singleCharacterExistFlag {
			tempSingleCharacterDirectoryNode := SingleCharacterDirectoryNode{
				Type:            1,
				Character:       eachTextByte,
				CharacterNubmer: singleCharacterDirectorySliceIndex,
			}
			singleCharacterDirectorySlice = append(singleCharacterDirectorySlice, &tempSingleCharacterDirectoryNode)
			singleCharacterDirectorySliceIndex++
		}

		// 处理decodingDirectorySlice的生成
		tempDecodingDirectoryCharacterByteSlice = append(tempDecodingDirectoryCharacterByteSlice, eachTextByte)
		characterNodeExistFlag := false
		for index, eachCharacterNode := range decodingDirectorySlice {
			if bytes.Equal(tempDecodingDirectoryCharacterByteSlice, eachCharacterNode.Character) {
				characterNodeExistFlag = true
				// if len(tempSegmentNumber)
				tempSegmentNumber = index + 1
				break
			}
		}

		// 在decodingDirectorySlice里，如果不存在则加入
		if !characterNodeExistFlag {

			if len(tempDecodingDirectoryCharacterByteSlice) == 1 {
				tempSegmentNumber = 0
			}

			var tempLastCharacterNumber uint8

			for index, each := range singleCharacterDirectorySlice {
				if eachTextByte == each.Character {
					tempLastCharacterNumber = uint8(index)
					break
				}
			}

			tempCharacterNode := DecodingDirectoryNode{
				Type:                0,
				Character:           tempDecodingDirectoryCharacterByteSlice,
				SegmentNumber:       tempSegmentNumber,
				LastCharacterNumber: tempLastCharacterNumber,
			}
			decodingDirectorySlice = append(decodingDirectorySlice, &tempCharacterNode)

			tempDecodingDirectoryCharacterByteSlice = []byte{}
		}
	}

	if len(tempDecodingDirectoryCharacterByteSlice) != 0 {
		tempCharacterNode := DecodingDirectoryNode{
			Type:                0,
			Character:           tempDecodingDirectoryCharacterByteSlice,
			SegmentNumber:       tempSegmentNumber,
			LastCharacterNumber: 0,
		}
		decodingDirectorySlice = append(decodingDirectorySlice, &tempCharacterNode)
		lastIsSpecialFlag = true
	}

	return
}

func calculateBinaryDigitsNumber(length int) (binaryDigitsNumber uint8) {

	limit := 1
	binaryDigitsNumber = 1
	for limit < length {
		limit = (limit << 1) + 1
		binaryDigitsNumber++
	}

	return
}

func WriteBinaryToFile(decodingDirectorySlice []*DecodingDirectoryNode, singleCharacterDirectorySlice []*SingleCharacterDirectoryNode, decodingDirectoryLength uint8, singleCharacterDirectoryLength uint8, lastIsSpecialFlag bool, outputByteChannel chan<- byte) {

	// 写入头文件
	// 2 LZ编码
	util.WriteFlag(2, outputByteChannel)

	writeSingleCharacterDirectoryNumber(uint8(len(singleCharacterDirectorySlice)), outputByteChannel)
	writeSingleCharacterDirectory(singleCharacterDirectorySlice, outputByteChannel)
	writeDecodingDirectoryLength(decodingDirectoryLength, outputByteChannel)
	writeSingleCharacterDirectoryLength(singleCharacterDirectoryLength, outputByteChannel)
	calculatedPaddingLength := calculatePaddingLength(len(decodingDirectorySlice), decodingDirectoryLength, singleCharacterDirectoryLength)
	// fmt.Println("calculatedPaddingLength", calculatedPaddingLength)
	writePaddingLength(calculatedPaddingLength, outputByteChannel)
	wirteLastIsSpecialFlag(lastIsSpecialFlag, outputByteChannel)

	writeCode(decodingDirectorySlice, decodingDirectoryLength, singleCharacterDirectoryLength, calculatedPaddingLength, outputByteChannel)

}

func writeCode(decodingDirectorySlice []*DecodingDirectoryNode, decodingDirectoryLength uint8, singleCharacterDirectoryLength uint8, calculatedPaddingLength uint8, outputByteChannel chan<- byte) {

	bitChannel := make(chan bool, int(math.Pow(2, 16)))

	go func() {
		for _, each := range decodingDirectorySlice {

			// fmt.Println(each.Character)
			// 处理segmentNumber
			for i := decodingDirectoryLength; i > 0; i-- {
				if tempBit := each.SegmentNumber & (1 << (i - 1)); tempBit == 0 {
					bitChannel <- false
					// fmt.Print(0)
				} else {
					bitChannel <- true
					// fmt.Print(1)
				}
			}
			// fmt.Print("+")
			// 处理lastCharacterNumber
			for i := singleCharacterDirectoryLength; i > 0; i-- {
				if tempBit := each.LastCharacterNumber & (1 << (i - 1)); tempBit == 0 {
					bitChannel <- false
					// fmt.Print(0)
				} else {
					bitChannel <- true
					// fmt.Print(1)
				}
			}
			// fmt.Println()
			// time.Sleep(1 * time.Second)
		}
		close(bitChannel)
	}()

	operatedPaddingLength := util.ConvertCodeBitToCodeByte(bitChannel, outputByteChannel)
	// fmt.Println("operatedPaddingLength", operatedPaddingLength)
	if calculatedPaddingLength != operatedPaddingLength {
		fmt.Println(calculatedPaddingLength, operatedPaddingLength)
		fmt.Println("calculatePaddingLength != operatedPaddingLength")
	}
}

func writeSingleCharacterDirectoryNumber(singleCharacterDirectoryNumber uint8, outputByteChannel chan<- byte) {
	outputByteChannel <- singleCharacterDirectoryNumber
}

func writeSingleCharacterDirectory(singleCharacterDirectorySlice []*SingleCharacterDirectoryNode, outputByteChannel chan<- byte) {
	for _, each := range singleCharacterDirectorySlice {
		outputByteChannel <- each.Character
	}
}

func writeDecodingDirectoryLength(decodingDirectoryLength uint8, outputByteChannel chan<- byte) {
	temp4byteArray := util.CounvertUint32ToByteSlice(uint32(decodingDirectoryLength))
	for i := 0; i < 4; i++ {
		outputByteChannel <- temp4byteArray[i]
	}
}

func writeSingleCharacterDirectoryLength(singleCharacterDirectoryLength uint8, outputByteChannel chan<- byte) {
	outputByteChannel <- singleCharacterDirectoryLength
}

func calculatePaddingLength(decodingDirectoryNumber int, decodingDirectoryLength uint8, singleCharacterDirectoryLength uint8) (paddingLength uint8) {
	paddingLength = (8 - uint8((decodingDirectoryNumber*int(decodingDirectoryLength+singleCharacterDirectoryLength))%8)) % 8
	return
}

func writePaddingLength(calculatePaddingLength uint8, outputByteChannel chan<- byte) {
	outputByteChannel <- calculatePaddingLength
}

func wirteLastIsSpecialFlag(lastIsSpecialFlag bool, outputByteChannel chan<- byte) {
	if lastIsSpecialFlag {
		outputByteChannel <- 0x01
	} else {
		outputByteChannel <- 0x00
	}
}
