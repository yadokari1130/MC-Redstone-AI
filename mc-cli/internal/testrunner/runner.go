package testrunner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"time"

	"mc-cli/internal/model"
)

// Runner はテストの実行を管理します。
type Runner struct {
	ServerURL string
}

// NewRunner は新しいRunnerを生成します。
func NewRunner(serverURL string) *Runner {
	return &Runner{ServerURL: serverURL}
}

// RunTest は1つのTestCaseを実行し、TestResultを返します。
func (r *Runner) RunTest(tc TestCase) TestResult {
	result := TestResult{Name: tc.Name}

	// 1. セットアップ（回路配置）
	if err := r.runSetup(tc.Setup); err != nil {
		result.Error = fmt.Sprintf("セットアップ失敗: %v", err)
		return result
	}

	// 2. ステップの実行
	for i, step := range tc.Steps {
		if err := r.runStep(step); err != nil {
			result.Error = fmt.Sprintf("ステップ %d (%s) 失敗: %v", i+1, step.Action, err)
			return result
		}
	}

	// 3. アサーション
	if len(tc.Assert) > 0 {
		// アサーション対象の最小バウンディングボックスを計算してブロックを取得
		blockMap, err := r.fetchBlocksForAssertions(tc.Assert)
		if err != nil {
			result.Error = fmt.Sprintf("アサーション用ブロック取得失敗: %v", err)
			return result
		}

		failures := checkAssertions(tc.Assert, blockMap)
		result.Failures = failures
		result.Passed = len(failures) == 0
	} else {
		// アサーションが定義されていない場合は合格扱い
		result.Passed = true
	}

	return result
}

// runSetup はセットアップ（回路配置）を実行します。
func (r *Runner) runSetup(setup Setup) error {
	if setup.BlocksFile == "" && setup.Blocks == nil {
		// セットアップなし
		return nil
	}

	var req model.PlaceRequest
	if setup.BlocksFile != "" {
		data, err := os.ReadFile(setup.BlocksFile)
		if err != nil {
			return fmt.Errorf("ファイル読み込み失敗 (%s): %v", setup.BlocksFile, err)
		}
		if err := json.Unmarshal(data, &req); err != nil {
			return fmt.Errorf("JSONパース失敗 (%s): %v", setup.BlocksFile, err)
		}
	} else {
		req = *setup.Blocks
	}

	return r.placeRequest(req)
}

// runStep は1つのステップを実行します。
func (r *Runner) runStep(step Step) error {
	switch step.Action {
	case "interact_block":
		return r.runInteractBlock(step)
	case "wait":
		if step.Ms > 0 {
			time.Sleep(time.Duration(step.Ms) * time.Millisecond)
		}
		return nil
	case "place_blocks":
		return r.runPlaceBlocks(step)
	case "fill":
		return r.runFill(step)
	default:
		return fmt.Errorf("不明なアクション: %s", step.Action)
	}
}

// runInteractBlock はinteract_blockステップを実行します。
// target_stateが指定されている場合は現在の状態を確認し、必要な場合のみインタラクトします。
func (r *Runner) runInteractBlock(step Step) error {
	if len(step.TargetState) > 0 {
		// 現在の状態を取得
		block, err := r.fetchSingleBlock(step.X, step.Y, step.Z)
		if err != nil {
			return fmt.Errorf("ブロック状態取得失敗 (%d,%d,%d): %v", step.X, step.Y, step.Z, err)
		}

		// 目標状態と現在の状態が一致しているか確認
		if block != nil && matchesTargetState(block.State, step.TargetState) {
			// 既に目標状態なので何もしない
			return nil
		}
	}

	// インタラクト実行
	return r.interact(step.X, step.Y, step.Z)
}

// runPlaceBlocks はplace_blocksステップを実行します。
func (r *Runner) runPlaceBlocks(step Step) error {
	var req model.PlaceRequest
	if step.BlocksFile != "" {
		data, err := os.ReadFile(step.BlocksFile)
		if err != nil {
			return fmt.Errorf("ファイル読み込み失敗 (%s): %v", step.BlocksFile, err)
		}
		if err := json.Unmarshal(data, &req); err != nil {
			return fmt.Errorf("JSONパース失敗 (%s): %v", step.BlocksFile, err)
		}
	} else if step.Blocks != nil {
		req = *step.Blocks
	} else {
		return fmt.Errorf("place_blocks: blocks_file または blocks のどちらかを指定してください")
	}
	return r.placeRequest(req)
}

