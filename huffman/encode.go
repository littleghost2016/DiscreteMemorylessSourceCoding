package huffman

import (
	"DiscreteMemorylessSourceCoding/util"
	"fmt"
	"os"
)

// EncodeHandler huffman编码处理
func EncodeHandler(filePath string) {

	// 读取文件
	textByteSlice := util.ReadText(filePath)
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

	bitChannel := make(chan bool, 32)

	// 此处用[]byte想为以后可能的编码非ASCII码做准备
	// 如只用做ASCII编码，[]byte可改为byte
	byteSliceChannel := make(chan []byte, 8)

	go encodeTextFromTreeNodeMap(textByteSlice, treeNodeMap, bitChannel)

	// go convertCodeStringToCodeByte(codeStringChannel, byteSliceChannel)
	go util.ConvertCodeStringToCodeByte(bitChannel, byteSliceChannel)

	util.WriteCodeToBinaryFile("testOutput.bin", byteSliceChannel)

	// printEncodeResult(text, treeNodeMap)
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
