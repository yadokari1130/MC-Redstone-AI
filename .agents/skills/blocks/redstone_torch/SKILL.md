---
name: minecraft_block_redstone_torch
description: レッドストーントーチの特性、配置方法、および機能に関するガイドラインです。
---

# レッドストーントーチの詳細仕様

## Block ID
- 床置き: `minecraft:redstone_torch`
- 壁掛け: `minecraft:redstone_wall_torch`

## プロパティ (Block States)
- `lit` (boolean): トーチが点灯しているかどうか。
- `facing` (direction): 壁掛け（`redstone_wall_torch`）の際、どの方向を向いているか（`north`, `south`, `east`, `west`）。

## 配置ガイドライン
- **設置場所**: 不透明で導電性のあるブロック（石など）の上面または側面に設置できます。
- **制限**: ブロックの下側（底面）には設置できません。
- **向き**: 指定した面とIDに応じて、床置きか壁掛けかが決定されます。壁掛けの場合は、その背面が支持ブロックとなります。
- **設置方法**: `place_blocks` で配置する際は基本的に `attaches` を使用してください。これにより、土台となるブロックに合わせて `facing` が適切に設定され、`redstone_torch` と `redstone_wall_torch` の ID も自動的に選択されます。

## 操作 (Interaction)
- レッドストーントーチはプレイヤーによる直接的なインタラクト操作（右クリックなどでの状態変更）は行えません。

## 機能・特性
- **反転論理 (NOT Gate)**: アタッチされているブロックが動力化されると、トーチは消灯します。詳細は [logic_gates スキル](../../concepts/logic_gates/SKILL.md) を参照してください。
- **信号出力**:
    - 真上のブロックを **強動力 (Strong Power)** します。
    - 隣接する 5 方向（アタッチ先を除く）に信号を伝達します。詳細は [strong_weak_power スキル](../../concepts/strong_weak_power/SKILL.md) を参照してください。
- **遅延**: 状態変化時に **1 レッドストーンチック (2 ゲームチック)** の遅延が発生します。
- **焼き切れ (Burn-out)**: 短時間（3秒以内）に 8 回以上 ON/OFF を繰り返すと焼き切れ、一定時間（約5秒）機能しなくなります。

## 使用方法
- **論理反転**: 信号のON/OFFを逆転させるNOTゲートとして使用されます。
- **垂直伝達 (Torch Tower)**: トーチを垂直に積み重ねることで、信号を上方へ効率よく伝達できます。
- **信号源**: 常にONの信号を供給する安定した信号源として使用されます。
