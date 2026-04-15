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
	invPosStr     string
	invItemsInput string
)

var setInventoryCmd = &cobra.Command{
	Use:   "set-inventory",
	Short: "ブロックのインベントリにアイテムをセットする",
	Long:  `指定された --pos 座標にあるブロックのインベントリにアイテムをセットします。既存のアイテムは消去されます。`,
	Run: func(cmd *cobra.Command, args []string) {
		pos, err := parsePos(invPosStr)
		if err != nil {
			printError(err.Error())
		}
		var items []model.ItemInfo
		// ... (reading itemsInput logic)
		if invItemsInput != "" {
			var inputData []byte
			if strings.HasPrefix(invItemsInput, "@") {
				// ファイルから読み込む
				filePath := strings.TrimPrefix(invItemsInput, "@")
				inputData, err = os.ReadFile(filePath)
				if err != nil {
					printError(fmt.Sprintf("ファイルの読み込み失敗: %v", err))
					return
				}
			} else {
				// 直接 JSON 文字列として扱う
				inputData = []byte(invItemsInput)
			}

			if err := json.Unmarshal(inputData, &items); err != nil {
				printError(fmt.Sprintf("アイテムデータの JSON パース失敗: %v", err))
				return
			}
		}

		// 空配列の場合はアイテム削除なので、items が空でもリクエストを送る
		// ただし、もし --items が指定されなかった場合は明示的に空として扱う

		req := model.InventoryRequest{
			X:     pos[0],
			Y:     pos[1],
			Z:     pos[2],
			Items: items,
		}
// ... (rest same)
		jsonData, err := json.Marshal(req)
		if err != nil {
			printError(fmt.Sprintf("JSON エンコード失敗: %v", err))
			return
		}

		url := fmt.Sprintf("%s/api/inventory", serverURL)
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

		if resp.StatusCode != http.StatusOK && resp.StatusCode != 207 {
			printError(fmt.Sprintf("API エラー (ステータス: %d): %s", resp.StatusCode, string(body)))
			return
		}

		printJSON(model.CommandResult{
			Success: true,
			Message: string(body),
		})
	},
}

func init() {
	rootCmd.AddCommand(setInventoryCmd)

	setInventoryCmd.Flags().StringVar(&invPosStr, "pos", "0,0,0", "対象ブロックの座標 [x,y,z]")
	setInventoryCmd.Flags().StringVar(&invItemsInput, "items", "[]", "アイテム情報の JSON 配列（[{\"id\":\"...\",\"amount\":...}, ...]）、または @file.json 形式のパス")
}
