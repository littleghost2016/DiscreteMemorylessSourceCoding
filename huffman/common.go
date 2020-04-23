package huffman

import (
	"fmt"
	"sort"
)

// GenerateHuffmanTreeNode 生成节点
func GenerateHuffmanTreeNode(characterFrequencyMap map[byte]uint32) (treeNodeMap map[byte]*TreeNode) {
	treeNodeMap = make(map[byte]*TreeNode)
	for k, v := range characterFrequencyMap {
		if _, ok := treeNodeMap[k]; !ok {
			treeNodeMap[k] = &TreeNode{
				Character: k,
				Weight:    v,
				FNode:     nil,
				LNode:     nil,
				RNode:     nil,
				Code:      "",
				LeafNode:  true,
			}
		}
	}
	return
}

// GenerateHuffmanTree 生成（动词）树
func GenerateHuffmanTree(treeNodeMap map[byte]*TreeNode) {

	// 验证：map是无序的
	// for k, v := range treeNodeMap {
	//  fmt.Println(k, v)
	// }

	tns := treeNodeMapToTreeNodeSlice(treeNodeMap)
	// 排序，当排序完成后，出现次数小的在前面
	// 出现次数相同时，ASCII码小的排在前面
	sort.Sort(TreeNodeSlice(tns))

	for _, each := range tns {
		fmt.Println(*each)
	}

	if len(tns) == 1 {
		panic("There is none in TreeNodeSlice!")
	}

	for len(tns) != 1 {
		tempNode := TreeNode{
			Character: 0,
			Weight:    tns[0].Weight + tns[1].Weight,
			LNode:     tns[0],
			RNode:     tns[1],
			FNode:     nil,
			Code:      "",
			LeafNode:  false,
		}
		tns[0].FNode = &tempNode
		tns[1].FNode = tns[0].FNode
		tns = append(tns[2:], &tempNode)
		sort.Sort(TreeNodeSlice(tns))
	}

	// 测试：输出tns的每一个
	// for _, each := range tns {
	//  fmt.Println(*each)
	// }

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
	// 在tn里面，相同出现次数的不同字符 的 排序后所在位置会影响编码的分配。
	// 此处增加后面的条件是保证第一次排序时，相同出现次数的字符，ASCII码小的在前面
	// 可解决互为兄弟节点的节点无法每次都绝对确定自己是0还是1的问题
	// 这个bug排查了我半个小时
	return (tn[i].Weight < tn[j].Weight) || (tn[i].Weight == tn[j].Weight && tn[i].Character < tn[j].Character)
}

// 交换函数
func (tn TreeNodeSlice) Swap(i, j int) {
	tn[i], tn[j] = tn[j], tn[i]
}

// map转slice
func treeNodeMapToTreeNodeSlice(tnm map[byte]*TreeNode) (trs []*TreeNode) {
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

// PrintTreeMap 输出树map
func PrintTreeMap(tnm map[byte]*TreeNode) {
	for k, v := range tnm {
		fmt.Println(k, v)
	}
}
