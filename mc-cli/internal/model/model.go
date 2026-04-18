package model

// BlockData はMinecraftのブロックの座標、種類、および状態を表します。
type BlockData struct {
	X     int               `json:"x" yaml:"x"`
	Y     int               `json:"y" yaml:"y"`
	Z     int               `json:"z" yaml:"z"`
	Block string            `json:"block" yaml:"block"`
	State map[string]string `json:"state,omitempty" yaml:"state,omitempty"`
}

// ToCompact は [Block, [X, Y, Z], State] の形式に変換します。
func (b BlockData) ToCompact() any {
	return []any{b.Block, []int{b.X, b.Y, b.Z}, b.State}
}

type AttachesData struct {
	Pos       []int  `json:"pos" yaml:"pos"`
	Component string `json:"component" yaml:"component"`
	Base      []int  `json:"base" yaml:"base"`
}

// ToCompact は [Component, Pos, Base] の形式に変換します。
func (a AttachesData) ToCompact() any {
	return []any{a.Component, a.Pos, a.Base}
}

type ConnectsData struct {
	From      []int  `json:"from" yaml:"from"`
	To        []int  `json:"to" yaml:"to"`
	Component string `json:"component" yaml:"component"`
}

// ToCompact は [Component, From, To] の形式に変換します。
func (c ConnectsData) ToCompact() any {
	return []any{c.Component, c.From, c.To}
}

type PlaceRequest struct {
	Blocks   []BlockData    `json:"blocks" yaml:"blocks"`
	Attaches []AttachesData `json:"attaches" yaml:"attaches"`
	Connects []ConnectsData `json:"connects" yaml:"connects"`
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
