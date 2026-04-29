---
name: minecraft_block_sticky_piston
description: 粘着ピストン（Sticky Piston）の特性、配置方法、および操作に関するガイドラインです。
---

# 粘着ピストン（Sticky Piston）の詳細仕様

## Block ID
- `minecraft:sticky_piston`

## プロパティ (Block States)
- `facing`: 粘着ピストンのヘッド（動き出す面）が向いている方向を指定します。
  - 値: `down`, `up`, `north`, `south`, `west`, `east`
  - **【重要】向きの仕様 (Minecraft準拠)**: `facing` は **「ヘッド（出力側）」** が向いている方角を指します。
  - **具体例**:
    - `facing: "up"` の場合：ヘッドは**上方向**へ伸びます。
    - `facing: "east"` の場合：ヘッドは**東方向**へ伸びます。
- `extended`: ピストンが伸びている状態かどうかを示します。
  - 値: `true` (伸びている), `false` (縮んでいる)

## 配置ガイドライン
- `place_blocks` スキルを使用する際は、`facing` プロパティでブロックを動かしたい方向を指定します。
- 通常のピストンと同様、配置された場所から `facing` 方向に向かってヘッドが伸びます。

## 操作 (Interaction)
- 直接的な右クリック操作（インタラクト）による状態の変化はありません。

## 機能・特性
- **ブロックの押し出しと引き戻し**: 動力を受けると前方のブロックを押し出し、動力を失うとそのブロックを引き戻します。詳細は [piston_slime_mechanics スキル](../../concepts/piston_slime_mechanics/SKILL.md) を参照してください。
- **1ティックパルス挙動 (1-tick Pulse)**: 非常に短いパルスを受けた際、ブロックを「置いてくる / 回収する」という特殊な挙動をします。詳細は [pulse_control_conversion スキル](../../concepts/pulse_control_conversion/SKILL.md) を参照してください。
- **準接続 (Quasi-connectivity)**: 詳細は [quasi_connectivity スキル](../../concepts/quasi_connectivity/SKILL.md) を参照してください。
- **動作時間**:
    - **伸長 (Extending)**: 2 ゲームチック (1 レッドストーンチック)。
    - **収縮 (Retracting)**: 0 ゲームチック（瞬時、ただしアニメーションには 2 ゲームチックかかる）。
- **不動ブロック**: [不動ブロックリスト](../../concepts/piston_slime_mechanics/SKILL.md#3-ブロックの可動性分類-java-edition) を参照してください。

## 使用方法
- **隠し扉**: 壁の一部を引き込んで通路を作る、一般的な隠し扉のコア部品として使用されます。
- **Tフリップフロップ**: 1チックパルス挙動を利用して、ボタン信号をレバーのようなトグル動作（ON/OFF保持）に変換する回路に使用されます。
- **フライングマシン**: 観察者、スライムブロックと組み合わせて、永続的に移動し続ける機構を作成する際に不可欠です。
