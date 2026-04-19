---
name: mc_cli_get_blocks
description: Minecraft の世界から指定した範囲のブロック情報を取得します。
---

# Minecraft ブロック取得 Skill

このスキルは、Minecraft サーバー（Fabric）の HTTP API を使用して、指定された座標範囲（直方体）に含まれるブロックの情報を JSON 形式で取得し、解析するためのものです。

## 使用方法

`mc-cli` ツールを使用して、以下のコマンドを実行します。

```bash
mc-cli get-blocks --pos1 <開始座標> --pos2 <終了座標>
```

### 引数
- `--pos1`: 範囲の開始座標 `"x,y,z"`（整数）。
- `--pos2`: 範囲の終了座標 `"x,y,z"`（整数）。
- `--interval`: (任意) 実行間隔（ゲームチック、1=50ms）。デフォルトは `0`。
- `--count`: (任意) 実行回数。デフォルトは `1`。
- `--url`: (任意) サーバーの URL。デフォルトは `http://localhost:8080`。

## 出力形式

コマンドの実行結果は、解析しやすいように JSON 形式で出力されます。
`data` フィールドには、各実行回ごとのブロックデータリストを格納した**二次元配列**（`[][]BlockData`）がセットされます。

### 成功時のレスポンス例 (`count=1` の場合)
```json
{
  "success": true,
  "data": [
    [
      ["minecraft:stone", [100, 64, 100], {}],
      ["minecraft:redstone_wire", [101, 64, 100], {"power": "15", "north": "side", "south": "side"}]
    ]
  ]
}
```

### 成功時のレスポンス例 (`count=2, interval=10` の場合)
```json
{
  "success": true,
  "data": [
    [
      ["minecraft:redstone_wire", [100, 64, 100], {"power": "0"}]
    ],
    [
      ["minecraft:redstone_wire", [100, 64, 100], {"power": "15"}]
    ]
  ]
}
```

### データの構造
各ブロックの情報は、以下の順序の配列として表現されます。
`[BlockID, [X, Y, Z], Properties]`

- `BlockID`: ブロックの種類を表す文字列。
- `[X, Y, Z]`: ブロックの絶対座標を表す数値配列。
- `Properties`: ブロックの状態を表すオブジェクト。

### プロパティの解説
- `x, y, z`: ブロックの絶対座標。
- `block`: ブロックの ID（例: `minecraft:lever`）。
- `properties`: ブロックの状態（向き、電力、オン/オフの状態など）。

## TIPS
- レッドストーン回路の解析を行う前に、このスキルを使用して周囲の状況を把握してください。
- 範囲が広すぎると API のレスポンスが遅くなる可能性があるため、必要な範囲に限定して取得することをお勧めします。
