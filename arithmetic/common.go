package arithmetic

// 使用sort.Sort()需要实现以下三个函数: Len() Less() Swap
// 长度
func (cns CharacterNodeSlice) Len() int {
	return len(cns)
}

// 默认从小到大排序，因此此处函数名为Less
func (cns CharacterNodeSlice) Less(i, j int) bool {
	return cns[i].Character < cns[j].Character
}

// 交换函数
func (cns CharacterNodeSlice) Swap(i, j int) {
	cns[i], cns[j] = cns[j], cns[i]
}
