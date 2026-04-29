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

## 6.2. `PlaceRequest`
`place-blocks` コマンドで受け付けるルート構造。
```json
{
  "blocks": [ { "x": 10, "y": 64, "z": 10, "block": "minecraft:stone" } ],
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
- **`blocks`**: 絶対座標指定。
- **`attaches`**: 土台となる `base` 座標との位置関係から向きを自動計算。
- **`connects`**: `from` と `to` の中間（距離2が必要）に向きを自動計算して配置。
- **`fills`**: `from` と `to` で指定された矩形範囲内をすべて指定されたブロックで埋める。

## 6.3. `ItemInfo`
アイテムIDと数量を表す。
```json
{
  "id": "minecraft:redstone",
  "amount": 64
}
```
