package arithmetic

import (
	"DiscreteMemorylessSourceCoding/util"
	"fmt"
	"os"
	"sort"

	"github.com/shopspring/decimal"
)

// EncodeHandler arithmetic编码处理
func EncodeHandler(filePath string, textByteSlice []byte) {

	if len(textByteSlice) == 0 {
		fmt.Println("There is no character in text!")
		os.Exit(1)
	}

	// 统计字符出现次数
	characterFrequencyMap := util.CountCharacterFromText(textByteSlice)

	var totalFrequency decimal.Decimal
	for _, v := range characterFrequencyMap {
		totalFrequency = totalFrequency.Add(decimal.NewFromInt32(int32(v)))
	}

	fmt.Println("totalFrequency", totalFrequency)

	characterMap := make(map[byte]*CharacterNode)

	var characterNodeSlice []*CharacterNode

	for k, v := range characterFrequencyMap {
		frequencyBigint := decimal.NewFromInt32(int32(v))
		weight := frequencyBigint.Div(totalFrequency)
		tempCharacterNode := CharacterNode{
			Character: k,
			Frequency: v,
			Weight:    weight,
		}

		if _, ok := characterMap[k]; !ok {
			characterMap[k] = &tempCharacterNode
		}

		characterNodeSlice = append(characterNodeSlice, &tempCharacterNode)
	}

	// for k, v := range characterMap {
	// 	fmt.Println(k, v)
	// }
	// fmt.Println("111")
	sort.Sort(CharacterNodeSlice(characterNodeSlice))
	// 分配每个字符所在的区间端点
	tempLeftBounded := decimal.Zero
	for _, eachCharacterNode := range characterNodeSlice {
		eachCharacterNode.LeftBounded = tempLeftBounded
		tempLeftBounded = tempLeftBounded.Add(eachCharacterNode.Weight)
		eachCharacterNode.RightBounded = tempLeftBounded
	}

	// for _, each := range characterNodeSlice {
	// 	fmt.Println(each)
	// }

	byteChannelToFile := make(chan byte, 64)

	go writeBinaryFile(textByteSlice, characterMap, characterNodeSlice, byteChannelToFile)

	// 构造输出文件名
	binaryFileName := fmt.Sprintf("%s.bin", filePath)
	fmt.Println("binaryFileName", binaryFileName)
	util.WriteByteToFile(binaryFileName, byteChannelToFile)
}

func writeBinaryFile(textByteSlice []byte, characterMap map[byte]*CharacterNode, characterNodeSlice []*CharacterNode, byteChannelToFile chan<- byte) {

	// 1: arithmetic
	util.WriteFlag(1, byteChannelToFile)

	writeCodeNumber(characterMap, byteChannelToFile)
	writeCodeMap(characterNodeSlice, byteChannelToFile)
	writeTotalTextByteNumber(textByteSlice, byteChannelToFile)
	writeCode(textByteSlice, characterMap, byteChannelToFile)

	close(byteChannelToFile)
}

func writeCodeNumber(characterMap map[byte]*CharacterNode, byteChannelToFile chan<- byte) {
	codeNumber := uint8(len(characterMap)) // 类型转换原因同上
	fmt.Println("codeNumber", codeNumber)

	byteChannelToFile <- codeNumber
}

func writeCodeMap(characterNodeSlice []*CharacterNode, byteChannelToFile chan<- byte) {
	for _, eachNode := range characterNodeSlice {
		byteChannelToFile <- eachNode.Character
		temp4byteArray := util.CounvertUint32ToByteSlice(eachNode.Frequency)
		for i := 0; i < 4; i++ {
			byteChannelToFile <- temp4byteArray[i]
		}
	}
}

func writeTotalTextByteNumber(textByteSlice []byte, byteChannelToFile chan<- byte) {
	temp4byteArray := util.CounvertUint32ToByteSlice(uint32(len(textByteSlice)))
	for i := 0; i < 4; i++ {
		byteChannelToFile <- temp4byteArray[i]
	}
}

