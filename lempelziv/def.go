package lempelziv

type DecodingDirectoryNode struct {
	Type                uint8  `json:"type"`
	Character           []byte `json:"character"`
	SegmentNumber       int    `json:"segmentnumber"`
	LastCharacterNumber uint8  `json:"lastcharacternumber"`
	Code                string `json:"code"`
}

type SingleCharacterDirectoryNode struct {
	Type            uint8  `json:"type"`
	Character       byte   `json:"character"`
	CharacterNubmer uint8  `json:"characternubmer"`
	Code            string `json:"code"`
}
