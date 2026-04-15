package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"

	"mc-cli/internal/model"
	"github.com/spf13/cobra"
)

var blocksInput string

var placeBlocksCmd = &cobra.Command{
	Use:   "place-blocks",
	Short: "ブロックを世界の中に配置する",
	Long:  `JSON 配列またはファイルパス（@file.json）からブロックデータを読み込み、Minecraft の世界に配置します。`,
	Run: func(cmd *cobra.Command, args []string) {
		var inputData []byte
		var err error

		if strings.HasPrefix(blocksInput, "@") {
			// ファイルから読み込む
			filePath := strings.TrimPrefix(blocksInput, "@")
			inputData, err = os.ReadFile(filePath)
			if err != nil {
				printError(fmt.Sprintf("ファイルの読み込み失敗: %v", err))
				return
			}
		} else {
			// 直接 JSON 文字列として扱う
			inputData = []byte(blocksInput)
		}

		if len(inputData) == 0 {
			printError("ブロックデータが指定されていません")
			return
		}

		// 有効な JSON かチェック
		var req model.PlaceRequest
		if err := json.Unmarshal(inputData, &req); err != nil {
			printError(fmt.Sprintf("JSON パース失敗: %v", err))
			return
		}

		// 1. attaches バリデーション
		if len(req.Attaches) > 0 {
			minX, minY, minZ := math.MaxInt32, math.MaxInt32, math.MaxInt32
			maxX, maxY, maxZ := math.MinInt32, math.MinInt32, math.MinInt32

			for _, a := range req.Attaches {
				if a.Base[0] < minX { minX = a.Base[0] }
				if a.Base[1] < minY { minY = a.Base[1] }
				if a.Base[2] < minZ { minZ = a.Base[2] }
				if a.Base[0] > maxX { maxX = a.Base[0] }
				if a.Base[1] > maxY { maxY = a.Base[1] }
				if a.Base[2] > maxZ { maxZ = a.Base[2] }
			}

			blocks, err := getBlocksRange(minX, minY, minZ, maxX, maxY, maxZ)
			if err != nil {
				printError(fmt.Sprintf("ベースブロックの取得失敗: %v", err))
				return
			}
			blockMap := make(map[string]string)
			for _, b := range blocks {
				key := fmt.Sprintf("%d,%d,%d", b.X, b.Y, b.Z)
				blockMap[key] = b.Block
			}

			// リクエスト内の blocks に含まれる予定のブロックも考慮する
			for _, b := range req.Blocks {
				key := fmt.Sprintf("%d,%d,%d", b.X, b.Y, b.Z)
				blockMap[key] = b.Block
			}

			for _, a := range req.Attaches {
				key := fmt.Sprintf("%d,%d,%d", a.Base[0], a.Base[1], a.Base[2])
				b, ok := blockMap[key]
				if !ok || strings.Contains(b, "air") {
					printError(fmt.Sprintf("ベースブロックが存在しません: x=%d, y=%d, z=%d", a.Base[0], a.Base[1], a.Base[2]))
					return
				}
			}
		}

		// 2. connects バリデーション
		for _, c := range req.Connects {
			distX := int(math.Abs(float64(c.From[0] - c.To[0])))
			distY := int(math.Abs(float64(c.From[1] - c.To[1])))
			distZ := int(math.Abs(float64(c.From[2] - c.To[2])))

			if (distX == 2 && distY == 0 && distZ == 0) ||
				(distX == 0 && distY == 2 && distZ == 0) ||
				(distX == 0 && distY == 0 && distZ == 2) {
				// OK
			} else {
				printError(fmt.Sprintf("connects の from と to の間がちょうど1マスではありません: from=(%d,%d,%d), to=(%d,%d,%d)", c.From[0], c.From[1], c.From[2], c.To[0], c.To[1], c.To[2]))
				return
			}
		}

		// 3. attaches 計算
		var attachesBlocks []model.BlockData
		for _, a := range req.Attaches {
			state := make(map[string]string)
			component := a.Component

			// ブロックの種類によってプロパティを使い分ける (faceプロパティを持つタイプ)
			isFaceType := strings.Contains(component, "lever") || strings.Contains(component, "button") || strings.Contains(component, "grindstone")
			// レッドストーントーチの判定
			isTorch := strings.Contains(component, "redstone_torch") || strings.Contains(component, "redstone_wall_torch")

			if a.Pos[1] > a.Base[1] {
				// 上面 (floor)
				if isFaceType {
					state["face"] = "floor"
					state["facing"] = "north" // デフォルトの向き
				} else if isTorch {
					component = "minecraft:redstone_torch"
					// 床置きトーチには facing プロパティはない
				} else {
					state["facing"] = "up"
				}
			} else if a.Pos[1] < a.Base[1] {
				// 下面 (ceiling)
				if isFaceType {
					state["face"] = "ceiling"
					state["facing"] = "north" // デフォルトの向き
				} else if isTorch {
					// レッドストーントーチは天井には設置できないが、デフォルトで床置きとして扱うか
					component = "minecraft:redstone_torch"
				} else {
					state["facing"] = "down"
				}
			} else {
				// 側面 (wall)
				facing := ""
				if a.Pos[0] > a.Base[0] {
					facing = "east"
				} else if a.Pos[0] < a.Base[0] {
					facing = "west"
				} else if a.Pos[2] > a.Base[2] {
					facing = "south"
				} else if a.Pos[2] < a.Base[2] {
					facing = "north"
				}

				if facing != "" {
					if isFaceType {
						state["face"] = "wall"
					} else if isTorch {
						component = "minecraft:redstone_wall_torch"
					}
					state["facing"] = facing
				}
			}

			attachesBlocks = append(attachesBlocks, model.BlockData{
				X:     a.Pos[0],
				Y:     a.Pos[1],
				Z:     a.Pos[2],
				Block: component,
				State: state,
			})
		}

		// 4. connects 計算
		var connectsBlocks []model.BlockData
		for _, c := range req.Connects {
			facing := ""
			if c.To[0] > c.From[0] {
				facing = "east"
			} else if c.To[0] < c.From[0] {
				facing = "west"
			} else if c.To[2] > c.From[2] {
				facing = "south"
			} else if c.To[2] < c.From[2] {
				facing = "north"
			} else if c.To[1] > c.From[1] {
				facing = "up"
			} else if c.To[1] < c.From[1] {
				facing = "down"
			}

			state := make(map[string]string)
			if facing != "" {
				state["facing"] = facing
			}

			connectsBlocks = append(connectsBlocks, model.BlockData{
				X:     (c.From[0] + c.To[0]) / 2,
				Y:     (c.From[1] + c.To[1]) / 2,
				Z:     (c.From[2] + c.To[2]) / 2,
				Block: c.Component,
				State: state,
			})
		}

		// 5. APIリクエスト送信
		if err := sendBlocks(req.Blocks); err != nil {
			printError(fmt.Sprintf("blocks の配置に失敗しました: %v", err))
			return
		}
		if err := sendBlocks(attachesBlocks); err != nil {
			printError(fmt.Sprintf("attaches の配置に失敗しました: %v", err))
			return
		}
		if err := sendBlocks(connectsBlocks); err != nil {
			printError(fmt.Sprintf("connects の配置に失敗しました: %v", err))
			return
		}

		// 6. ブロックアップデート
		var updateBlocks []model.BlockData
		updateBlocks = append(updateBlocks, req.Blocks...)
		updateBlocks = append(updateBlocks, attachesBlocks...)
		updateBlocks = append(updateBlocks, connectsBlocks...)

		if len(updateBlocks) > 0 {
			// 重複排除（同じ座標を何度もアップデートしないようにする）
			uniqueBlocks := make([]model.BlockData, 0)
			seen := make(map[string]bool)
			for _, b := range updateBlocks {
				key := fmt.Sprintf("%d,%d,%d", b.X, b.Y, b.Z)
				if !seen[key] {
					seen[key] = true
					uniqueBlocks = append(uniqueBlocks, model.BlockData{X: b.X, Y: b.Y, Z: b.Z})
				}
			}

			if err := sendUpdateBlocks(uniqueBlocks); err != nil {
				printError(fmt.Sprintf("ブロックアップデートに失敗しました: %v", err))
				return
			}
		}

		printJSON(model.CommandResult{
			Success: true,
			Message: "ブロックの配置に成功しました",
		})
	},
}

func getBlocksRange(x1, y1, z1, x2, y2, z2 int) ([]model.BlockData, error) {
	urlStr := fmt.Sprintf("%s/api/blocks?x1=%d&y1=%d&z1=%d&x2=%d&y2=%d&z2=%d", serverURL, x1, y1, z1, x2, y2, z2)
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

	return blocks, nil
}

func sendBlocks(blocks []model.BlockData) error {
	if len(blocks) == 0 {
		return nil
	}

	data, err := json.Marshal(blocks)
	if err != nil {
		return err
	}

	urlStr := fmt.Sprintf("%s/api/blocks", serverURL)
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

func sendUpdateBlocks(blocks []model.BlockData) error {
	if len(blocks) == 0 {
		return nil
	}

	data, err := json.Marshal(blocks)
	if err != nil {
		return err
	}

	urlStr := fmt.Sprintf("%s/api/update-blocks", serverURL)
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

func init() {
	rootCmd.AddCommand(placeBlocksCmd)

	placeBlocksCmd.Flags().StringVar(&blocksInput, "blocks", "", "ブロックデータの JSON 文字列、または @file.json 形式のパス (PlaceRequest形式)")
}
