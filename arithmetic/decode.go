package arithmetic

import (
	"DiscreteMemorylessSourceCoding/util"
	"fmt"
	"sort"

	"github.com/shopspring/decimal"
)

func DecodeHandler(filePath string, fileContent []byte) {
	byteChannelFromBinaryFile := make(chan byte, 1024)

	// 读取文件
	go func() {
		for _, eachByte := range fileContent {
			byteChannelFromBinaryFile <- eachByte
		}
		close(byteChannelFromBinaryFile)
	}()

	codeNumber := util.ReadCodeNumber(byteChannelFromBinaryFile)
	characterMap := readCodeFrequencyAndGenerateCharacterMap(codeNumber, byteChannelFromBinaryFile)
	totalTextByteNumber := int(readTotalTextByteNumber(byteChannelFromBinaryFile))
	// fmt.Println(characterMap)
	// fmt.Println(totalTextByteNumber) // 5

	finalCodeByteSlice := readCode(byteChannelFromBinaryFile)

	// TODO:
	// fmt.Println(len(finalCodeByteSlice), finalCodeByteSlice)

	var finalCodeNumber decimal.Decimal
	finalCodeNumber.UnmarshalText(finalCodeByteSlice)

	byteChannelToTextFile := make(chan byte, 1024)

	go decodeTextFromFinalCodeNumber(totalTextByteNumber, finalCodeNumber, characterMap, byteChannelToTextFile)

	fmt.Println("filePath", filePath)
	// 译码内容写入文件
	util.WriteByteToFile(filePath, byteChannelToTextFile)
}

func readCodeFrequencyAndGenerateCharacterMap(codeNumber uint8, byteChannelFromBinaryFile <-chan byte) (characterMap map[byte]*CharacterNode) {

	characterFrequencyMap := make(map[byte]uint32)
	var loopNumber uint16
	// 当为0时，说明有256个字符需要计入统计
	if codeNumber != 0 {
		loopNumber = uint16(codeNumber)
	} else {
		loopNumber = uint16(256)
	}

	for i := uint16(0); i < loopNumber; i++ {
		character := <-byteChannelFromBinaryFile
		var frequencyArray [4]byte
		for j := uint(0); j < 4; j++ {
			frequencyArray[j] = <-byteChannelFromBinaryFile
		}
		// fmt.Println(frequencyArray)
		frequencyInt := util.Couvert4ByteArrayToUint32(frequencyArray)
		characterFrequencyMap[character] = frequencyInt
	}

	var totalFrequency decimal.Decimal
	for _, v := range characterFrequencyMap {
		totalFrequency = totalFrequency.Add(decimal.NewFromInt32(int32(v)))
	}

	fmt.Println("totalFrequency", totalFrequency)

	characterMap = make(map[byte]*CharacterNode)

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

	return
}

func readTotalTextByteNumber(byteChannelFromBinaryFile <-chan byte) (totalTextByteNumber uint32) {

	var totalTextByteNumber4Array [4]byte
	for j := uint(0); j < 4; j++ {
		totalTextByteNumber4Array[j] = <-byteChannelFromBinaryFile
	}
	// fmt.Println(frequencyArray)
	totalTextByteNumber = util.Couvert4ByteArrayToUint32(totalTextByteNumber4Array)

	return
}

func readCode(byteChannelFromBinaryFile <-chan byte) (finalCodeByteSlice []byte) {

	finalCodeByteSlice = []byte{0x30, 0x2e}

	receivedByte, ok := <-byteChannelFromBinaryFile

	for ok {
		firstSemiByte := ((receivedByte & 0xf0) >> 4) + 0x30
		finalCodeByteSlice = append(finalCodeByteSlice, firstSemiByte)
		secondSemiByte := (receivedByte & 0x0f) + 0x30
		if secondSemiByte != 0x3a {
			finalCodeByteSlice = append(finalCodeByteSlice, secondSemiByte)
		}
		receivedByte, ok = <-byteChannelFromBinaryFile
	}

	return
}

func decodeTextFromFinalCodeNumber(totalTextByteNumber int, finalCodeNumber decimal.Decimal, characterMap map[byte]*CharacterNode, byteChannelToTextFile chan<- byte) {

	low := decimal.Zero
	high := decimal.NewFromInt(1)

	var characterNodeSlice []*CharacterNode

	for _, v := range characterMap {
		characterNodeSlice = append(characterNodeSlice, v)
	}

	sort.Sort(CharacterNodeSlice(characterNodeSlice))
	// for _, each := range characterNodeSlice {
	// 	fmt.Println(each)
	// }
	for i := 0; i < totalTextByteNumber; i++ {

		zeroDecimal := decimal.Zero
		oneDecimal := decimal.NewFromInt(1)
		// fmt.Println("low", low)
		// fmt.Println(finalCodeNumber)
		// fmt.Println("high", high)
		lowByteSlice, _ := low.MarshalText()
		highByteSlice, _ := high.MarshalText()
		// fmt.Println("lowByteSlice", lowByteSlice)
		// fmt.Println("highByteSlice", highByteSlice)

		lowByteSliceLength := len(lowByteSlice)
		highByteSliceLength := len(highByteSlice)

		truncatePosition := 0
		// FIXME: decimal变[]byte 可能不只是简单地直接变，还要再研究一下对应规则
		if lowByteSliceLength <= highByteSliceLength {
			for i := 2; i < lowByteSliceLength; i++ {
				if lowByteSlice[i] != highByteSlice[i] {
					truncatePosition = i - 2
					break
				}
			}
		} else {
			for i := 2; i < lowByteSliceLength; i++ {
				if highByteSlice[i] != lowByteSlice[i] {
					truncatePosition = i - 2
					break
				}
			}

		}
		// fmt.Println("truncatePosition", truncatePosition)
		if (!low.Equal(zeroDecimal)) || (!high.Equal(oneDecimal)) {
			lowShifted := low.Shift(int32(truncatePosition))
			// fmt.Println("lowShifted", lowShifted)
			low = lowShifted.Sub(lowShifted.Floor()).Truncate(20)
			// fmt.Println("low-update", low)
			highShifted := high.Shift(int32(truncatePosition))
			// fmt.Println("highShifted", highShifted)
			high = highShifted.Sub(highShifted.Floor()).Truncate(20)
			// fmt.Println("high-update", high)
			finalCodeNumberShifted := finalCodeNumber.Shift(int32(truncatePosition))
			finalCodeNumber = finalCodeNumberShifted.Sub(finalCodeNumberShifted.Floor()).Truncate(15)
		}

		// fmt.Println("low", low)
		// fmt.Println("high", high)
		for _, eachCharacterNode := range characterNodeSlice {
			tempLow := low.Add(eachCharacterNode.LeftBounded.Mul(high.Sub(low)))
			// fmt.Println("tempLow", tempLow)
			if finalCodeNumber.GreaterThanOrEqual(tempLow) {
				tempHigh := low.Add(eachCharacterNode.RightBounded.Mul(high.Sub(low)))
				// fmt.Println("tempHigh", tempHigh)
				if finalCodeNumber.LessThanOrEqual(tempHigh) {
					low = tempLow
					high = tempHigh
					// fmt.Printf("%c", eachCharacterNode.Character)
					byteChannelToTextFile <- eachCharacterNode.Character
					break
				}
			}
		}
	}
	close(byteChannelToTextFile)
}
