---
name: minecraft_block_piston
description: ピストン（Piston）の特性、配置方法、および操作に関するガイドラインです。
---

# ピストン（Piston）の詳細仕様

## Block ID
- `minecraft:piston`

## プロパティ (Block States)
- `facing`: ピストンのヘッド（動き出す面）が向いている方向を指定します。
  - 値: `down`, `up`, `north`, `south`, `west`, `east`
- `extended`: ピストンが伸びている状態かどうかを示します。
  - 値: `true` (伸びている), `false` (縮んでいる)

## 配置ガイドライン
- `place_blocks` スキルを使用する際は、`facing` プロパティでブロックを押し出したい方向を指定します。
- 配置された場所から `facing` 方向に向かってヘッドが伸びます。

## 操作 (Interaction)
- 直接的な右クリック操作（インタラクト）による状態の変化はありません。

## 機能・特性
- **ブロックの押し出し**: 動力を受けると `facing` 方向に最大12個までのブロックを押し出します。詳細は [piston_slime_mechanics スキル](../../concepts/piston_slime_mechanics/SKILL.md) を参照してください。
- **準接続 (Quasi-connectivity)**: Java版特有の仕様です。詳細は [quasi_connectivity スキル](../../concepts/quasi_connectivity/SKILL.md) を参照してください。
- **動作時間**:
    - **伸長 (Extending)**: 2 ゲームチック (1 レッドストーンチック)。
    - **収縮 (Retracting)**: 0 ゲームチック（瞬時、ただしアニメーションには 2 ゲームチックかかる）。
- **不動ブロック**: 岩盤、黒曜石、[不動ブロックリスト](../../concepts/piston_slime_mechanics/SKILL.md#3-ブロックの可動性分類-java-edition) にあるブロックは押し出すことができません。

## 使用方法
- **ブロックの移動**: 隠し扉やエスカレーター、自動収穫機などでブロックを物理的に動かすために使用されます。
- **信号の遮断**: 導体ブロックを移動させることで、レッドストーン信号の伝達を物理的にON/OFFするスイッチとして機能させることができます。
