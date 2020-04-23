package util

// ConvertCodeStringToCodeByte 将string类型的编码转变成可写入文件的[]byte
func ConvertCodeStringToCodeByte(bc <-chan bool, bsc chan<- []byte) (paddingLength uint8) {
	byteLength := uint8(8)
	tempByte := uint8(0) // 结合tempBytePointer

	paddingLength = uint8(0)

	var tempBytePointer *uint8 = &tempByte // 重复使用
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
			if byteLength != 0 {
				paddingLength = byteLength
				for i := byteLength; i != 0; i-- {
					*tempBytePointer <<= 1
				}
				bsc <- []byte{*tempBytePointer}
			}
			break
		}
	}
	return
}

// TODO
func ConvertCodeByteToCodeString() {
}

// CounvertUint32ToByteSlice ...
func CounvertUint32ToByteSlice(in uint32) (out []byte) {
	out = []byte{}
	for i := 3; i >= 0; i-- {
		eachByte := uint8(((uint32(0x000f) << (i * 8)) & in) >> (i * 8))
		out = append(out, eachByte)
	}
	return
}

// Couvert4ByteArrayToUint32 ...
func Couvert4ByteArrayToUint32(in [4]byte) (out uint32) {
	for i := uint8(0); i < 4; i++ {
		out = uint32(in[i]) << ((3 - i) * 8)
	}
	return out
}
