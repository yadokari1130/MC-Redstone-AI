package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"mc-cli/internal/model"
	"github.com/spf13/cobra"
)

var (
	x1, y1, z1 int
	x2, y2, z2 int
)

var getBlocksCmd = &cobra.Command{
	Use:   "get-blocks",
	Short: "指定された座標範囲のブロック状況を取得する",
	Long:  `指定された (x1, y1, z1) から (x2, y2, z2) までの範囲に含まれるブロックの情報を取得します。`,
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("%s/api/blocks?x1=%d&y1=%d&z1=%d&x2=%d&y2=%d&z2=%d", 
			serverURL, x1, y1, z1, x2, y2, z2)
		
		resp, err := http.Get(url)
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

		var blocks []model.BlockData
		if err := json.Unmarshal(body, &blocks); err != nil {
			printError(fmt.Sprintf("JSON デコード失敗: %v", err))
			return
		}

		printJSON(model.CommandResult{
			Success: true,
			Data:    blocks,
		})
	},
}

func init() {
	rootCmd.AddCommand(getBlocksCmd)

	getBlocksCmd.Flags().IntVar(&x1, "x1", 0, "開始 X 座標")
	getBlocksCmd.Flags().IntVar(&y1, "y1", 0, "開始 Y 座標")
	getBlocksCmd.Flags().IntVar(&z1, "z1", 0, "開始 Z 座標")
	getBlocksCmd.Flags().IntVar(&x2, "x2", 0, "終了 X 座標")
	getBlocksCmd.Flags().IntVar(&y2, "y2", 0, "終了 Y 座標")
	getBlocksCmd.Flags().IntVar(&z2, "z2", 0, "終了 Z 座標")
}
