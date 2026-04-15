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
	ibPosStr, ibPos1Str, ibPos2Str string
	idelay                         int
)

var interactBlockCmd = &cobra.Command{
	Use:   "interact-block",
	Short: "指定した座標のブロックを操作する",
	Long:  `指定した --pos 座標にあるブロック（レバー、ボタンなど）に対してインタラクト操作を実行します。`,
	Run: func(cmd *cobra.Command, args []string) {
		pos, err := parsePos(ibPosStr)
		if err != nil {
			printError(err.Error())
		}

		req := model.InteractionRequest{
			X: pos[0],
			Y: pos[1],
			Z: pos[2],
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

		if cmd.Flags().Changed("pos1") {
			if idelay > 0 {
				time.Sleep(time.Duration(idelay) * 50 * time.Millisecond)
			}

			pos1, err := parsePos(ibPos1Str)
			if err != nil {
				printError(err.Error())
			}
			pos2, err := parsePos(ibPos2Str)
			if err != nil {
				printError(err.Error())
			}

			blocksUrl := fmt.Sprintf("%s/api/blocks?x1=%d&y1=%d&z1=%d&x2=%d&y2=%d&z2=%d",
				serverURL, pos1[0], pos1[1], pos1[2], pos2[0], pos2[1], pos2[2])

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

	interactBlockCmd.Flags().StringVar(&ibPosStr, "pos", "0,0,0", "対象の座標 [x,y,z]")

	interactBlockCmd.Flags().StringVar(&ibPos1Str, "pos1", "0,0,0", "取得範囲の開始座標 [x,y,z]")
	interactBlockCmd.Flags().StringVar(&ibPos2Str, "pos2", "0,0,0", "取得範囲の終了座標 [x,y,z]")
	interactBlockCmd.Flags().IntVar(&idelay, "delay", 0, "操作後の待機時間 (ゲームチック単位)")
}
