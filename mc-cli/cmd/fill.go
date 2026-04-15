package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"mc-cli/internal/model"
	"github.com/spf13/cobra"
)

var fillStateInput string

var fillCmd = &cobra.Command{
	Use:   "fill <pos1> <pos2> <block>",
	Short: "指定された範囲を特定のブロックで埋める",
	Long:  `座標 pos1 (x,y,z) から pos2 (x,y,z) までの矩形範囲を、指定されたブロック ID で埋めます。`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		// 引数のパース
		pos1, err := parsePos(args[0])
		if err != nil {
			printError(err.Error())
		}
		pos2, err := parsePos(args[1])
		if err != nil {
			printError(err.Error())
		}
		blockID := args[2]

		// ブロック状態のパース
		var state map[string]string
		if fillStateInput != "" {
			if err := json.Unmarshal([]byte(fillStateInput), &state); err != nil {
				printError(fmt.Sprintf("ブロック状態 (state) の JSON パース失敗: %v", err))
				return
			}
		}

		// 範囲の最小値・最大値を算出
		minX, maxX := minMax(pos1[0], pos2[0])
		minY, maxY := minMax(pos1[1], pos2[1])
		minZ, maxZ := minMax(pos1[2], pos2[2])

		// ブロックデータの生成
		var blocks []model.BlockData
		for x := minX; x <= maxX; x++ {
			for y := minY; y <= maxY; y++ {
				for z := minZ; z <= maxZ; z++ {
					blocks = append(blocks, model.BlockData{
						X:     x,
						Y:     y,
						Z:     z,
						Block: blockID,
						State: state,
					})
				}
			}
		}

		if len(blocks) == 0 {
			printError("配置対象のブロックがありません")
			return
		}

		// JSON エンコード
		jsonData, err := json.Marshal(blocks)
		if err != nil {
			printError(fmt.Sprintf("JSON エンコード失敗: %v", err))
			return
		}

		// API リクエスト
		url := fmt.Sprintf("%s/api/blocks", serverURL)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			printError(fmt.Sprintf("API リクエスト失敗: %v", err))
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			printError(fmt.Sprintf("レスポンス読み取り失敗: %v", err))
			return
		}

		if resp.StatusCode != http.StatusOK {
			printError(fmt.Sprintf("API エラー (ステータス: %d): %s", resp.StatusCode, string(body)))
			return
		}

		printJSON(model.CommandResult{
			Success: true,
			Message: fmt.Sprintf("%d 個のブロックを配置しました", len(blocks)),
		})
	},
}

func minMax(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}

func init() {
	rootCmd.AddCommand(fillCmd)

	fillCmd.Flags().StringVar(&fillStateInput, "state", "", "ブロックの状態を指定する JSON 文字列 (例: '{\"facing\":\"north\"}')")
}
