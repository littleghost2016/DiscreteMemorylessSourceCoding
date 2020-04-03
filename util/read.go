package operatefile

import (
	"fmt"
	"io/ioutil"
	"os"
)

// ReadText 从文件读取内容
func ReadText(filePath string) string {
	inputFile, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "os.Open() failed.")
		return ""
	}

	defer inputFile.Close()

	text, err := ioutil.ReadAll(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ioutil.ReadAll() failed.")
		return ""
	}

	return string(text)
}
