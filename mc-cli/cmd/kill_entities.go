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
	kePos1Str, kePos2Str string
	keType              string
)

var killEntitiesCmd = &cobra.Command{
	Use:   "kill-entities",
	Short: "指定範囲のエンティティを削除する",
	Long:  `指定された --pos1 から --pos2 までの範囲に含まれるエンティティを削除します。--type でエンティティタイプを指定できます。`,
	Run: func(cmd *cobra.Command, args []string) {
		pos1, err := parsePos(kePos1Str)
		if err != nil {
			printError(err.Error())
		}
		pos2, err := parsePos(kePos2Str)
		if err != nil {
			printError(err.Error())
		}

		req := struct {
			X1   int    `json:"x1"`
			Y1   int    `json:"y1"`
			Z1   int    `json:"z1"`
			X2   int    `json:"x2"`
			Y2   int    `json:"y2"`
			Z2   int    `json:"z2"`
			Type string `json:"type,omitempty"`
		}{
			X1:   pos1[0],
			Y1:   pos1[1],
			Z1:   pos1[2],
			X2:   pos2[0],
			Y2:   pos2[1],
			Z2:   pos2[2],
			Type: keType,
		}

		jsonData, err := json.Marshal(req)
		if err != nil {
			printError(fmt.Sprintf("JSON エンコード失敗: %v", err))
			return
		}

		url := fmt.Sprintf("%s/api/kill-entities", serverURL)
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

		// レスポンスから削除件数を取得
		var result model.CommandResult
		if err := json.Unmarshal(body, &result); err != nil {
			// JSONでない場合はテキストとして返す
			printJSON(model.CommandResult{
				Success: true,
				Message: string(body),
			})
			return
		}

		printJSON(result)
	},
}

func init() {
	rootCmd.AddCommand(killEntitiesCmd)

	killEntitiesCmd.Flags().StringVar(&kePos1Str, "pos1", "0,0,0", "開始座標 [x,y,z]")
	killEntitiesCmd.Flags().StringVar(&kePos2Str, "pos2", "0,0,0", "終了座標 [x,y,z]")
	killEntitiesCmd.Flags().StringVar(&keType, "type", "", "削除対象のエンティティタイプ（例: minecraft:boat）。未指定時はすべてのエンティティを対象とする")
}
