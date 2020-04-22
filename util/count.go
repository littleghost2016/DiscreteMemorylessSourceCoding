package util

// CountCharacterFromText 统计自个字符的出现次数
func CountCharacterFromText(text []byte) (characterFrequencyMap map[byte]int) {
	characterFrequencyMap = make(map[byte]int)

	for _, eachCharacter := range text {
		if _, ok := characterFrequencyMap[eachCharacter]; ok {
			characterFrequencyMap[eachCharacter]++
		} else {
			characterFrequencyMap[eachCharacter] = 1
		}
	}

	return
}
