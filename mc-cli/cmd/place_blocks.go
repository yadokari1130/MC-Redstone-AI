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
				if a.BaseX < minX { minX = a.BaseX }
				if a.BaseY < minY { minY = a.BaseY }
				if a.BaseZ < minZ { minZ = a.BaseZ }
				if a.BaseX > maxX { maxX = a.BaseX }
				if a.BaseY > maxY { maxY = a.BaseY }
				if a.BaseZ > maxZ { maxZ = a.BaseZ }
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
				key := fmt.Sprintf("%d,%d,%d", a.BaseX, a.BaseY, a.BaseZ)
				b, ok := blockMap[key]
				if !ok || strings.Contains(b, "air") {
					printError(fmt.Sprintf("ベースブロックが存在しません: x=%d, y=%d, z=%d", a.BaseX, a.BaseY, a.BaseZ))
					return
				}
			}
		}

		// 2. connects バリデーション
		for _, c := range req.Connects {
			distX := int(math.Abs(float64(c.FromX - c.ToX)))
			distY := int(math.Abs(float64(c.FromY - c.ToY)))
			distZ := int(math.Abs(float64(c.FromZ - c.ToZ)))

			if (distX == 2 && distY == 0 && distZ == 0) ||
				(distX == 0 && distY == 2 && distZ == 0) ||
				(distX == 0 && distY == 0 && distZ == 2) {
				// OK
			} else {
				printError(fmt.Sprintf("connects の from と to の間がちょうど1マスではありません: from=(%d,%d,%d), to=(%d,%d,%d)", c.FromX, c.FromY, c.FromZ, c.ToX, c.ToY, c.ToZ))
				return
			}
		}

		// 3. attaches 計算
		var attachesBlocks []model.BlockData
		for _, a := range req.Attaches {
			facing := ""
			if a.ComponentX > a.BaseX {
				facing = "east"
			} else if a.ComponentX < a.BaseX {
				facing = "west"
			} else if a.ComponentZ > a.BaseZ {
				facing = "south"
			} else if a.ComponentZ < a.BaseZ {
				facing = "north"
			} else if a.ComponentY > a.BaseY {
				facing = "up"
			} else if a.ComponentY < a.BaseY {
				facing = "down"
			}

			state := make(map[string]string)
			if facing != "" {
				state["facing"] = facing
			}

			attachesBlocks = append(attachesBlocks, model.BlockData{
				X:     a.ComponentX,
				Y:     a.ComponentY,
				Z:     a.ComponentZ,
				Block: a.Component,
				State: state,
			})
		}

		// 4. connects 計算
		var connectsBlocks []model.BlockData
		for _, c := range req.Connects {
			facing := ""
			if c.ToX > c.FromX {
				facing = "east"
			} else if c.ToX < c.FromX {
				facing = "west"
			} else if c.ToZ > c.FromZ {
				facing = "south"
			} else if c.ToZ < c.FromZ {
				facing = "north"
			} else if c.ToY > c.FromY {
				facing = "up"
			} else if c.ToY < c.FromY {
				facing = "down"
			}

			state := make(map[string]string)
			if facing != "" {
				state["facing"] = facing
			}

			connectsBlocks = append(connectsBlocks, model.BlockData{
				X:     (c.FromX + c.ToX) / 2,
				Y:     (c.FromY + c.ToY) / 2,
				Z:     (c.FromZ + c.ToZ) / 2,
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

func init() {
	rootCmd.AddCommand(placeBlocksCmd)

	placeBlocksCmd.Flags().StringVar(&blocksInput, "blocks", "", "ブロックデータの JSON 文字列、または @file.json 形式のパス (PlaceRequest形式)")
}
