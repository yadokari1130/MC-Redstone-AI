package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
		var blocks []model.BlockData
		if err := json.Unmarshal(inputData, &blocks); err != nil {
			printError(fmt.Sprintf("JSON パース失敗: %v", err))
			return
		}

		// API リクエスト
		url := fmt.Sprintf("%s/api/blocks", serverURL)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(inputData))
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
			Message: "ブロックの配置に成功しました",
		})
	},
}

func init() {
	rootCmd.AddCommand(placeBlocksCmd)

	placeBlocksCmd.Flags().StringVar(&blocksInput, "blocks", "", "ブロックデータの JSON 文字列、または @file.json 形式のパス")
}
