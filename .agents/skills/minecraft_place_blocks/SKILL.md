---
name: minecraft_place_blocks
description: Minecraft の世界に指定したブロックを配置します。
---

# Minecraft ブロック配置 Skill

このスキルは、Minecraft サーバー（Fabric）の HTTP API を使用して、指定した座標にブロックを一括で配置するためのものです。

## 使用方法

`mc-cli` ツールを使用して、以下のコマンドを実行します。

```bash
./mc-cli/mc-cli place-blocks --blocks '<JSON配列>'
```

または、ファイルを指定して配置します。

```bash
./mc-cli/mc-cli place-blocks --blocks '@path/to/blocks.json'
```

### 引数
- `--blocks`: 配置するブロック情報の JSON 配列。または、`@` を付けたファイルパス。
- `--url`: (任意) サーバーの URL。デフォルトは `http://localhost:8080`。

## 入力形式 (JSON)

配置するブロックのリストを JSON 形式で提供します。

### 例: 2つのブロックを配置する
```json
[
  {
    "x": 100,
    "y": 64,
    "z": 100,
    "block": "minecraft:lever",
    "properties": {
      "facing": "north",
      "powered": "false"
    }
  },
  {
    "x": 101,
    "y": 64,
    "z": 100,
    "block": "minecraft:redstone_wire",
    "properties": {
      "power": "0"
    }
  }
]
```

## TIPS
- 回路を構築する場合、まずは設計図（JSON）をファイルとして保存し、`@filename` 形式で流し込むのが確実です。
- 既存のブロックを上書きして配置します。
- 配置が完了すると、成功メッセージが JSON 形式で出力されます。
