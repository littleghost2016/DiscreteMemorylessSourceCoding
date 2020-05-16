package huffman

import (
	"DiscreteMemorylessSourceCoding/util"
	"fmt"
	"math"
	"os"
)

// EncodeHandler huffman编码处理
func EncodeHandler(filePath string, textByteSlice []byte) {

	if len(textByteSlice) == 0 {
		fmt.Println("There is no character in text!")
		os.Exit(1)
	}

	// 统计字符出现次数
	characterFrequencyMap := util.CountCharacterFromText(textByteSlice)

	// 生成霍夫曼树节点
	treeNodeMap := GenerateHuffmanTreeNode(characterFrequencyMap)

	// 生成霍夫曼树
	GenerateHuffmanTree(treeNodeMap)

	// 仅做测试用
	// PrintTreeMap(treeNodeMap)

	byteChannelToFile := make(chan byte, 64)

	// 准备二进制文件所需的数据
	go writeBinaryFile(treeNodeMap, characterFrequencyMap, textByteSlice, byteChannelToFile)

	// 构造输出文件名
	binaryFileName := fmt.Sprintf("%s.bin", filePath)
	fmt.Println("binaryFileName", binaryFileName)
	util.WriteByteToFile(binaryFileName, byteChannelToFile)
}

func encodeTextFromTreeNodeMap(text []byte, tnm map[byte]*TreeNode, bc chan<- bool) {

	for _, eachCharacher := range text {
		for _, eachCharacter := range tnm[eachCharacher].Code {
			// fmt.Printf("%T  %v\n", eachCharacter, eachCharacter)
			if eachCharacter == rune(48) {
				bc <- false
			} else if eachCharacter == rune(49) {
				bc <- true
			} else {
				fmt.Println("There is a wrong code that isn't 0 or isn't 1!")
				os.Exit(1)
			}
		}
	}
	close(bc)
}

func writeBinaryFile(tnm map[byte]*TreeNode, cfm map[byte]uint32, tbs []byte, byteChannelToFile chan<- byte) {

	// 0: huffman
	util.WriteFlag(0, byteChannelToFile)

	writeCodeNumber(tnm, byteChannelToFile)
	writeCodeMap(cfm, byteChannelToFile)
	paddingLength := calculatePaddingLength(tnm, cfm)
	util.WritePaddingLength(paddingLength, byteChannelToFile)
	writeCode(tbs, tnm, byteChannelToFile, paddingLength)
}

func writeCodeMap(cfm map[byte]uint32, byteChannelToFile chan<- byte) {

	for k, v := range cfm {
		byteChannelToFile <- k
		temp4byteArray := util.CounvertUint32ToByteSlice(uint32(v))
		for i := 0; i < 4; i++ {
			byteChannelToFile <- temp4byteArray[i]
		}
	}
}

func calculatePaddingLength(tnm map[byte]*TreeNode, cfm map[byte]uint32) (paddingLength uint8) {
	// paddingLength能用出现的频数和编码长度算出来
	var lastByteValidCodeLength uint8
	for k, v := range cfm {
		lastByteValidCodeLength += uint8(len(tnm[k].Code) * int(v))
	}
	// 注意此处两次模8，如果没有后面的一次，则可能算出来paddingLength为8，本来这时候不用再进行任何填充，但8表示又填了一个byte
	paddingLength = (8 - (lastByteValidCodeLength % 8)) % 8
	return
}

func writeCodeNumber(tnm map[byte]*TreeNode, byteChannelToFile chan<- byte) {
	codeNumber := uint8(len(tnm)) // 类型转换原因同上

	byteChannelToFile <- codeNumber
}

func writeCode(tbs []byte, tnm map[byte]*TreeNode, byteChannelToFile chan<- byte, calculatedPaddingLength uint8) {
	bitChannel := make(chan bool, int(math.Pow(2, 16)))
	go encodeTextFromTreeNodeMap(tbs, tnm, bitChannel)
	operatedPaddingLength := util.ConvertCodeBitToCodeByte(bitChannel, byteChannelToFile)
	if calculatedPaddingLength != operatedPaddingLength {
		fmt.Println(calculatedPaddingLength, operatedPaddingLength)
		fmt.Println("calculatePaddingLength != operatedPaddingLength")
	}
}
