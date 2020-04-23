package huffman

import (
	"DiscreteMemorylessSourceCoding/util"
	"fmt"
)

// DecodeHandler ...
func DecodeHandler(fileContent []byte) {

	byteChannel := make(chan byte, 1024)

	// 读取文件
	go func(byteChannel chan<- byte) {
		for _, eachByte := range fileContent {
			byteChannel <- eachByte
		}
	}(byteChannel)

	codeNumber := util.ReadCodeNumber(byteChannel)
	characterFrequencyMap := readCodeFrequencyAndGenerateCharacterFrequencyMap(codeNumber, byteChannel)

	// 测试characterFrequencyMap
	// for k, v := range characterFrequencyMap {
	// 	fmt.Println(k, v)
	// }

	// 生成霍夫曼树节点
	treeNodeMap := GenerateHuffmanTreeNode(characterFrequencyMap)
	// 生成霍夫曼树
	GenerateHuffmanTree(treeNodeMap)

	// 仅做测试用
	PrintTreeMap(treeNodeMap)

	paddingLength := util.ReadPaddingLength(byteChannel)
	fmt.Println("paddingLength", paddingLength)

	writeByteChannel := make(chan byte, 1024)

	// util.ConvertCodeByteToCodeString(paddingLength, byteChannel, writeByteChannel)
	util.ConvertCodeByteToCodeString()

}

func readCodeFrequencyAndGenerateCharacterFrequencyMap(codeNumber uint8, byteChannel <-chan byte) (cfm map[byte]uint32) {

	cfm = make(map[byte]uint32)

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
	return
}
