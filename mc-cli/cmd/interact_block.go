package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"mc-cli/internal/model"
	"github.com/spf13/cobra"
)

var (
	xi, yi, zi                             int
	ix1, iy1, iz1, ix2, iy2, iz2, idelay int
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

		if cmd.Flags().Changed("x1") {
			if idelay > 0 {
				time.Sleep(time.Duration(idelay) * 50 * time.Millisecond)
			}

			blocksUrl := fmt.Sprintf("%s/api/blocks?x1=%d&y1=%d&z1=%d&x2=%d&y2=%d&z2=%d",
				serverURL, ix1, iy1, iz1, ix2, iy2, iz2)

			blocksResp, err := http.Get(blocksUrl)
			if err != nil {
				printError(fmt.Sprintf("ブロック取得リクエスト失敗: %v", err))
				return
			}
			defer blocksResp.Body.Close()

			blocksBody, err := io.ReadAll(blocksResp.Body)
			if err != nil {
				printError(fmt.Sprintf("ブロック取得レスポンス読み取り失敗: %v", err))
				return
			}

			if blocksResp.StatusCode != http.StatusOK {
				printError(fmt.Sprintf("ブロック取得 API エラー (ステータス: %d): %s", blocksResp.StatusCode, string(blocksBody)))
				return
			}

			var blocks []model.BlockData
			if err := json.Unmarshal(blocksBody, &blocks); err != nil {
				printError(fmt.Sprintf("JSON デコード失敗 (ブロック取得): %v", err))
				return
			}

			printJSON(model.CommandResult{
				Success: true,
				Data:    [][]model.BlockData{blocks},
			})
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

	interactBlockCmd.Flags().IntVar(&ix1, "x1", 0, "取得範囲の開始 X 座標")
	interactBlockCmd.Flags().IntVar(&iy1, "y1", 0, "取得範囲の開始 Y 座標")
	interactBlockCmd.Flags().IntVar(&iz1, "z1", 0, "取得範囲の開始 Z 座標")
	interactBlockCmd.Flags().IntVar(&ix2, "x2", 0, "取得範囲の終了 X 座標")
	interactBlockCmd.Flags().IntVar(&iy2, "y2", 0, "取得範囲の終了 Y 座標")
	interactBlockCmd.Flags().IntVar(&iz2, "z2", 0, "取得範囲の終了 Z 座標")
	interactBlockCmd.Flags().IntVar(&idelay, "delay", 0, "操作後の待機時間 (ゲームチック単位)")
}
