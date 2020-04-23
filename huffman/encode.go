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
	PrintTreeMap(treeNodeMap)

	// 此处用[]byte想为以后可能的编码非ASCII码做准备
	// 如只用做ASCII编码，[]byte可改为byte
	// 是否可以不改而直接适应非ASCII编码还未测试
	byteSliceChannel := make(chan []byte, 64)

	// 写二进制文件
	go writeBinaryFile(treeNodeMap, characterFrequencyMap, textByteSlice, byteSliceChannel)

	// binaryFileName := fmt.Sprintf("%s.bin", strings.Split(filePath, ".")[0])
	binaryFileName := fmt.Sprintf("%s.bin", filePath)
	fmt.Println(binaryFileName)
	util.WriteToFile(binaryFileName, byteSliceChannel)

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

func writeBinaryFile(tnm map[byte]*TreeNode, cfm map[byte]uint32, tbs []byte, bsc chan<- []byte) {

	// 0: huffman
	writeFlag(0, bsc)

	writeCodeNumber(tnm, bsc)
	writeCodeMap(cfm, bsc)
	paddingLength := calculatePaddingLength(tnm, cfm)
	writePaddingLength(paddingLength, bsc)
	writeCode(tbs, tnm, bsc, paddingLength)
	close(bsc)
}

// 写二进制文件标志
func writeFlag(encodeType uint8, bsc chan<- []byte) {
	bsc <- []byte{0x19, 0x15}
	// 上面这一句也可以写成如下，因为[]byte的初始化需要byte类型，而Go中uint8和byte是一样的。
	// 需要注意：Go中int8与byte不一样，因为int8的取值范围为-128~127
	// bsc <- []byte{uint8(0x19), uint8(0x15)}
	// 但不可以写成这样，会提示无法将int转化为byte
	// flag19 := 0x19
	// flag15 := 0x15
	// bsc <- []byte{flag19, flag15}

	switch encodeType {
	case 0:
		// huffman编码
		bsc <- []byte{0x00}
	case 1:
		// 算术编码
		bsc <- []byte{0x01}
	case 2:
		// LZ编码
		bsc <- []byte{0x02}
	}
}

func writeCodeNumber(tnm map[byte]*TreeNode, bsc chan<- []byte) {
	codeNumber := uint8(len(tnm)) // 类型转换原因同上

	bsc <- []byte{codeNumber}
}

func writeCodeMap(cfm map[byte]uint32, bsc chan<- []byte) {

	for k, v := range cfm {
		bsc <- []byte{k}
		bsc <- util.CounvertUint32ToByteSlice(uint32(v))
	}
}

func calculatePaddingLength(tnm map[byte]*TreeNode, cfm map[byte]uint32) (paddingLength uint8) {
	// paddingLength能用出现的频数和编码长度算出来
	var lastByteValidCodeLength uint8
	for k, v := range cfm {
		lastByteValidCodeLength += uint8(len(tnm[k].Code) * int(v))
	}
	paddingLength = 8 - (lastByteValidCodeLength % 8)
	return
}

func writePaddingLength(pl uint8, bsc chan<- []byte) {
	bsc <- []byte{pl}
}

func writeCode(tbs []byte, tnm map[byte]*TreeNode, bsc chan<- []byte, calculatedPaddingLength uint8) {
	bitChannel := make(chan bool, int(math.Pow(2, 32)))
	go encodeTextFromTreeNodeMap(tbs, tnm, bitChannel)
	operatedPaddingLength := util.ConvertCodeStringToCodeByte(bitChannel, bsc)
	if calculatedPaddingLength != operatedPaddingLength {
		fmt.Println(calculatedPaddingLength, operatedPaddingLength)
		fmt.Println("calculatePaddingLength != operatedPaddingLength")
	}
}
