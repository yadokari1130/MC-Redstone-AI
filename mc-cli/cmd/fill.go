package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"mc-cli/internal/model"
	"github.com/spf13/cobra"
)

var fillStateInput string

var fillCmd = &cobra.Command{
	Use:   "fill <x1> <y1> <z1> <x2> <y2> <z2> <block>",
	Short: "指定された範囲を特定のブロックで埋める",
	Long:  `座標 (x1, y1, z1) から (x2, y2, z2) までの矩形範囲を、指定されたブロック ID で埋めます。`,
	Args:  cobra.ExactArgs(7),
	Run: func(cmd *cobra.Command, args []string) {
		// 引数のパース
		x1, err1 := strconv.Atoi(args[0])
		y1, err2 := strconv.Atoi(args[1])
		z1, err3 := strconv.Atoi(args[2])
		x2, err4 := strconv.Atoi(args[3])
		y2, err5 := strconv.Atoi(args[4])
		z2, err6 := strconv.Atoi(args[5])
		blockID := args[6]

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
			printError("座標は整数で指定してください")
			return
		}

		// ブロック状態のパース
		var state map[string]string
		if fillStateInput != "" {
			if err := json.Unmarshal([]byte(fillStateInput), &state); err != nil {
				printError(fmt.Sprintf("ブロック状態 (state) の JSON パース失敗: %v", err))
				return
			}
		}

		// 範囲の最小値・最大値を算出
		minX, maxX := minMax(x1, x2)
		minY, maxY := minMax(y1, y2)
		minZ, maxZ := minMax(z1, z2)

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
