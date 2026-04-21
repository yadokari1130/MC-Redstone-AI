---
name: minecraft_block_detector_rail
description: ディテクターレールの特性、配置方法、および操作に関するガイドラインです。
---

# ディテクターレールの詳細仕様

## Block ID
- `minecraft:detector_rail`

## プロパティ (Block States)
- `powered`: `true`, `false`
  - トロッコがレールの上にあるかどうかを示します。
- `shape`: `north_south`, `east_west`, `ascending_east`, `ascending_west`, `ascending_north`, `ascending_south`
  - レールの向きと傾斜を指定します。
- `waterlogged`: `true`, `false`
  - ブロック内に水が含まれているかどうかを指定します。

## 配置ガイドライン
- `shape` プロパティを使用して向きを指定します。パワードレールと同様、カーブを作ることはできません。
- 設置は `place_blocks` を行いますが、`powered` 状態はトロッコの有無によって自動的に決まるため、通常は `false` で設置します。

## 操作 (Interaction)
- プレイヤーによる直接の右クリック操作で状態を変えることはできません。
- トロッコがこのレールの上に乗ることで `powered` が `true` になり、周囲にレッドストーン信号を出力します。

## 機能・特性
- 感圧版のレール版として機能します。
- トロッコを検知すると、隣接するブロックに強度 15 の信号を送ります。
- また、**レッドストーンコンパレーター** を接続すると、その上のトロッコ（コンテナ付き）のアイテム量に応じた強度の信号を出力します。詳細は [inventory_detection スキル](../../concepts/inventory_detection/SKILL.md) および [signal_strength_management スキル](../../concepts/signal_strength_management/SKILL.md) を参照してください。

## 使用方法
- トロッコの通過を検知して信号を送る、自動仕分けシステムや自動駅のトリガーとして使用されます。
- また、コンパレータと組み合わせて、アイテムの積み込みが完了したことを検知してトロッコを出発させるなどの高度な制御に利用されます。
