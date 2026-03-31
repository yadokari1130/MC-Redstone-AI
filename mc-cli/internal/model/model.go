package model

// BlockData はMinecraftのブロックの座標、種類、および状態を表します。
type BlockData struct {
	X     int               `json:"x"`
	Y     int               `json:"y"`
	Z     int               `json:"z"`
	Block string            `json:"block"`
	State map[string]string `json:"state,omitempty"`
}

// InteractionRequest はブロック操作のリクエストを表します。
type InteractionRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

// CommandResult はコマンドの実行結果をAIがパースしやすい形式で表します。
type CommandResult struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}
