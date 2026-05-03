---
name: mc_cli_block_design
description: Minecraft のレッドストーン回路や構造物を設計するための、JSON データモデルの仕様と設計ガイドラインです。
---

# Minecraft ブロック設計 Skill

このスキルは、Minecraft サーバーに配置するためのブロック設計データ（JSON）を作成するためのガイドラインと仕様を提供します。
`mc-cli place-blocks` コマンドや `mc-cli test` コマンドの入力データとして使用する JSON フォーマットの詳細を定義しています。

## 設計データ構造 (JSON)

配置データは、`blocks`, `entities`, `attaches`, `connects`, `fills` の5つのリストを持つオブジェクトです。1つのオブジェクト内では `blocks` -> `entities` -> `attaches` -> `connects` の順序で実行され、`fills` は事前に `blocks` リストに展開された上で配置されます。

また、**配列形式（フェーズ）**を使用することで、配置を複数の段階に分けて順番に実行できます。これはレッドストーン回路（ピストンの準接続やオブザーバーの更新検知など）で「配置される順番」が重要な場合に必要です。配列の各要素（フェーズ）が順番にAPIに送信され、フェーズ間には自動的に1tickの待機が入ります。

### 1. `blocks` (基本配置)
指定した座標に直接ブロックを配置します。

- **`x`, `y`, `z`**: 配置座標（整数）。
- **`block`**: ブロック ID（例: `minecraft:iron_block`）。
- **`state`**: (任意) ブロックの状態（例: `{"powered": "true"}`, `{"lit": "true"}`）。

### 2. `attaches` (相対配置/アタッチ)
土台となるブロックに対して部品を取り付けます。土台との位置関係から、適切な向き（`facing` 等）が自動計算されます。

- **`pos`**: 配置する部品の座標 `[x, y, z]`。
- **`component`**: 配置するブロック ID。
  - **対応例**: `minecraft:redstone_torch`, `minecraft:lever`, `minecraft:stone_button`, `minecraft:oak_pressure_plate`, `minecraft:tripwire_hook`, `minecraft:ladder` 等。
- **`base`**: 土台となるブロックの座標 `[x, y, z]`。

**自動計算の例:**
- `base` が `[100, 64, 100]`、`pos` が `[100, 65, 100]`（真上）の場合：上向き、または床置きとして配置されます。
- `base` が `[100, 64, 100]`、`pos` が `[101, 64, 100]`（横）の場合：壁掛けとして配置されます。

### 3. `connects` (接続配置)
2つのブロックの間に部品を配置します。`from` から `to` の方向を向くように向きが自動計算されます。

- **`from`**: 開始地点の座標 `[x, y, z]`。
- **`to`**: 終了地点の座標 `[x, y, z]`。
- **`component`**: 配置するブロック ID。
  - **対応例**: `minecraft:repeater`, `minecraft:comparator`, `minecraft:observer`, `minecraft:powered_rail` 等。

**制約:**
- `from` と `to` は同じ軸上（X, Y, または Z）にある必要があります。
- `from` と `to` の間にはちょうど1マスの空き（距離が2）が必要です。部品はその中間のマスに配置されます（既存のブロックは上書きされます）。

### 4. `fills` (範囲埋め)
指定された2点（`from` と `to`）で囲まれる直方体範囲内を、すべて指定されたブロックで埋めます。

- **`from`**: 範囲の始点座標 `[x, y, z]`。
- **`to`**: 範囲の終点座標 `[x, y, z]`。
- **`block`**: 埋めるブロック ID。
- **`state`**: (任意) ブロックの状態。

**注意:** `fills` は送信前に内部で `blocks` リストに自動展開されるため、最終的には通常の `blocks` 配置と同じAPI呼び出しで処理されます。

### 5. `entities` (エンティティ配置)
指定した座標にエンティティをスポーンさせます。ブロックと同時に配置できます。

- **`uuid`**: (任意) エンティティの一意な識別子。未指定時は自動生成されます。
- **`type`**: エンティティの種類（例: `minecraft:boat`, `minecraft:minecart`, `minecraft:item_frame`）。
- **`x`, `y`, `z`**: スポーン座標（小数）。
- **`yaw`**: (任意) 水平方向の向き（度数法）。
- **`pitch`**: (任意) 垂直方向の向き（度数法）。
- **`nbt`**: (任意) エンティティに適用するNBTデータ。JSON形式で指定します（例: `{"Type": "oak"}`）。

**注意:** `entities` は `blocks` の後に処理されます。エンティティをスポーンさせる前に、土台となるブロックが既に配置されていることを確認してください。

### 6. フェーズ（配列形式）
単一のオブジェクトではなく、オブジェクトの配列 `[ {...}, {...} ]` として設計データを記述できます。配列の各要素が「フェーズ」として順番に処理されます。

**用途例:**
- フェーズ1: 土台ブロックを配置
- フェーズ2: レッドストーンダストとトーチを配置（土台の更新を確実にさせる）
- フェーズ3: ピストンやオブザーバーを配置（準接続や更新検知を正しく動作させる）

## 設計ガイドライン

1.  **土台の優先**: 複雑な回路を設計する場合、まず `blocks` で不透過ブロック（石、鉄ブロックなど）を配置し、その後に `attaches` や `connects` でレッドストーン部品を配置するようにリストを構成してください。
2.  **座標の整合性**: 回路全体を特定の座標（例: `0, 0, 0`）からの相対座標で設計し、配置時にオフセットを加算するようにすると再利用性が高まります。
3.  **向きの自動化**: 手動で `state` の `facing` を計算するのではなく、可能な限り `attaches` や `connects` を使用して、物理的な位置関係から向きを決定させてください。これにより、回路を回転・反転させた際のミスを防げます。
4.  **ブロック状態の指定**: `state` を使用して、詳細な設定を明示的に指定できます。
    - `minecraft:repeater`: `{"delay": "1"}` (1〜4)
    - `minecraft:comparator`: `{"mode": "compare"}` (`compare` または `subtract`)
    - `minecraft:lever`: `{"powered": "true"}`
    - `minecraft:redstone_lamp`: `{"lit": "true"}`
    - `minecraft:piston`: `{"extended": "false"}`

## 設計例

### 単一オブジェクト形式

```json
{
  "blocks": [
    { "x": 10, "y": 64, "z": 10, "block": "minecraft:stone" },
    { "x": 11, "y": 64, "z": 10, "block": "minecraft:stone" }
  ],
  "entities": [
    {
      "type": "minecraft:boat",
      "x": 10.5,
      "y": 64.5,
      "z": 10.5,
      "yaw": 90.0,
      "nbt": { "Type": "oak" }
    }
  ],
  "attaches": [
    {
      "pos": [10, 65, 10],
      "component": "minecraft:redstone_torch",
      "base": [10, 64, 10]
    }
  ],
  "connects": [
    {
      "from": [10, 65, 10],
      "to": [12, 65, 10],
      "component": "minecraft:repeater"
    }
  ]
}
```

### フェーズ（配列）形式

以下は、土台を先に配置し、その後レッドストーン部品を配置する2フェーズの例です。

```json
[
  {
    "fills": [
      {
        "from": [10, 64, 10],
        "to": [12, 64, 10],
        "block": "minecraft:stone"
      }
    ]
  },
  {
    "attaches": [
      {
        "pos": [10, 65, 10],
        "component": "minecraft:lever",
        "base": [10, 64, 10]
      }
    ],
    "connects": [
      {
        "from": [10, 65, 10],
        "to": [12, 65, 10],
        "component": "minecraft:repeater"
      }
    ]
  }
]
```
