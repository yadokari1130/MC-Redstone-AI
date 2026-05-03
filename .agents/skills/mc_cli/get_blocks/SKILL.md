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
- `--include-entities`: (任意) 指定すると、同じ範囲内のエンティティ情報も取得します。
- `--interval`: (任意) 実行間隔（ゲームチック、1=50ms）。デフォルトは `0`。
- `--count`: (任意) 実行回数。デフォルトは `1`。
- `--url`: (任意) サーバーの URL。デフォルトは `http://localhost:8080`。

## 出力形式

コマンドの実行結果は、解析しやすいように JSON 形式で出力されます。

### `--include-entities` 未指定時
`data` フィールドには、各実行回ごとのブロックデータリストを格納した**二次元配列**（`[][]BlockData`）がセットされます。

#### 成功時のレスポンス例 (`count=1` の場合)
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

#### 成功時のレスポンス例 (`count=2, interval=10` の場合)
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

### `--include-entities` 指定時
`data` フィールドには、`blocks` と `entities` を含むオブジェクトの配列（`[]BlocksAndEntities`）がセットされます。

#### 成功時のレスポンス例
```json
{
  "success": true,
  "data": [
    {
      "blocks": [
        ["minecraft:stone", [100, 64, 100], {}]
      ],
      "entities": [
        {
          "uuid": "550e8400-e29b-41d4-a716-446655440000",
          "type": "minecraft:boat",
          "x": 100.5,
          "y": 64.0,
          "z": 100.5,
          "yaw": 90.0,
          "pitch": 0.0,
          "nbt": { "Type": "oak" }
        }
      ]
    }
  ]
}
```

### データの構造

#### ブロックデータ
各ブロックの情報は、以下の順序の配列として表現されます。
`[BlockID, [X, Y, Z], Properties]`

- `BlockID`: ブロックの種類を表す文字列（例: `minecraft:lever`）。
- `[X, Y, Z]`: ブロックの絶対座標を表す数値配列。
- `Properties`: ブロックの状態（向き、動力強度、ON/OFF等）を表すオブジェクト（例: `{"power": "15", "facing": "north"}`）。

#### エンティティデータ (`--include-entities` 指定時)
各エンティティの情報は、以下のフィールドを持つオブジェクトとして表現されます。

- `uuid`: エンティティの一意な識別子。
- `type`: エンティティの種類（例: `minecraft:boat`, `minecraft:minecart`）。
- `x`, `y`, `z`: エンティティの座標（小数）。
- `yaw`: (任意) 水平方向の向き（度数法）。
- `pitch`: (任意) 垂直方向の向き（度数法）。
- `nbt`: (任意) エンティティのNBTデータ（例: `{"Type": "oak"}`）。

## TIPS
- レッドストーン回路の解析を行う前に、このスキルを使用して周囲の状況を把握してください。
- 範囲が広すぎると API のレスポンスが遅くなる可能性があるため、必要な範囲に限定して取得することをお勧めします。
- 座標引数 (`--pos1`, `--pos2`) は `"x,y,z"` (整数、カンマ区切り) の形式で指定します。
- 出力データの `data` フィールドは、`count` が 1 の場合でも「1回分の実行結果を格納した配列」を含む二次元配列（`data[0]` が1回目の結果）になります。
- `--include-entities` を指定した場合、プレイヤーエンティティ（`minecraft:player`）は取得対象から除外されます。
