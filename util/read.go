package util

import (
    "fmt"
    "io/ioutil"
    "os"
)

// ReadText 从文件读取内容
func ReadText(filePath string) []byte {

    fileContent, err := ioutil.ReadFile(filePath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "ioutil.ReadFile() : %s\n", err)
    }

    // return string(fileContent)
    return fileContent
}
