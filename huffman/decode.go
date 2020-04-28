package huffman

import (
	"DiscreteMemorylessSourceCoding/util"
	"fmt"
	"math"
)

// DecodeHandler ...
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

	paddingLength := util.ReadPaddingLength(byteChannelFromBinaryFile)
	fmt.Println("paddingLength", paddingLength)

	bitChannel := make(chan bool, int(math.Pow(2, 16)))

	go util.ConvertCodeByteToCodeBit(paddingLength, byteChannelFromBinaryFile, bitChannel)

	byteChannelToTextFile := make(chan byte, 1024)

	go decodeTextFromTreeNodeMap(rootNode, bitChannel, byteChannelToTextFile)

	fmt.Println("filePath", filePath)
	util.WriteByteToFile(filePath, byteChannelToTextFile)

}

func readCodeFrequencyAndGenerateCharacterFrequencyMap(codeNumber uint8, byteChannel <-chan byte) (cfm map[byte]uint32) {

	cfm = make(map[byte]uint32)

	if codeNumber != 0 {
		for i := uint8(0); i < codeNumber; i++ {
			character := <-byteChannel
			var frequencyArray [4]byte
			for j := uint(0); j < 4; j++ {
				frequencyArray[j] = <-byteChannel
			}
			// fmt.Println(frequencyArray)
			frequencyInt := util.Couvert4ByteArrayToUint32(frequencyArray)
			cfm[character] = frequencyInt
		}
	} else {
		for i := uint16(0); i < 256; i++ {
			character := <-byteChannel
			var frequencyArray [4]byte
			for j := uint(0); j < 4; j++ {
				frequencyArray[j] = <-byteChannel
			}
			// fmt.Println(frequencyArray)
			frequencyInt := util.Couvert4ByteArrayToUint32(frequencyArray)
			cfm[character] = frequencyInt
		}
	}
	return
}

func decodeTextFromTreeNodeMap(rootNode *TreeNode, bitChannel <-chan bool, byteChannel chan<- byte) {
	currentNode := rootNode
	for each := range bitChannel {
		if each == false {
			currentNode = currentNode.LNode
		} else {
			currentNode = currentNode.RNode
		}
		if currentNode.IsLeafNode {
			byteChannel <- currentNode.Character
			// fmt.Println(currentNode.Character)
			currentNode = rootNode
		}
	}
	close(byteChannel)
}
