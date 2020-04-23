package util

import (
	"fmt"
	"io/ioutil"
	"os"
)

func ReadFromFile(filePath string) (fileContent []byte) {

	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ioutil.ReadFile() : %s\n", err)
	}

	return
}

// ReadFlag 读取文件头标志
// 经过本脚本编码的文件，均以0x19 0x15开头
func ReadFlag(flag2Array []byte) (flagExist bool) {

	var byte19 uint8 = 0x19
	var byte15 uint8 = 0x15

	if flag2Array[0] != byte19 {
		fmt.Fprintf(os.Stderr, "The first byte is not 0x19. This file may not been encoded by my tools.")
		return false
	}
	if flag2Array[1] != byte15 {
		fmt.Fprintf(os.Stderr, "The second byte is not 0x15. This file may not been encoded by my tools.")
		return false
	}
	return true
}

// ReadCodeType 读取编码类型
func ReadCodeType(in byte) (codeType uint8) {
	codeType = in
	return
}

// ReadCodeNumber 读取被编码的字符个数
func ReadCodeNumber(byteChannel <-chan byte) (codeNumber uint8) {
	codeNumber = <-byteChannel
	fmt.Println("code number", codeNumber)
	return
}

func ReadPaddingLength(byteChannel <-chan byte) (paddingLength uint8) {
	paddingLength = <-byteChannel
	return
}
