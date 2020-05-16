package lempelziv

import (
	"DiscreteMemorylessSourceCoding/util"
	"fmt"
	"math"
)

func DecodeHandler(filePath string, fileContent []byte) {

	byteChannelFromBinaryFile := make(chan byte, 1024)

	// 读取文件
	go func(byteChannelFromBinaryFile chan<- byte) {
		for _, eachByte := range fileContent {
			byteChannelFromBinaryFile <- eachByte
		}
		close(byteChannelFromBinaryFile)
	}(byteChannelFromBinaryFile)

	codeNumber := util.ReadCodeNumber(byteChannelFromBinaryFile)
	// singleCharacterDirectorySlice, singleCharacterDirectory := readSingleCharacterDirectory(codeNumber, byteChannelFromBinaryFile)
	singleCharacterDirectorySlice, _ := readSingleCharacterDirectory(codeNumber, byteChannelFromBinaryFile)
	// _, _ = readSingleCharacterDirectory(codeNumber, byteChannelFromBinaryFile)

	// decodingDirectoryLength = readDecodingDirectoryLength(byteChannelFromBinaryFile)
	segmentLength := readSegmentLength(byteChannelFromBinaryFile)
	lastCharacterLength := readLastCharacterLength(byteChannelFromBinaryFile)
	lastIsSpecialFlag := readLastIsSpecialFlag(byteChannelFromBinaryFile)
	paddingLength := util.ReadPaddingLength(byteChannelFromBinaryFile)
	// fmt.Println(segmentLength, lastCharacterLength, lastIsSpecialFlag, paddingLength)

	bitChannel := make(chan bool, int(math.Pow(2, 16)))

	// 读取编码文件的编码部分
	go util.ConvertCodeByteToCodeBit(paddingLength, byteChannelFromBinaryFile, bitChannel)

	byteChannelToTextFile := make(chan byte, 1024)

	// 译码
	go decodeTextFromSingleCharacterDirectorySlice(singleCharacterDirectorySlice, segmentLength, lastCharacterLength, lastIsSpecialFlag, bitChannel, byteChannelToTextFile)

	fmt.Println("filePath", filePath)
	// 译码内容写入文件
	util.WriteByteToFile(filePath, byteChannelToTextFile)
}

func readSingleCharacterDirectory(codeNumber uint8, byteChannelFromBinaryFile <-chan byte) (singleCharacterDirectorySlice []*SingleCharacterDirectoryNode, singleCharacterDirectory map[byte]*SingleCharacterDirectoryNode) {

	singleCharacterDirectory = make(map[byte]*SingleCharacterDirectoryNode)

	var singleCharacterDirectorySliceIndex uint8 = 0
	var loopNumber int
	if codeNumber == 0 {
		loopNumber = 256
	} else {
		loopNumber = int(codeNumber)
	}
	for i := 0; i < loopNumber; i++ {
		eachTextByte := <-byteChannelFromBinaryFile
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
	}

	return
}

func readSegmentLength(byteChannelFromBinaryFile <-chan byte) (segmentLength uint32) {
	var segmentLengthArray [4]byte
	for i := uint8(0); i < 4; i++ {
		segmentLengthArray[i] = <-byteChannelFromBinaryFile
	}
	segmentLength = util.Couvert4ByteArrayToUint32(segmentLengthArray)
	return
}

func readLastCharacterLength(byteChannelFromBinaryFile <-chan byte) (lastCharacterLength uint8) {
	lastCharacterLength = <-byteChannelFromBinaryFile
	return
}

func readLastIsSpecialFlag(byteChannelFromBinaryFile <-chan byte) (lastIsSpecialFlag bool) {

	// TODO: 判断接收到是0还是1
	receivedBit := <-byteChannelFromBinaryFile
	if receivedBit == 0x00 {
		lastIsSpecialFlag = false
	} else {
		lastIsSpecialFlag = true
	}
	return

}

