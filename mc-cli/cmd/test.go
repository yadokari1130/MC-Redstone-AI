package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mc-cli/internal/testrunner"

	"gopkg.in/yaml.v3"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test <file_or_dir>",
	Short: "テストファイルを実行する",
	Long: `指定した YAML テストファイル（*_test.yaml / *_test.yml）を実行し、
回路の動作を検証します。ディレクトリを指定した場合は配下の
テストファイルをすべてスキャンして実行します。`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]

		// テストファイルの収集
		files, err := collectTestFiles(target)
		if err != nil {
			printError(fmt.Sprintf("テストファイルの収集失敗: %v", err))
			return
		}
		if len(files) == 0 {
			fmt.Println("テストファイルが見つかりませんでした。")
			return
		}

		runner := testrunner.NewRunner(serverURL)

		totalPassed := 0
		totalFailed := 0
		totalError := 0

		for _, filePath := range files {
			fmt.Printf("\n=== %s ===\n", filePath)

			// YAMLファイルのパース
			tf, err := parseTestFile(filePath)
			if err != nil {
				fmt.Printf("  ⚠️  ファイルのパース失敗: %v\n", err)
				totalError++
				continue
			}

			if tf.Name != "" {
				fmt.Printf("テスト名: %s\n", tf.Name)
			}
			if tf.Description != "" {
				fmt.Printf("説明: %s\n", tf.Description)
			}
			fmt.Println()

			passed := 0
			failed := 0
			errorCount := 0

			for _, tc := range tf.Tests {
				result := runner.RunTest(tc)

				if result.Error != "" {
					fmt.Printf("  ⚠️  %s\n", result.Name)
					fmt.Printf("     エラー: %s\n", result.Error)
					errorCount++
					totalError++
				} else if result.Passed {
					fmt.Printf("  ✅ %s\n", result.Name)
					passed++
					totalPassed++
				} else {
					fmt.Printf("  ❌ %s\n", result.Name)
					for _, f := range result.Failures {
						fmt.Printf("     アサーション失敗: (%d, %d, %d)\n", f.X, f.Y, f.Z)
						// 期待値
						if f.Expected.Block != "" {
							fmt.Printf("       期待 block: %s\n", f.Expected.Block)
						}
						for k, v := range f.Expected.State {
							fmt.Printf("       期待 state.%s: %s\n", k, v)
						}
						// 実際の値
						if f.Actual == nil {
							fmt.Printf("       実際: ブロックが存在しません\n")
						} else {
							fmt.Printf("       実際 block: %s\n", f.Actual.Block)
							for k, v := range f.Actual.State {
								fmt.Printf("       実際 state.%s: %s\n", k, v)
							}
						}
						fmt.Printf("       理由: %s\n", f.Reason)
					}
					failed++
					totalFailed++
				}
			}

			fmt.Printf("\n  %d件中 %d件合格 (失敗: %d件, エラー: %d件)\n",
				len(tf.Tests), passed, failed, errorCount)
		}

		// 全体サマリ
		if len(files) > 1 {
			total := totalPassed + totalFailed + totalError
			fmt.Printf("\n========== 全体サマリ ==========\n")
			fmt.Printf("  %d件中 %d件合格 (失敗: %d件, エラー: %d件)\n",
				total, totalPassed, totalFailed, totalError)
		}

		// 失敗またはエラーがあった場合は終了コード1
		if totalFailed > 0 || totalError > 0 {
			os.Exit(1)
		}
	},
}

// collectTestFiles は指定されたパス（ファイルまたはディレクトリ）から
// *_test.yaml / *_test.yml ファイルを収集します。
func collectTestFiles(target string) ([]string, error) {
	info, err := os.Stat(target)
	if err != nil {
		return nil, fmt.Errorf("パスが存在しません: %s", target)
	}

	if !info.IsDir() {
		// 単一ファイルの場合はそのまま返す
		return []string{target}, nil
	}

	// ディレクトリの場合は再帰的にスキャン
	var files []string
	err = filepath.WalkDir(target, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		name := d.Name()
		if strings.HasSuffix(name, "_test.yaml") || strings.HasSuffix(name, "_test.yml") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

// parseTestFile はYAMLテストファイルをパースします。
func parseTestFile(path string) (*testrunner.TestFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ファイル読み込み失敗: %v", err)
	}

	var tf testrunner.TestFile
	if err := yaml.Unmarshal(data, &tf); err != nil {
		return nil, fmt.Errorf("YAMLパース失敗: %v", err)
	}

	return &tf, nil
}

func init() {
	rootCmd.AddCommand(testCmd)
}
