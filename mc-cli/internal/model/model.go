package model

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

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
	Pos       []int             `json:"pos" yaml:"pos"`
	Component string            `json:"component" yaml:"component"`
	Base      []int             `json:"base" yaml:"base"`
	State     map[string]string `json:"state,omitempty" yaml:"state,omitempty"`
}

// ToCompact は [Component, Pos, Base] の形式に変換します。
func (a AttachesData) ToCompact() any {
	return []any{a.Component, a.Pos, a.Base}
}

type ConnectsData struct {
	From      []int             `json:"from" yaml:"from"`
	To        []int             `json:"to" yaml:"to"`
	Component string            `json:"component" yaml:"component"`
	State     map[string]string `json:"state,omitempty" yaml:"state,omitempty"`
}

// ToCompact は [Component, From, To] の形式に変換します。
func (c ConnectsData) ToCompact() any {
	return []any{c.Component, c.From, c.To}
}

type FillsData struct {
	From  []int             `json:"from" yaml:"from"`
	To    []int             `json:"to" yaml:"to"`
	Block string            `json:"block" yaml:"block"`
	State map[string]string `json:"state,omitempty" yaml:"state,omitempty"`
}

// ToCompact は [Block, From, To, State] の形式に変換します。
func (f FillsData) ToCompact() any {
	return []any{f.Block, f.From, f.To, f.State}
}

type PlaceRequest struct {
	Blocks   []BlockData    `json:"blocks" yaml:"blocks"`
	Attaches []AttachesData `json:"attaches" yaml:"attaches"`
	Connects []ConnectsData `json:"connects" yaml:"connects"`
	Fills    []FillsData    `json:"fills" yaml:"fills"`
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

// RawPlaceRequest はJSON/YAMLの生バイト列を保持するヘルパー型です。
// 単一オブジェクトと配列の両方を受け入れ、後で ParsePlaceRequestPhases で解析します。
type RawPlaceRequest struct {
	Data []byte
}

// UnmarshalYAML はyaml.v3ノードから生データを抽出します。
func (r *RawPlaceRequest) UnmarshalYAML(node *yaml.Node) error {
	var raw interface{}
	if err := node.Decode(&raw); err != nil {
		return err
	}
	data, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	r.Data = data
	return nil
}

// UnmarshalJSON はJSONバイト列をそのまま保持します。
func (r *RawPlaceRequest) UnmarshalJSON(data []byte) error {
	r.Data = data
	return nil
}

// ParsePlaceRequestPhases はJSONバイト列を解析し、PlaceRequestのフェーズ配列を返します。
// 単一オブジェクト { ... } の場合は1要素のスライスにラップします。
// 配列 [ { ... }, { ... } ] の場合はそのまま返します。
func ParsePlaceRequestPhases(data []byte) ([]PlaceRequest, error) {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return nil, fmt.Errorf("データが空です")
	}

	if data[0] == '[' {
		var reqs []PlaceRequest
		if err := json.Unmarshal(data, &reqs); err != nil {
			return nil, fmt.Errorf("PlaceRequest配列のパース失敗: %w", err)
		}
		return reqs, nil
	}

	var req PlaceRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("PlaceRequestのパース失敗: %w", err)
	}
	return []PlaceRequest{req}, nil
}
