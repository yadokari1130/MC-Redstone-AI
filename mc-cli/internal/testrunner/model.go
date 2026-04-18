package testrunner

import "mc-cli/internal/model"

// TestFile はテストYAMLファイルのトップレベル構造を表します。
type TestFile struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Tests       []TestCase `yaml:"tests"`
}

// TestCase は1つのテストケースを表します。
type TestCase struct {
	Name   string     `yaml:"name"`
	Setup  Setup      `yaml:"setup"`
	Steps  []Step     `yaml:"steps"`
	Assert []Assertion `yaml:"assert"`
}

// Setup はテスト開始前に実行する回路配置設定を表します。
// blocks_file（ファイルパス）またはblocks（JSONインライン）のどちらか一方を指定します。
type Setup struct {
	BlocksFile string           `yaml:"blocks_file"` // PlaceRequest形式のJSONファイルパス
	Blocks     *model.PlaceRequest `yaml:"blocks"`      // インラインで指定するPlaceRequest
}

// Step はテストの1つの操作ステップを表します。
type Step struct {
	// action: "interact_block" | "wait" | "place_blocks" | "fill"
	Action string `yaml:"action"`

	// --- interact_block 用 ---
	X           int               `yaml:"x"`
	Y           int               `yaml:"y"`
	Z           int               `yaml:"z"`
	TargetState map[string]string `yaml:"target_state"` // 目標とするブロック状態（省略時は無条件にインタラクト）

	// --- wait 用 ---
	Ms int `yaml:"ms"` // 待機時間（ミリ秒）

	// --- place_blocks 用 ---
	BlocksFile string           `yaml:"blocks_file"` // PlaceRequest形式のJSONファイルパス
	Blocks     *model.PlaceRequest `yaml:"blocks"`      // インラインで指定するPlaceRequest

	// --- fill 用 ---
	X1    int    `yaml:"x1"`
	Y1    int    `yaml:"y1"`
	Z1    int    `yaml:"z1"`
	X2    int    `yaml:"x2"`
	Y2    int    `yaml:"y2"`
	Z2    int    `yaml:"z2"`
	Block string `yaml:"block"` // fill で埋めるブロックID
}

// Assertion は1つのアサーション条件を表します。
type Assertion struct {
	X     int               `yaml:"x"`
	Y     int               `yaml:"y"`
	Z     int               `yaml:"z"`
	Block string            `yaml:"block"` // 省略可：省略時はブロックIDチェックをスキップ
	State map[string]string `yaml:"state"` // 部分一致チェック（指定したキーのみ検証）
}

// AssertionFailure は1つのアサーション失敗の詳細を表します。
type AssertionFailure struct {
	X        int
	Y        int
	Z        int
	Expected Assertion
	Actual   *model.BlockData // nilの場合はブロックが見つからなかった
	Reason   string
}

// TestResult は1つのテストケースの実行結果を表します。
type TestResult struct {
	Name     string
	Passed   bool
	Failures []AssertionFailure
	Error    string // セットアップやステップでエラーが発生した場合
}
