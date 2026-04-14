/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"mc-cli/internal/model"
	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mc-cli",
	Short: "Minecraft サーバー API への CLI ラッパー",
	Long: `Minecraft サーバー（Fabric）に組み込まれた HTTP API を通じて、
ワールド情報の取得、ブロックの配置、およびインタラクト操作を行うための CLI ツールです。
AI エージェントが解析しやすいよう、基本的に JSON 形式で結果を出力します。`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	serverURL string
	verbose   bool
)

func init() {
	// サーバーのベース URL を指定するための永続フラグ。デフォルトは localhost:8080
	rootCmd.PersistentFlags().StringVar(&serverURL, "url", "http://localhost:8080", "Minecraft サーバーの HTTP API ベース URL")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "冗長な形式で JSON を出力する")
}

// printJSON はデータを JSON 形式で標準出力します。
func printJSON(data any) {
	output := data
	if !verbose {
		output = tryCompact(data)
	}

	b, err := json.Marshal(output)
	if err != nil {
		printError(fmt.Sprintf("JSON エンコードエラー: %v", err))
		return
	}
	fmt.Println(string(b))
}

// tryCompact はデータを可能な限りコンパクトな形式に変換します。
func tryCompact(data any) any {
	switch v := data.(type) {
	case model.CommandResult:
		v.Data = tryCompact(v.Data)
		return v
	case []model.BlockData:
		res := make([]any, len(v))
		for i, b := range v {
			res[i] = b.ToCompact()
		}
		return res
	case [][]model.BlockData:
		res := make([]any, len(v))
		for i, list := range v {
			inner := make([]any, len(list))
			for j, b := range list {
				inner[j] = b.ToCompact()
			}
			res[i] = inner
		}
		return res
	case []model.AttachesData:
		res := make([]any, len(v))
		for i, a := range v {
			res[i] = a.ToCompact()
		}
		return res
	case []model.ConnectsData:
		res := make([]any, len(v))
		for i, c := range v {
			res[i] = c.ToCompact()
		}
		return res
	}
	return data
}

// printError はエラー情報を JSON 形式で標準出力し、プロセスを終了します。
func printError(msg string) {
	result := map[string]any{
		"success": false,
		"message": msg,
	}
	b, _ := json.Marshal(result)
	fmt.Println(string(b))
	os.Exit(1)
}


