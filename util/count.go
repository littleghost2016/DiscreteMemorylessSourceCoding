package util

// CountCharacter 统计自个字符的出现次数
func CountCharacter(text string) (characterFrequencyMap map[rune]int) {
	// func CountCharacter(text string) map[rune]int {
	characterFrequencyMap = make(map[rune]int)

	for _, eachCharacter := range text {
		if _, ok := characterFrequencyMap[eachCharacter]; ok {
			characterFrequencyMap[eachCharacter]++
		} else {
			characterFrequencyMap[eachCharacter] = 1
		}
	}

	return
}
