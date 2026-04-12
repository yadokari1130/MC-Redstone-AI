package model

// BlockData はMinecraftのブロックの座標、種類、および状態を表します。
type BlockData struct {
	X     int               `json:"x"`
	Y     int               `json:"y"`
	Z     int               `json:"z"`
	Block string            `json:"block"`
	State map[string]string `json:"state,omitempty"`
}

type AttachesData struct {
	ComponentX int    `json:"component_x"`
	ComponentY int    `json:"component_y"`
	ComponentZ int    `json:"component_z"`
	Component  string `json:"component"`
	BaseX      int    `json:"base_x"`
	BaseY      int    `json:"base_y"`
	BaseZ      int    `json:"base_z"`
}

type ConnectsData struct {
	FromX     int    `json:"from_x"`
	FromY     int    `json:"from_y"`
	FromZ     int    `json:"from_z"`
	ToX       int    `json:"to_x"`
	ToY       int    `json:"to_y"`
	ToZ       int    `json:"to_z"`
	Component string `json:"component"`
}

type PlaceRequest struct {
	Blocks   []BlockData    `json:"blocks"`
	Attaches []AttachesData `json:"attaches"`
	Connects []ConnectsData `json:"connects"`
}


// InteractionRequest はブロック操作のリクエストを表します。
type InteractionRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

// ItemInfo はドロップするアイテムの情報（IDと数量）を表します。
type ItemInfo struct {
	ID     string `json:"id"`
	Amount int    `json:"amount"`
}

// DropItemsRequest はアイテムドロップのリクエストを表します。
type DropItemsRequest struct {
	X     float64    `json:"x"`
	Y     float64    `json:"y"`
	Z     float64    `json:"z"`
	Items []ItemInfo `json:"items"`
}

// InventoryRequest はインベントリ設定のリクエストを表します。
type InventoryRequest struct {
	X     int        `json:"x"`
	Y     int        `json:"y"`
	Z     int        `json:"z"`
	Items []ItemInfo `json:"items"`
}

// CommandResult はコマンドの実行結果をAIがパースしやすい形式で表します。
type CommandResult struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}
