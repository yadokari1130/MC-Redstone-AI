---
name: minecraft_block_door
description: ドア（オークのドア等）の特性、配置方法、および操作に関するガイドラインです。
---

# ドアの詳細仕様

## Block ID
- `minecraft:oak_door` (ほか、各種木材、鉄、銅のドアが存在)
    - 木製: `spruce_door`, `birch_door`, `jungle_door`, `acacia_door`, `dark_oak_door`, `mangrove_door`, `cherry_door`, `bamboo_door`, `crimson_door`, `warped_door`
    - 金属製: `iron_door`, `copper_door` (および酸化バリエーション)

## プロパティ (Block States)
- `facing`: `north`, `south`, `east`, `west` (ドアの正面が向いている方向)
- `half`: `lower`, `upper` (上下どちらのパーツか)
- `hinge`: `left`, `right` (蝶番の位置)
- `open`: `true`, `false` (開閉状態)
- `powered`: `true`, `false` (レッドストーン信号を受けているか)

## 配置ガイドライン
- `place_blocks` を使用して設置します。ドアは上下2ブロックを占有するため、`lower` と `upper` の両方を配置する必要があります。
- 設置には下に不透明なフルブロックが必要です。
- `facing` と `hinge` を適切に設定することで、観音開きのペアを作成したり、開閉方向を調整したりできます。

## 操作 (Interaction)
- **木製ドア**: 右クリック（インタラクト）することで `open` 状態を反転させることができます。
- **鉄製・銅製ドア**: 右クリックによる操作はできません。レッドストーン信号による制御のみ可能です。

## 機能・特性
- **レッドストーン制御**: 信号を受けると `powered` が `true` になり、`open` 状態が切り替わります。
- **上下連動**: 上下のブロック（`half=lower/upper`）は内部的に連動しており、一方の状態が変わると他方も同期します。

## 使用方法
- **回路の出力先**: レッドストーン信号を受けて動作するアクチュエーターとして、隠し扉やトラップ、自動ドアなどの構築に使用されます。
- **鉄のドア**: プレイヤーの直接操作を制限し、ボタンやレバー、感圧版などの論理回路を介したアクセス制御を行う際に多用されます。
