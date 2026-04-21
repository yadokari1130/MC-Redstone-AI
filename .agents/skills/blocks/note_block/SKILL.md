---
name: minecraft_block_note_block
description: 音符ブロックの特性、配置方法、および操作に関するガイドラインです。
---

# 音符ブロックの詳細仕様

## Block ID
- `minecraft:note_block`

## プロパティ (Block States)
- `instrument`: 楽器の種類。以下の値を取ります：
    - `harp` (デフォルト), `basedrum`, `snare`, `hat`, `bass`, `flute`, `bell`, `chime`, `guitar`, `iron_xylophone`, `cow_bell`, `didgeridoo`, `bit`, `banjo`, `pling`, `xylophone`, `zombie`, `skeleton`, `creeper`, `dragon`, `wither_skeleton`, `piglin`, `custom_head`
- `note`: 音程（0–24の整数。24で2オクターブに相当）。
- `powered`: レッドストーン信号を受けているか（`true`/`false`）。

## 配置ガイドライン
- 音符ブロックを鳴らすためには、ブロックの真上が空気ブロック（またはMobの頭などの一部の例外）である必要があります。
- 音符ブロックの下にあるブロックの種類によって、`instrument`（楽器）が決定されます。
- 配置時に `instrument` や `note` を直接指定することも可能です。

## 操作 (Interaction)
- **mc_cli_interact_block**: 右クリック操作を再現し、`note` プロパティの値を 1 つ上げます。
- **オブザーバー検知**: 音符ブロックの `note` 値が変化したり、音が鳴ったりする際の `powered` 状態の変化を、隣接するオブザーバーで検知できます。詳細は [block_update_detection スキル](../../concepts/block_update_detection/SKILL.md) を参照してください。

## 機能・特性
- レッドストーン信号を受け取った瞬間（`powered` が `true` に変化した時）に、設定された楽器と音程で音を鳴らします。
- 真上に不透過ブロックがある場合、信号を送っても音は鳴りません。

## 使用方法
- メロディの演奏や、回路の作動状況を音で知らせる通知システムに使用されます。
- 下に置くブロックを使い分けることで、ドラム、ギター、ピアノなどの様々な音色を奏でることができます。
