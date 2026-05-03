---
trigger: model_decision
description: API通信で利用されるブロックデータのJSONスキーマ・データモデルの定義
---

# 6. データモデル・JSON構造

## 6.1. `BlockData`
ブロックの座標、種類、およびレッドストーン状態を表す。
```json
{
  "x": 1,
  "y": 2,
  "z": 3,
  "block": "minecraft:redstone_repeater",
  "state": {
    "delay": "2",
    "facing": "north"
  }
}
```
- `verbose` オフ時のコンパクト形式: `["minecraft:id", [x, y, z], {state}]`

## 6.2. `EntityData`
エンティティのUUID、種類、座標、およびNBTデータを表す。
```json
{
  "uuid": "123e4567-e89b-12d3-a456-426614174000",
  "type": "minecraft:boat",
  "x": 10.5,
  "y": 64.0,
  "z": 12.3,
  "yaw": 90.0,
  "pitch": 0.0,
  "nbt": {
    "Type": "oak"
  }
}
```

## 6.3. `PlaceRequest`
`place-blocks` コマンドで受け付けるルート構造。

単一オブジェクト形式:
```json
{
  "blocks": [ { "x": 10, "y": 64, "z": 10, "block": "minecraft:stone" } ],
  "entities": [
    {
      "uuid": "123e4567-e89b-12d3-a456-426614174000",
      "type": "minecraft:boat",
      "x": 10.5,
      "y": 64.0,
      "z": 12.3
    }
  ],
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
  ],
  "fills": [
    {
      "from": [0, 0, 0],
      "to": [5, 5, 5],
      "block": "minecraft:stone",
      "state": {}
    }
  ]
}
```

フェーズ（配列）形式:
```json
[
  {
    "fills": [
      { "from": [10, 64, 10], "to": [12, 64, 10], "block": "minecraft:stone" }
    ]
  },
  {
    "attaches": [
      { "pos": [10, 65, 10], "component": "minecraft:lever", "base": [10, 64, 10] }
    ],
    "connects": [
      { "from": [10, 65, 10], "to": [12, 65, 10], "component": "minecraft:repeater" }
    ]
  }
]
```
- **`blocks`**: 絶対座標指定。
- **`entities`**: エンティティのスポーン情報。UUID、種類、座標、Yaw/Pitch、NBTデータを指定可能。
- **`attaches`**: 土台となる `base` 座標との位置関係から向きを自動計算。
- **`connects`**: `from` と `to` の中間（距離2が必要）に向きを自動計算して配置。
- **`fills`**: `from` と `to` で指定された矩形範囲内をすべて指定されたブロックで埋める。

**フェーズ機能**: 単一オブジェクトの他、配列 `[PlaceRequest]` としても受け付けます。配列の各要素（フェーズ）が順番にAPIに送信され、フェーズ間には自動的に1tickの待機が入ります。これにより、レッドストーン回路（ピストンの準接続やオブザーバーの更新検知など）で「配置される順番」が重要な場合に、正しい挙動を確保できます。

## 6.4. `BlocksAndEntities`
`get-blocks --include-entities` 実行時のレスポンス形式。
```json
{
  "blocks": [
    { "x": 1, "y": 2, "z": 3, "block": "minecraft:stone" }
  ],
  "entities": [
    { "uuid": "...", "type": "minecraft:boat", "x": 10.5, "y": 64.0, "z": 12.3 }
  ]
}
```

## 6.5. `ItemInfo`
アイテムIDと数量を表す。
```json
{
  "id": "minecraft:redstone",
  "amount": 64
}
```
