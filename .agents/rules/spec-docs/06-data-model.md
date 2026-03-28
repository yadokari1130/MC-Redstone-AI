---
trigger: model_decision
description: API通信で利用されるブロックデータのJSONスキーマ・データモデルの定義
---

# 6. データモデル・JSON構造

## 6.1. `BlockData`
ブロックの座標、種類、およびレッドストーン状態をやり取りするためのデータモデル。
```json
{
  "x": 1,
  "y": 2,
  "z": 3,
  "block": "minecraft:redstone_repeater",
  "state": {
    "delay": "2",
    "facing": "north",
    "powered": "true"
  }
}
```
- `x`, `y`, `z`: 絶対座標。
- `block`: Minecraftのブロック識別子（名前空間付き）。
- `state`: 各ブロック特有のBlockStateを表すプロパティのKey-Valueマップ。存在しない場合は空オブジェクト可。
