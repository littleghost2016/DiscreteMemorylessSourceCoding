package util

// CountCharacterFromText 统计自个字符的出现次数
func CountCharacterFromText(text []byte) (characterFrequencyMap map[byte]uint32) {
	characterFrequencyMap = make(map[byte]uint32)

	for _, eachCharacter := range text {
		if _, ok := characterFrequencyMap[eachCharacter]; ok {
			characterFrequencyMap[eachCharacter]++
		} else {
			characterFrequencyMap[eachCharacter] = 1
		}
	}

	return
}