func writeCode(textByteSlice []byte, characterMap map[byte]*CharacterNode, byteChannelToFile chan<- byte) {

	var finalCodeByteSlice []byte

	low := decimal.Zero
	high := decimal.NewFromInt(1)

	zeroDecimal := decimal.Zero
	oneDecimal := decimal.NewFromInt(1)

	lowChannel := make(chan decimal.Decimal)
	highChannel := make(chan decimal.Decimal)

	for _, eachByte := range textByteSlice {

		// fmt.Println("low", low)
		// fmt.Println("high", high)
		go func() {
			lowChannel <- low.Add(characterMap[eachByte].LeftBounded.Mul(high.Sub(low)))
		}()
		go func() {
			highChannel <- low.Add(characterMap[eachByte].RightBounded.Mul(high.Sub(low)))
		}()

		low = <-lowChannel
		high = <-highChannel

		lowByteSlice, _ := low.MarshalText()
		highByteSlice, _ := high.MarshalText()

		lowByteSliceLength := len(lowByteSlice)
		highByteSliceLength := len(highByteSlice)

		truncatePosition := 0
		// FIXME: decimal变[]byte 可能不只是简单地直接变，还要再研究一下对应规则
		if lowByteSliceLength <= highByteSliceLength {
			for i := 2; i < lowByteSliceLength; i++ {
				if lowByteSlice[i] == highByteSlice[i] {
					finalCodeByteSlice = append(finalCodeByteSlice, highByteSlice[i])
				} else {
					truncatePosition = i - 2
					break
				}
			}

		} else {
			for i := 2; i < highByteSliceLength; i++ {
				if highByteSlice[i] == lowByteSlice[i] {
					finalCodeByteSlice = append(finalCodeByteSlice, highByteSlice[i])
				} else {
					truncatePosition = i - 2
					break
				}
			}
		}
		// fmt.Println(truncatePosition)
		if (!low.Equal(zeroDecimal)) || (!high.Equal(oneDecimal)) {
			lowShifted := low.Shift(int32(truncatePosition))
			// 	// // fmt.Println("lowShifted", lowShifted)
			low = lowShifted.Sub(lowShifted.Floor()).Truncate(10)
			// 	// low = low.Truncate(20)
			highShifted := high.Shift(int32(truncatePosition))
			// 	// // fmt.Println("highShifted", highShifted)
			high = highShifted.Sub(highShifted.Floor()).Truncate(10)
			// 	// high = high.Truncate(20)
		}
	}
	lastByte, _ := high.MarshalText()
	finalCodeByteSlice = append(finalCodeByteSlice, lastByte[2])

	// TODO:
	// fmt.Println(len(finalCodeByteSlice), finalCodeByteSlice)

	convertDigitalToByte(finalCodeByteSlice, byteChannelToFile)
}

func convertDigitalToByte(finalCodeByteSlice []byte, byteChannelToFile chan<- byte) {

	// semiByte := 0x00
	semiByteChannel := make(chan byte, 16)
	go func() {
		for _, eachByte := range finalCodeByteSlice {
			semiByteChannel <- eachByte
		}
		close(semiByteChannel)
	}()

	receivedFirstSemiByte, ok1 := <-semiByteChannel
	receivedSecondSemiByte, ok2 := <-semiByteChannel

	for {
		if ok1 {
			if ok2 {
				tempFirstSemiByte := receivedFirstSemiByte % 48
				tempSecondSemiByte := receivedSecondSemiByte % 48
				tempByte := ((0x00 | tempFirstSemiByte) << 4) | tempSecondSemiByte
				byteChannelToFile <- tempByte
				receivedFirstSemiByte, ok1 = <-semiByteChannel
				receivedSecondSemiByte, ok2 = <-semiByteChannel
			} else {
				tempFirstSemiByte := receivedFirstSemiByte % 48
				var tempSecondSemiByte byte
				tempSecondSemiByte = 0x0a
				tempByte := ((0x00 | tempFirstSemiByte) << 4) | tempSecondSemiByte
				byteChannelToFile <- tempByte
				receivedFirstSemiByte, ok1 = <-semiByteChannel
				receivedSecondSemiByte, ok2 = <-semiByteChannel
				break
			}
		} else {
			// fmt.Println("none")

			break
		}
	}
}
