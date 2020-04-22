package util

import (
    "fmt"
    "os"
    "testing"
)

func TestConvertCodeStringToCodeByte(t *testing.T) {
    bc := make(chan bool, 16)
    bsc := make(chan []byte, 8)
    a := []string{
        "010",
        "01100",
        "0101",
        "1",
        "10",
        "11111111",
        "1",
        "1",
    }
    go func() {
        for _, each := range a {
            for _, eachCharacter := range each {
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
    }()
    go ConvertCodeStringToCodeByte(bc, bsc)

    // each, ok := <-bsc
    // t.Log("1", each, ok)
    for {
        each, ok := <-bsc
        t.Log("2", each, ok)
        if ok != true {
            break
        }
    }
}
