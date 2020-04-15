package util

// ConvertCodeStringToCodeByte 将string类型的编码转变成可写入文件的[]byte
func ConvertCodeStringToCodeByte(bc <-chan bool, bsc chan<- []byte) {
    byteLength := 8
    tempByte := uint8(0)
    var tempBytePointer *uint8 = &tempByte
    for {
        receivedBit, ok := <-bc
        if ok {
            byteLength--
            switch receivedBit {
            case false:
                *tempBytePointer <<= 1
            case true:
                *tempBytePointer = (*tempBytePointer << 1) + 1
            }
            if byteLength == 0 {
                bsc <- []byte{*tempBytePointer}
                *tempBytePointer = 0
                byteLength = 8
            }
        } else {
            for i := byteLength; i != 0; i-- {
                *tempBytePointer <<= 1
            }
            bsc <- []byte{*tempBytePointer}
            close(bsc)
            break
        }
    }
}