func decodeTextFromSingleCharacterDirectorySlice(singleCharacterDirectorySlice []*SingleCharacterDirectoryNode, segmentLength uint32, lastCharacterLength uint8, lastIsSpecialFlag bool, bitChannel <-chan bool, byteChannelToTextFile chan<- byte) {

	var decodingDirectorySlice []*DecodingDirectoryNode

	selfSegmentNumberIndex := 0

	if lastIsSpecialFlag {

		var tempSegmentNumber uint32
		var tempLastCharacterNumber uint8

		var receivedBit bool
		// 计算tempSegmentNumber
		for i := segmentLength; i > 0; i-- {
			receivedBit := <-bitChannel
			if receivedBit {
				tempSegmentNumber += (0x1 << (i - 1))
			}
		}
		// 计算tempLastCharacterNumber
		for i := lastCharacterLength; i > 0; i-- {
			receivedBit = <-bitChannel
			if receivedBit {
				tempLastCharacterNumber += (0x1 << (i - 1))
			}
		}
		// fmt.Println("段号", tempSegmentNumber, "字符号", tempLastCharacterNumber)
		var tempWaitingToBeSendNode1 TempWaitingToBeSendNode
		// var tempWaitingToBeSendNode2 tempWaitingToBeSendNode

		tempWaitingToBeSendNode1 = TempWaitingToBeSendNode{
			TempSegmentNumber:       tempSegmentNumber,
			TempLastCharacterNumber: tempLastCharacterNumber,
		}

		for {
			// fmt.Println(tempWaitingToBeSendNode1.TempSegmentNumber, tempWaitingToBeSendNode1.TempLastCharacterNumber)
			if receivedBit, ok := <-bitChannel; ok {

				// 发送tempWaitingToBeSendNode1
				convertDecimalToCharacterAndSendToChannel(tempWaitingToBeSendNode1.TempSegmentNumber, tempWaitingToBeSendNode1.TempLastCharacterNumber, &decodingDirectorySlice, singleCharacterDirectorySlice, selfSegmentNumberIndex, byteChannelToTextFile)

				// 段号自增
				selfSegmentNumberIndex++

				// 重置临时变量
				tempSegmentNumber = 0
				tempLastCharacterNumber = 0

				// 计算tempSegmentNumber
				if receivedBit {
					tempSegmentNumber += (0x1 << (segmentLength - 1))
				}
				for i := segmentLength - 1; i > 0; i-- {
					receivedBit = <-bitChannel
					if receivedBit {
						tempSegmentNumber += (0x1 << (i - 1))
					}
				}

				// 计算tempLastCharacterNumber
				for i := lastCharacterLength; i > 0; i-- {
					receivedBit = <-bitChannel
					if receivedBit {
						tempLastCharacterNumber += (0x1 << (i - 1))
					}
				}
				// fmt.Println("段号", tempSegmentNumber, "字符号", tempLastCharacterNumber)

				tempWaitingToBeSendNode1 = TempWaitingToBeSendNode{
					TempSegmentNumber:       tempSegmentNumber,
					TempLastCharacterNumber: tempLastCharacterNumber,
				}

			} else {
				// 只发送前面的段发送至channel
				// (*decodingDirectorySlice)要加括号，表示优先运算*
				tempSegmentCharacter := decodingDirectorySlice[tempSegmentNumber-1].Character
				for _, each := range tempSegmentCharacter {
					byteChannelToTextFile <- each
				}
				break
			}
		}
	} else {

		// for _, each := range decodingDirectorySlice {
		// 	fmt.Println("111", each)
		// }

		var tempSegmentNumber uint32
		var tempLastCharacterNumber uint8
		for {
			if receivedBit, ok := <-bitChannel; ok {
				// 计算tempSegmentNumber
				if receivedBit {
					tempSegmentNumber += (0x1 << (segmentLength - 1))
				}
				for i := segmentLength - 1; i > 0; i-- {
					receivedBit = <-bitChannel
					if receivedBit {
						tempSegmentNumber += (0x1 << (i - 1))
					}
				}
				// 计算tempLastCharacterNumber
				for i := lastCharacterLength; i > 0; i-- {
					receivedBit = <-bitChannel
					if receivedBit {
						tempLastCharacterNumber += (0x1 << (i - 1))
					}
				}
			} else {
				break
			}

			// 注意：此处因为要对decodingDirectorySlice进行修改，所以要传入地址以使修改成功应用到原处
			convertDecimalToCharacterAndSendToChannel(tempSegmentNumber, tempLastCharacterNumber, &decodingDirectorySlice, singleCharacterDirectorySlice, selfSegmentNumberIndex, byteChannelToTextFile)

			// 段号自增
			selfSegmentNumberIndex++

			// 重置临时变量，忘记好几次了
			tempSegmentNumber = 0
			tempLastCharacterNumber = 0
		}
	}

	// for _, each := range decodingDirectorySlice {
	// 	fmt.Println(each)
	// }
	// 关闭channel
	close(byteChannelToTextFile)
}

