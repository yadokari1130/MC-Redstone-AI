---
name: minecraft_block_button
description: ボタン（オーク、石）の特性、配置方法、および操作に関するガイドラインです。
---

# ボタンの詳細仕様

## Block ID
- `minecraft:stone_button`, `minecraft:polished_blackstone_button` (石系)
- `minecraft:oak_button`, `minecraft:spruce_button`, `minecraft:birch_button`, `minecraft:jungle_button`, `minecraft:acacia_button`, `minecraft:dark_oak_button`, `minecraft:mangrove_button`, `minecraft:cherry_button`, `minecraft:bamboo_button`, `minecraft:crimson_button`, `minecraft:warped_button` (木系)

## プロパティ (Block States)
- `face`: ボタンが取り付けられている面。
    - `floor`: 床面（上面）
    - `wall`: 壁面（側面）
    - `ceiling`: 天井面（下面）
- `facing`: ボタンが向いている方位。
    - `north`, `south`, `east`, `west`
- `powered`: ボタンのON/OFF状態（押されているかどうか）。
    - `true`: 押された状態（信号出力中）
    - `false`: 通常状態（信号停止中）

## 配置ガイドライン
- **設置場所**: 不透過のフルブロックの上面、側面、下面に設置可能です。
- 配置は `.agents/skills/mc_cli/place_blocks/SKILL.md` を使用して行います。基本的に `attaches` を使用して配置してください。これにより、土台となるブロックに合わせて `facing` や `face` が自動的に計算されます。

## 操作 (Interaction)
- **mc_cli_interact_block** を使用してボタンを押すことができます。
- **パルス時間**:
    - **石系ボタン**: 10 レッドストーンチック (20 ゲームチック / 1.0秒)
    - **木系ボタン**: 15 レッドストーンチック (30 ゲームチック / 1.5秒)
- **木製ボタンの特性**: 矢、三叉槍、釣竿の浮きなどがヒットした際にも起動します。

## 機能・特性
- **信号強度**: 強度 **15** を出力します。
- **動力化**: 
    - 隣接するレッドストーンコンポーネントに信号を伝達します。
    - 取り付けられているブロックを **強動力化 (Strong Power)** します。

## 使用方法
- **パルス信号の生成**: 回路のトリガーや、手動スイッチとして使用されます。
