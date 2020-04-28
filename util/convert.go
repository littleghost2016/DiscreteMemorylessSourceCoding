package util

import (
	"fmt"
)

// ConvertCodeStringToCodeByte 将string类型的编码转变成可写入文件的[]byte
func ConvertCodeStringToCodeByte(bc <-chan bool, byteChannelToFile chan<- byte) (paddingLength uint8) {
	byteLength := uint8(0)
	tempByte := uint8(0) // 结合tempBytePointer

	var tempBytePointer *uint8 = &tempByte // 重复使用
	for {
		receivedBit, ok := <-bc
		if ok {
			// fmt.Println("1", byteLength)
			switch receivedBit {
			case false:
				*tempBytePointer <<= 1
			case true:
				*tempBytePointer = (*tempBytePointer << 1) + 1
			}
			byteLength++
			if byteLength == 8 {
				byteChannelToFile <- *tempBytePointer
				*tempBytePointer = 0
				byteLength = 0
			}
		} else {
			if byteLength != 0 {
				paddingLength = 8 - byteLength
				for i := byteLength; i != 8; i++ {
					*tempBytePointer <<= 1
				}
				byteChannelToFile <- *tempBytePointer
			}
			close(byteChannelToFile)
			break
		}
	}
	return
}

// ConvertCodeByteToCodeBit ...
func ConvertCodeByteToCodeBit(paddingLength uint8, inByteChannel <-chan byte, outBitChannel chan<- bool) {
	var tempByte byte
	receivedByte1, ok1 := <-inByteChannel
	receivedByte2, ok2 := <-inByteChannel
	for {
		if ok1 && ok2 {
			for i := uint8(8); i > uint8(0); i-- {
				tempByte = receivedByte1 & (uint8(1) << (i - 1))
				if tempByte != uint8(0) {
					outBitChannel <- true
				} else {
					outBitChannel <- false
				}
			}
			receivedByte1, ok1 = receivedByte2, true
			receivedByte2, ok2 = <-inByteChannel
		} else {
			for i := uint8(8); i > paddingLength; i-- {
				tempByte = receivedByte1 & (uint8(1) << (i - 1))
				if tempByte != uint8(0) {
					outBitChannel <- true
				} else {
					outBitChannel <- false
				}
			}
			close(outBitChannel)
			break
		}

	}
	fmt.Println("3")

}

// CounvertUint32ToByteSlice ...
func CounvertUint32ToByteSlice(in uint32) (out [4]byte) {
	out = [4]byte{}
	// out在内存中的存放顺序是out[0] out[1] out[2] out[3]
	for i := 3; i >= 0; i-- {
		eachByte := uint8(((uint32(0x000000ff) << (i * 8)) & in) >> (i * 8))

		// 此处应该按照顺序放置，对应如下
		// 00000000 00000000 00000000 00000000
		// out[0]   out[1]   out[2]   out[3]
		out[3-i] = eachByte
	}
	return
}

// Couvert4ByteArrayToUint32 ...
func Couvert4ByteArrayToUint32(in [4]byte) (out uint32) {
	// fmt.Println(in)
	for i := uint8(0); i < 4; i++ {
		out = (out << 8) | uint32(in[i])
	}
	return out
}
