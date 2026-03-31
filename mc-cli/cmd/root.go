/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

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

var serverURL string

func init() {
	// サーバーのベース URL を指定するための永続フラグ。デフォルトは localhost:8080
	rootCmd.PersistentFlags().StringVar(&serverURL, "url", "http://localhost:8080", "Minecraft サーバーの HTTP API ベース URL")
}

// printJSON はデータを JSON 形式で標準出力します。
func printJSON(data any) {
	b, err := json.Marshal(data)
	if err != nil {
		printError(fmt.Sprintf("JSON エンコードエラー: %v", err))
		return
	}
	fmt.Println(string(b))
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


