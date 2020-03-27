package operatefile

import (
	"fmt"
	"io/ioutil"
	"os"
)

// ReadText 从文件读取内容
func ReadText(filePath string) []uint8 {
	inputFile, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "os.Open() failed.")
		return nil
	}

	defer inputFile.Close()

	text, err := ioutil.ReadAll(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ioutil.ReadAll() failed.")
		return nil
	}

	return text
}