// 此处decodingDirectorySlice类型为[]*DecodingDirectoryNode的指针，即*[]*DecodingDirectoryNode
func convertDecimalToCharacterAndSendToChannel(tempSegmentNumber uint32, tempLastCharacterNumber uint8, decodingDirectorySlice *[]*DecodingDirectoryNode, singleCharacterDirectorySlice []*SingleCharacterDirectoryNode, selfSegmentNumberIndex int, byteChannelToTextFile chan<- byte) {
	var tempSegmentCharacter []byte
	var tempLastCharacter byte

	if tempSegmentNumber == 0 {
		// 发送至channel
		tempLastCharacter = singleCharacterDirectorySlice[tempLastCharacterNumber].Character
		byteChannelToTextFile <- tempLastCharacter
	} else {
		// 注意此处应重新申请地址，为下面的copy做准备
		// 若不重新申请，则tempSegmentCharacter的地址均为(*decodingDirectorySlice)[tempSegmentNumber-1].Character的地址
		// 即只想同意切片，后续的修改会更改前面的，这个BUG排查了我一个晚上+第二天的一个上午....
		// (*decodingDirectorySlice)要加括号，表示优先运算*
		tempSegmentCharacter = make([]byte, len((*decodingDirectorySlice)[tempSegmentNumber-1].Character))

		// 发送至channel
		copy(tempSegmentCharacter, ((*decodingDirectorySlice)[tempSegmentNumber-1].Character))
		for _, each := range tempSegmentCharacter {
			byteChannelToTextFile <- each
		}
		tempLastCharacter = singleCharacterDirectorySlice[tempLastCharacterNumber].Character
		byteChannelToTextFile <- tempLastCharacter

	}

	newTempCharacter := tempSegmentCharacter
	newTempCharacter = append(newTempCharacter, tempLastCharacter)
	tempDecodingDirectoryNode := DecodingDirectoryNode{
		Type:                0,
		Character:           newTempCharacter,
		SelfSegmentNubmer:   selfSegmentNumberIndex + 1,
		SegmentNumber:       int(tempSegmentNumber),
		LastCharacterNumber: tempLastCharacterNumber,
		Code:                "",
	}
	// 排查上面说的bug用的
	// fmt.Printf("tempDecodingDirectoryNode: %p\n", tempDecodingDirectoryNode.Character)
	// fmt.Printf("newTempCharacter: %p\n", newTempCharacter)
	// fmt.Printf("tempSegmentCharacter: %p\n", tempSegmentCharacter)
	// fmt.Printf("tempLastCharacter: %p\n", tempLastCharacter)

	*decodingDirectorySlice = append(*decodingDirectorySlice, &tempDecodingDirectoryNode)
	// 排查上面说的bug用的
	// for _, each := range *decodingDirectorySlice {
	// 	fmt.Println(each)
	// }
	// fmt.Println()
}
