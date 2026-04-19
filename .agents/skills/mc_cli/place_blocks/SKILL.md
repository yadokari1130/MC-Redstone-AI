---
name: minecraft_place_blocks
description: Minecraft の世界に指定したブロックを配置します。
---

# Minecraft ブロック配置 Skill

このスキルは、Minecraft サーバー（Fabric）の HTTP API を使用して、指定した座標にブロックを一括で配置するためのものです。

## 使用方法

`mc-cli` ツールを使用して、以下のコマンドを実行します。

```bash
mc-cli place-blocks --blocks '<JSON文字列>'
```

または、ファイルを指定して配置します。

```bash
mc-cli place-blocks --blocks '@path/to/blocks.json'
```

### 引数
- `--blocks`: 配置するブロック情報を含む JSON 文字列。または、`@` を付けたファイルパス。
- `--url`: (任意) サーバーの URL。デフォルトは `http://localhost:8080`。

## 入力形式 (JSON)

配置するブロックのデータ構造（`blocks`, `attaches`, `connects`）の詳細については、[block_design/SKILL.md](file:///home/yadokari/redstone_ai/.agents/skills/mc_cli/block_design/SKILL.md) を参照してください。


## TIPS
- アタッチするブロックや繋ぐブロックは自動で向きが計算され、API 側に送られるため JSON 内で明示的に指定する必要はありません。
- 既存のブロックを上書きして配置します。
- 配置完了後、設置したすべての座標に対して自動的にブロックアップデートが実行されます。これにより、更新が必要なブロックが即座に同期されます。
- 配置が完了すると、成功メッセージが JSON 形式で出力されます。
