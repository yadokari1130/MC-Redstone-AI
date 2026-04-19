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

配置するブロックのリストを JSON オブジェクト形式で提供します。
`blocks`, `attaches`, `connects` の3つのリストを含めることができ、この順番で配置処理が行われます。

- **`blocks`**: 通常のブロック配置。指定した座標にブロックを置きます。
- **`attaches`**: 土台となるブロック（`base`）に対して部品（`component`）を取り付けます。`base` からの相対位置によって自動的に向き（`facing` 等）が計算されます。（同じ JSON 内の `blocks` を土台として利用することも可能です）
- **`connects`**: 指定した2点（`from` と `to`）の間に部品（`component`）を配置し、`from` から `to` の方向を向くように自動計算します。（※ `from` と `to` は同じ軸上で距離がちょうど2マス、つまり間に1マスだけ空きがある状態である必要があります）

### 例: 通常ブロック、アタッチ、接続の複合配置
```json
{
  "blocks": [
    {
      "x": 100,
      "y": 64,
      "z": 100,
      "block": "minecraft:iron_block"
    },
    {
      "x": 101,
      "y": 64,
      "z": 100,
      "block": "minecraft:redstone_wire",
      "state": {
        "power": "0"
      }
    }
  ],
  "attaches": [
    {
      "component_x": 100,
      "component_y": 65,
      "component_z": 100,
      "component": "minecraft:redstone_torch",
      "base_x": 100,
      "base_y": 64,
      "base_z": 100
    }
  ],
  "connects": [
    {
      "from_x": 101,
      "from_y": 64,
      "from_z": 100,
      "to_x": 103,
      "to_y": 64,
      "to_z": 100,
      "component": "minecraft:repeater"
    }
  ]
}
```

## TIPS
- 回路を構築する場合、まずは設計図（JSON）をファイルとして保存し、`@filename` 形式で流し込むのが確実です。
- アタッチするブロックや繋ぐブロックは自動で向きが計算され、API 側に送られるため JSON 内で明示的に指定する必要はありません。
- 既存のブロックを上書きして配置します。
- 配置完了後、設置したすべての座標に対して自動的にブロックアップデートが実行されます。これにより、レッドストーン回路などの更新が必要なブロックが即座に同期されます。
- 配置が完了すると、成功メッセージが JSON 形式で出力されます。
