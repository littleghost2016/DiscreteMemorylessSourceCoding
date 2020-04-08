package huffman

import (
	"DiscreteMemorylessSourceCoding/util"
	"fmt"
	"os"
	"sort"
)

// EncodeHandler huffman编码处理
func EncodeHandler(filePath string) {

	// 读取文件
	text := util.ReadText(filePath)
	if len(text) == 0 {
		fmt.Println("There is no character in text!")
		os.Exit(1)
	}

	// 统计字符出现次数
	characterFrequencyMap := util.CountCharacter(text)

	// 生成霍夫曼树节点
	treeNodeMap := generateHuffmanTreeNode(characterFrequencyMap)

	// 生成霍夫曼树
	generateHuffmanTree(treeNodeMap)

	printEncodeResult(text, treeNodeMap)
}

// 生成节点
func generateHuffmanTreeNode(characterFrequencyMap map[rune]int) (treeNodeMap map[rune]*TreeNode) {
	treeNodeMap = make(map[rune]*TreeNode)
	for k, v := range characterFrequencyMap {
		if _, ok := treeNodeMap[k]; !ok {
			treeNodeMap[k] = &TreeNode{Character: k, Weight: v, FNode: nil, LNode: nil, RNode: nil, Code: ""}
		}
	}
	return
}

func generateHuffmanTree(treeNodeMap map[rune]*TreeNode) {
	tns := treeNodeMapToTreeNodeSlice(treeNodeMap)
	// 排序，当排序完成后，出现次数小的在前面
	sort.Sort(TreeNodeSlice(tns))

	if len(tns) == 1 {
		panic("There is none in TreeNodeSlice!")
	}

	for len(tns) != 1 {
		tempNode := TreeNode{Character: 0, Weight: tns[0].Weight + tns[1].Weight, LNode: tns[0], RNode: tns[1], FNode: nil, Code: ""}
		tns[0].FNode = &tempNode
		tns[1].FNode = tns[0].FNode
		tns = append(tns[2:], &tempNode)
		sort.Sort(TreeNodeSlice(tns))
	}

	rootNode := tns[0]
	distributeCode(rootNode)

	return
}

// 使用sort.Sort()需要实现以下三个函数: Len() Less() Swap
// 长度
func (tn TreeNodeSlice) Len() int {
	return len(tn)
}

// 默认从小到大排序，因此此处函数名为Less
func (tn TreeNodeSlice) Less(i, j int) bool {
	return tn[i].Weight < tn[j].Weight
}

// 交换函数
func (tn TreeNodeSlice) Swap(i, j int) {
	tn[i], tn[j] = tn[j], tn[i]
}

// map转slice
func treeNodeMapToTreeNodeSlice(tnm map[rune]*TreeNode) (trs []*TreeNode) {
	trs = make(TreeNodeSlice, 0)
	for _, v := range tnm {
		trs = append(trs, v)
	}
	return
}

func distributeCode(node *TreeNode) {
	if node.LNode != nil {
		node.LNode.Code = fmt.Sprintf("%s%d", node.Code, 0)
		distributeCode(node.LNode)
	}
	if node.RNode != nil {
		node.RNode.Code = fmt.Sprintf("%s%d", node.Code, 1)
		distributeCode(node.RNode)
	}
	return
}

func printEncodeResult(text string, tnm map[rune]*TreeNode) {
	for _, eachCharacher := range text {
		fmt.Print(tnm[eachCharacher].Code)
	}
}
