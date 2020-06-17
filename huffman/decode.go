package huffman

import (
	"DiscreteMemorylessSourceCoding/util"
	"fmt"
	"math"
)

// DecodeHandler 译码过程
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
	characterFrequencyMap := readCodeFrequencyAndGenerateCharacterFrequencyMap(codeNumber, byteChannelFromBinaryFile)

	// 测试characterFrequencyMap
	// for k, v := range characterFrequencyMap {
	// 	fmt.Println(k, v)
	// }

	// 生成霍夫曼树节点
	treeNodeMap := GenerateHuffmanTreeNode(characterFrequencyMap)
	// 生成霍夫曼树
	rootNode := GenerateHuffmanTree(treeNodeMap)

	// 仅做测试用
	// PrintTreeMap(treeNodeMap)

	// 获得填充长度
	paddingLength := util.ReadPaddingLength(byteChannelFromBinaryFile)
	fmt.Println("paddingLength", paddingLength)

	bitChannel := make(chan bool, int(math.Pow(2, 16)))

	// 读取编码文件的编码部分
	go util.ConvertCodeByteToCodeBit(paddingLength, byteChannelFromBinaryFile, bitChannel)

	byteChannelToTextFile := make(chan byte, 1024)

	// 译码
	go decodeTextFromTreeNodeMap(rootNode, bitChannel, byteChannelToTextFile)

	fmt.Println("filePath", filePath)
	// 译码内容写入文件
	util.WriteByteToFile(filePath, byteChannelToTextFile)
}

func readCodeFrequencyAndGenerateCharacterFrequencyMap(codeNumber uint8, byteChannel <-chan byte) (cfm map[byte]uint32) {

	cfm = make(map[byte]uint32)

	var loopNumber uint16
	// 当为0时，说明有256个字符需要计入统计
	if codeNumber != 0 {
		loopNumber = uint16(codeNumber)
	} else {
		loopNumber = uint16(256)
	}

	for i := uint16(0); i < loopNumber; i++ {
		character := <-byteChannel
		var frequencyArray [4]byte
		for j := uint(0); j < 4; j++ {
			frequencyArray[j] = <-byteChannel
		}
		// fmt.Println(frequencyArray)
		frequencyInt := util.Couvert4ByteArrayToUint32(frequencyArray)
		cfm[character] = frequencyInt
	}

	return
}

// 根据霍夫曼树将比特流译码成字符
func decodeTextFromTreeNodeMap(rootNode *TreeNode, bitChannel <-chan bool, byteChannel chan<- byte) {
	currentNode := rootNode
	for each := range bitChannel {
		if each == false {
			// 0则进入左子树
			currentNode = currentNode.LNode
		} else {
			// 1则进入右子树
			currentNode = currentNode.RNode
		}
		// 若是叶节点则输出字符
		if currentNode.IsLeafNode {
			byteChannel <- currentNode.Character
			// fmt.Println(currentNode.Character)
			currentNode = rootNode
		}
	}
	close(byteChannel)
}
