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

var (
	dropPosStr string
	itemsInput string
)

var dropItemsCmd = &cobra.Command{
	Use:   "drop-items",
	Short: "アイテムを指定した座標にドロップする",
	Long:  `JSON 配列またはファイルパス（@file.json）からアイテムデータを読み込み、指定した --pos 座標にドロップします。`,
	Run: func(cmd *cobra.Command, args []string) {
		pos, err := parsePosF(dropPosStr)
		if err != nil {
			printError(err.Error())
		}
		var items []model.ItemInfo
		// ... (reading itemsInput logic remains the same)
		if itemsInput != "" {
			var inputData []byte
			if strings.HasPrefix(itemsInput, "@") {
				// ファイルから読み込む
				filePath := strings.TrimPrefix(itemsInput, "@")
				inputData, err = os.ReadFile(filePath)
				if err != nil {
					printError(fmt.Sprintf("ファイルの読み込み失敗: %v", err))
					return
				}
			} else {
				// 直接 JSON 文字列として扱う
				inputData = []byte(itemsInput)
			}

			if err := json.Unmarshal(inputData, &items); err != nil {
				printError(fmt.Sprintf("アイテムデータの JSON パース失敗: %v", err))
				return
			}
		}

		if len(items) == 0 {
			printError("ドロップするアイテムが指定されていません。--items フラグで指定してください。")
			return
		}

		req := model.DropItemsRequest{
			X:     pos[0],
			Y:     pos[1],
			Z:     pos[2],
			Items: items,
		}
// ... (rest of Run same)
		jsonData, err := json.Marshal(req)
		if err != nil {
			printError(fmt.Sprintf("JSON エンコード失敗: %v", err))
			return
		}

		url := fmt.Sprintf("%s/api/drop-items", serverURL)
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
			Message: "アイテムのドロップに成功しました",
		})
	},
}

func init() {
	rootCmd.AddCommand(dropItemsCmd)

	dropItemsCmd.Flags().StringVar(&dropPosStr, "pos", "0,0,0", "ドロップ先の座標 [x,y,z]")
	dropItemsCmd.Flags().StringVar(&itemsInput, "items", "", "アイテム情報の JSON 配列（[{\"id\":\"...\",\"amount\":...}, ...]）、または @file.json 形式のパス")
}
