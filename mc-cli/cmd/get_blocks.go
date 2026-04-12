package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"mc-cli/internal/model"
	"github.com/spf13/cobra"
)

var (
	x1, y1, z1 int
	x2, y2, z2 int
	interval   int
	count      int
)

var getBlocksCmd = &cobra.Command{
	Use:   "get-blocks",
	Short: "指定された座標範囲のブロック状況を取得する",
	Long:  `指定された (x1, y1, z1) から (x2, y2, z2) までの範囲に含まれるブロックの情報を取得します。`,
	Run: func(cmd *cobra.Command, args []string) {
		var allResults [][]model.BlockData

		for i := 0; i < count; i++ {
			url := fmt.Sprintf("%s/api/blocks?x1=%d&y1=%d&z1=%d&x2=%d&y2=%d&z2=%d",
				serverURL, x1, y1, z1, x2, y2, z2)

			resp, err := http.Get(url)
			if err != nil {
				printError(fmt.Sprintf("API リクエスト失敗 (回数 %d): %v", i+1, err))
				return
			}

			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				printError(fmt.Sprintf("レスポンス読み取り失敗 (回数 %d): %v", i+1, err))
				return
			}

			if resp.StatusCode != http.StatusOK {
				printError(fmt.Sprintf("API エラー (回数 %d, ステータス: %d): %s", i+1, resp.StatusCode, string(body)))
				return
			}

			var blocks []model.BlockData
			if err := json.Unmarshal(body, &blocks); err != nil {
				printError(fmt.Sprintf("JSON デコード失敗 (回数 %d): %v", i+1, err))
				return
			}

			allResults = append(allResults, blocks)

			if i < count-1 && interval > 0 {
				time.Sleep(time.Duration(interval) * 50 * time.Millisecond)
			}
		}

		printJSON(model.CommandResult{
			Success: true,
			Data:    allResults,
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
	getBlocksCmd.Flags().IntVar(&interval, "interval", 0, "ゲームチック間隔 (1チック=50ms)")
	getBlocksCmd.Flags().IntVar(&count, "count", 1, "実行回数")
}
