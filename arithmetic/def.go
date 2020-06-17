package arithmetic

import "github.com/shopspring/decimal"

// CharacterNode 字符节点
type CharacterNode struct {
	Character    uint8           `json:"character"`
	Frequency    uint32          `json:"frequency"`
	Weight       decimal.Decimal `json:"weight"`
	LeftBounded  decimal.Decimal `json:"leftbounded"`
	RightBounded decimal.Decimal `json:"rightbounded"`
}

// CharacterNodeSlice 字符节点切片
type CharacterNodeSlice []*CharacterNode