// runFill はfillステップを実行します。
func (r *Runner) runFill(step Step) error {
	if step.Block == "" {
		return fmt.Errorf("fill: block を指定してください")
	}

	minX, maxX := minMax(step.X1, step.X2)
	minY, maxY := minMax(step.Y1, step.Y2)
	minZ, maxZ := minMax(step.Z1, step.Z2)

	var blocks []model.BlockData
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			for z := minZ; z <= maxZ; z++ {
				blocks = append(blocks, model.BlockData{
					X:     x,
					Y:     y,
					Z:     z,
					Block: step.Block,
				})
			}
		}
	}

	return r.sendBlocks(blocks)
}

// matchesTargetState はブロックの現在の状態が目標状態と一致するかを確認します。
func matchesTargetState(current, target map[string]string) bool {
	for k, v := range target {
		if current[k] != v {
			return false
		}
	}
	return true
}

// fetchSingleBlock は1つのブロックの状態を取得します。
// ブロックが存在しない（空気）場合はnilを返します。
func (r *Runner) fetchSingleBlock(x, y, z int) (*model.BlockData, error) {
	urlStr := fmt.Sprintf("%s/api/blocks?x1=%d&y1=%d&z1=%d&x2=%d&y2=%d&z2=%d",
		r.ServerURL, x, y, z, x, y, z)
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var blocks []model.BlockData
	if err := json.Unmarshal(body, &blocks); err != nil {
		return nil, err
	}

	if len(blocks) == 0 {
		return nil, nil // 空気ブロック
	}
	return &blocks[0], nil
}

// fetchBlocksForAssertions はアサーション対象のブロックをまとめて取得します。
// 全アサーション座標を包む最小バウンディングボックスで一括取得します。
func (r *Runner) fetchBlocksForAssertions(assertions []Assertion) (map[string]model.BlockData, error) {
	if len(assertions) == 0 {
		return nil, nil
	}

	// バウンディングボックスを計算
	minX, minY, minZ := math.MaxInt32, math.MaxInt32, math.MaxInt32
	maxX, maxY, maxZ := math.MinInt32, math.MinInt32, math.MinInt32
	for _, a := range assertions {
		if a.X < minX { minX = a.X }
		if a.Y < minY { minY = a.Y }
		if a.Z < minZ { minZ = a.Z }
		if a.X > maxX { maxX = a.X }
		if a.Y > maxY { maxY = a.Y }
		if a.Z > maxZ { maxZ = a.Z }
	}

	urlStr := fmt.Sprintf("%s/api/blocks?x1=%d&y1=%d&z1=%d&x2=%d&y2=%d&z2=%d",
		r.ServerURL, minX, minY, minZ, maxX, maxY, maxZ)
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var blocks []model.BlockData
	if err := json.Unmarshal(body, &blocks); err != nil {
		return nil, err
	}

	blockMap := make(map[string]model.BlockData)
	for _, b := range blocks {
		blockMap[blockKey(b.X, b.Y, b.Z)] = b
	}
	return blockMap, nil
}

// interact は指定座標のブロックをインタラクトします。
func (r *Runner) interact(x, y, z int) error {
	req := model.InteractionRequest{X: x, Y: y, Z: z}
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	urlStr := fmt.Sprintf("%s/api/interact", r.ServerURL)
	resp, err := http.Post(urlStr, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

// placeRequest はPlaceRequestをAPIに送信してブロックを配置します。
func (r *Runner) placeRequest(req model.PlaceRequest) error {
	// blocksを配置
	if err := r.sendBlocks(req.Blocks); err != nil {
		return fmt.Errorf("blocks配置失敗: %v", err)
	}

	// attaches を変換して配置
	attachesBlocks := resolveAttaches(req.Attaches)
	if err := r.sendBlocks(attachesBlocks); err != nil {
		return fmt.Errorf("attaches配置失敗: %v", err)
	}

	// connects を変換して配置
	connectsBlocks := resolveConnects(req.Connects)
	if err := r.sendBlocks(connectsBlocks); err != nil {
		return fmt.Errorf("connects配置失敗: %v", err)
	}

	return nil
}

// sendBlocks はブロック配列をAPIに送信します。
func (r *Runner) sendBlocks(blocks []model.BlockData) error {
	if len(blocks) == 0 {
		return nil
	}

	data, err := json.Marshal(blocks)
	if err != nil {
		return err
	}

	urlStr := fmt.Sprintf("%s/api/blocks", r.ServerURL)
	resp, err := http.Post(urlStr, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

// blockKey は座標を "x,y,z" 形式の文字列キーに変換します。
func blockKey(x, y, z int) string {
	return fmt.Sprintf("%d,%d,%d", x, y, z)
}

// minMax はaとbの小さい方と大きい方を返します。
func minMax(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}
