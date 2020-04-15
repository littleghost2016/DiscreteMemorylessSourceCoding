package util

import (
	"fmt"
	"os"
)

// WriteCodeToBinaryFile ...
func WriteCodeToBinaryFile(fileName string, bsc <-chan []byte) {

	fileObject, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("WriteEncodedBinaryFile function failed.")
	}

	defer fileObject.Close()

	// 使用buffer时
	// binBuffer := new(bytes.Buffer)

	for eachByteSlice := range bsc {
		fileObject.Write(eachByteSlice)
	}

	// 使用buffer时
	// //使用Write方法,需要使用Writer对象的Flush方法将buffer中的数据刷到磁盘
	// for eachByteSlice := range bsc {
	// 	fmt.Println(eachByteSlice)
	// 	binary.Write(binBuffer, binary.BigEndian, eachByteSlice)
	// 	fileObject.Write(binBuffer.Bytes())
	// 	binBuffer.Reset()
	// }
}
