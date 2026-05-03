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
	gbPos1Str, gbPos2Str string
	interval             int
	count                int
	includeEntities      bool
)

var getBlocksCmd = &cobra.Command{
	Use:   "get-blocks",
	Short: "指定された座標範囲のブロック状況を取得する",
	Long:  `指定された --pos1 から --pos2 までの範囲に含まれるブロックの情報を取得します。`,
	Run: func(cmd *cobra.Command, args []string) {
		pos1, err := parsePos(gbPos1Str)
		if err != nil {
			printError(err.Error())
		}
		pos2, err := parsePos(gbPos2Str)
		if err != nil {
			printError(err.Error())
		}

		var allResults [][]model.BlockData
		var allResultsWithEntities []model.BlocksAndEntities

		for i := 0; i < count; i++ {
			url := fmt.Sprintf("%s/api/blocks?x1=%d&y1=%d&z1=%d&x2=%d&y2=%d&z2=%d&include_entities=%t",
				serverURL, pos1[0], pos1[1], pos1[2], pos2[0], pos2[1], pos2[2], includeEntities)

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

			if includeEntities {
				// まず BlocksAndEntities 形式でパースを試みる
				var result model.BlocksAndEntities
				if err := json.Unmarshal(body, &result); err != nil {
					// パース失敗した場合、レスポンスが配列（旧形式）かどうか確認
					if len(body) > 0 && body[0] == '[' {
						printError(fmt.Sprintf(
							"サーバーがエンティティ取得に対応していないようです。"+
							"Modの再ビルドとMinecraftサーバーの再起動が必要です。"+
							"(APIレスポンスが配列形式です: %s)", string(body)))
						return
					}
					printError(fmt.Sprintf("JSON デコード失敗 (回数 %d): %v", i+1, err))
					return
				}
				allResultsWithEntities = append(allResultsWithEntities, result)
			} else {
				var blocks []model.BlockData
				if err := json.Unmarshal(body, &blocks); err != nil {
					printError(fmt.Sprintf("JSON デコード失敗 (回数 %d): %v", i+1, err))
					return
				}
				allResults = append(allResults, blocks)
			}

			if i < count-1 && interval > 0 {
				time.Sleep(time.Duration(interval) * 50 * time.Millisecond)
			}
		}

		if includeEntities {
			printJSON(model.CommandResult{
				Success: true,
				Data:    allResultsWithEntities,
			})
		} else {
			printJSON(model.CommandResult{
				Success: true,
				Data:    allResults,
			})
		}
	},
}

func init() {
	rootCmd.AddCommand(getBlocksCmd)

	getBlocksCmd.Flags().StringVar(&gbPos1Str, "pos1", "0,0,0", "開始座標 [x,y,z]")
	getBlocksCmd.Flags().StringVar(&gbPos2Str, "pos2", "0,0,0", "終了座標 [x,y,z]")
	getBlocksCmd.Flags().IntVar(&interval, "interval", 0, "ゲームチック間隔 (1チック=50ms)")
	getBlocksCmd.Flags().IntVar(&count, "count", 1, "実行回数")
	getBlocksCmd.Flags().BoolVar(&includeEntities, "include-entities", false, "エンティティ情報も取得する")
}
