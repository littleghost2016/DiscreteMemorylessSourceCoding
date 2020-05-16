package util

import (
	"bytes"
	"fmt"
	"os"
)

// WriteByteToFile ...
func WriteByteToFile(filePath string, byteChannel <-chan byte) {

	fileObject, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err, "WriteEncodedBinaryFile function failed.")
	}

	defer fileObject.Close()

	// 使用buffer时
	binBuffer := new(bytes.Buffer)

	for eachByte := range byteChannel {
		// fmt.Println(eachByte)
		binBuffer.WriteByte(eachByte)
	}
	fileObject.Write(binBuffer.Bytes())
}

// WriteFlag 写二进制文件标志
func WriteFlag(encodeType uint8, byteChannelToFile chan<- byte) {

	byteChannelToFile <- 0x19
	byteChannelToFile <- 0x15

	switch encodeType {
	case 0:
		// huffman编码
		byteChannelToFile <- 0x00
	case 1:
		// 算术编码
		byteChannelToFile <- 0x01
	case 2:
		// LZ编码
		byteChannelToFile <- 0x02
	}
}

// WritePaddingLength 写填充长度
func WritePaddingLength(pl uint8, byteChannelToFile chan<- byte) {
	byteChannelToFile <- pl
}

func WriteByteSliceToFile(filePath string, in []byte) {

	fileObject, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err, "WriteEncodedBinaryFile function failed.")
	}

	defer fileObject.Close()

	fileObject.Write(in)

}
