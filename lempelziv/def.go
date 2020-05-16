package lempelziv

type DecodingDirectoryNode struct {
	Type                uint8  `json:"type"`
	Character           []byte `json:"character"`
	SelfSegmentNubmer   int    `json:"selfsegmentnumber"`
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

type TempWaitingToBeSendNode struct {
	TempSegmentNumber       uint32 `json:"tempsegmentnumber"`
	TempLastCharacterNumber uint8  `json:"templastcharacternumber"`
}
