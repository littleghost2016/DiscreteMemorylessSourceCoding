package lempelziv

import (
	"DiscreteMemorylessSourceCoding/util"
	"fmt"
	"math"
	"os"
)

func EncodeHandler(filePath string, textByteSlice []byte) {

	if len(textByteSlice) == 0 {
		fmt.Println("There is no character in text!")
		os.Exit(1)
	}

	// decodingDirectorySlice, singleCharacterDirectorySlice, decodingDirectory, singleCharacterDirectory, lastIsSpecialFlag := generateDecodingDirectoryAndSingleCharacterDirectory(textByteSlice)
	decodingDirectorySlice, singleCharacterDirectorySlice, _, _, lastIsSpecialFlag := generateDecodingDirectoryAndSingleCharacterDirectory(textByteSlice)

	// for _, each := range decodingDirectorySlice {
	// 	fmt.Println(1, each)
	// }
	// for _, each := range singleCharacterDirectorySlice {
	// 	fmt.Println(2, each)
	// }
	// for _, each := range decodingDirectory {
	// 	// fmt.Println(31, len(decodingDirectory))
	// 	fmt.Println(3, each)
	// }
	// for _, each := range singleCharacterDirectory {
	// 	fmt.Println(4, each)
	// }
	// fmt.Println(lastIsSpecialFlag)

	decodingDirectoryLength := calculateBinaryDigitsNumber(len(decodingDirectorySlice))
	singleCharacterDirectoryLength := calculateBinaryDigitsNumber(len(singleCharacterDirectorySlice))

	// fmt.Println(decodingDirectoryLength, singleCharacterDirectoryLength)

	outputByteChannel := make(chan byte, 64)

	// 准备二进制文件所需的数据
	go WriteBinaryToFile(decodingDirectorySlice, singleCharacterDirectorySlice, decodingDirectoryLength, singleCharacterDirectoryLength, lastIsSpecialFlag, outputByteChannel)

	// 构造输出文件名
	binaryFileName := fmt.Sprintf("%s.bin", filePath)
	fmt.Println("binaryFileName", binaryFileName)
	util.WriteByteToFile(binaryFileName, outputByteChannel)
}

func generateDecodingDirectoryAndSingleCharacterDirectory(fileContent []byte) (decodingDirectorySlice []*DecodingDirectoryNode, singleCharacterDirectorySlice []*SingleCharacterDirectoryNode, decodingDirectory map[string]*DecodingDirectoryNode, singleCharacterDirectory map[byte]*SingleCharacterDirectoryNode, lastIsSpecialFlag bool) {

	decodingDirectory = make(map[string]*DecodingDirectoryNode)
	singleCharacterDirectory = make(map[byte]*SingleCharacterDirectoryNode)

	tempDecodingDirectoryCharacterByteSlice := []byte{}
	// tempSingleCharacterDirectoryNode

	var tempSegmentNumber int = 0
	var decodingDirectorySliceIndex int = 0
	var singleCharacterDirectorySliceIndex uint8 = 0

	lastIsSpecialFlag = false

	for _, eachTextByte := range fileContent {

		// 在singleCharacterDirectorySlice里，如果不存在则加入
		if _, ok := singleCharacterDirectory[eachTextByte]; !ok {
			tempSingleCharacterDirectoryNode := SingleCharacterDirectoryNode{
				Type:            1,
				Character:       eachTextByte,
				CharacterNubmer: singleCharacterDirectorySliceIndex,
			}
			singleCharacterDirectory[eachTextByte] = &tempSingleCharacterDirectoryNode
			singleCharacterDirectorySlice = append(singleCharacterDirectorySlice, &tempSingleCharacterDirectoryNode)
			singleCharacterDirectorySliceIndex++
		}

		// 处理decodingDirectorySlice的生成
		tempDecodingDirectoryCharacterByteSliceWithoutThelastOne := tempDecodingDirectoryCharacterByteSlice[:]
		tempDecodingDirectoryCharacterByteSlice = append(tempDecodingDirectoryCharacterByteSlice, eachTextByte)

		tempKey := fmt.Sprintf("%s", tempDecodingDirectoryCharacterByteSlice)
		// fmt.Printf("%s, %T\n", tempDecodingDirectoryCharacterByteSlice, tempDecodingDirectoryCharacterByteSlice)
		// for _, each := range decodingDirectorySlice {
		// 	fmt.Println(each)
		// }
		// fmt.Println()

		// _, ok := decodingDirectory[tempKey]
		// fmt.Println(ok, decodingDirectory)
		if _, ok := decodingDirectory[tempKey]; !ok {

			// 在decodingDirectorySlice里，如果不存在则加入
			if len(tempDecodingDirectoryCharacterByteSlice) == 1 {
				tempSegmentNumber = 0
			} else {
				// fmt.Println(tempDecodingDirectoryCharacterByteSlice[:len(tempDecodingDirectoryCharacterByteSlice)-1])
				tempKey1 := fmt.Sprintf("%s", tempDecodingDirectoryCharacterByteSliceWithoutThelastOne)
				// fmt.Println("---", tempKey)
				tempSegmentNumber = decodingDirectory[tempKey1].SelfSegmentNubmer
			}

			tempLastCharacterNumber := singleCharacterDirectory[eachTextByte].CharacterNubmer

			tempCharacterNode := DecodingDirectoryNode{
				Type:                0,
				Character:           tempDecodingDirectoryCharacterByteSlice,
				SelfSegmentNubmer:   decodingDirectorySliceIndex + 1,
				SegmentNumber:       tempSegmentNumber,
				LastCharacterNumber: tempLastCharacterNumber,
			}
			// fmt.Println(tempKey)
			decodingDirectory[tempKey] = &tempCharacterNode
			decodingDirectorySlice = append(decodingDirectorySlice, &tempCharacterNode)
			decodingDirectorySliceIndex++
			// time.Sleep(time.Second)
			// time.Sleep(1000 * time.Microsecond)
			tempDecodingDirectoryCharacterByteSlice = []byte{}
		}
	}

	if len(tempDecodingDirectoryCharacterByteSlice) != 0 {

		tempKey := fmt.Sprintf("%s", tempDecodingDirectoryCharacterByteSlice)
		// fmt.Println("---", tempKey)
		tempSegmentNumber = decodingDirectory[tempKey].SelfSegmentNubmer

		tempCharacterNode := DecodingDirectoryNode{
			Type:                0,
			Character:           tempDecodingDirectoryCharacterByteSlice,
			SelfSegmentNubmer:   decodingDirectorySliceIndex + 1,
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
