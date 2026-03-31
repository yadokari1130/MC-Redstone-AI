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

var (
	xi, yi, zi int
)

var interactBlockCmd = &cobra.Command{
	Use:   "interact-block",
	Short: "指定した座標のブロックを操作する",
	Long:  `指定した (x, y, z) 座標にあるブロック（レバー、ボタンなど）に対してインタラクト操作を実行します。`,
	Run: func(cmd *cobra.Command, args []string) {
		req := model.InteractionRequest{
			X: xi,
			Y: yi,
			Z: zi,
		}

		jsonData, err := json.Marshal(req)
		if err != nil {
			printError(fmt.Sprintf("JSON エンコード失敗: %v", err))
			return
		}

		url := fmt.Sprintf("%s/api/interact", serverURL)
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
			Message: "ブロックの操作に成功しました",
		})
	},
}

func init() {
	rootCmd.AddCommand(interactBlockCmd)

	interactBlockCmd.Flags().IntVar(&xi, "x", 0, "対象の X 座標")
	interactBlockCmd.Flags().IntVar(&yi, "y", 0, "対象の Y 座標")
	interactBlockCmd.Flags().IntVar(&zi, "z", 0, "対象の Z 座標")
}
