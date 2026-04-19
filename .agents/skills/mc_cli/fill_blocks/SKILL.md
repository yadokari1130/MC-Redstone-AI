---
name: minecraft_fill_blocks
description: Minecraft の世界で指定した範囲を特定のブロックで埋めます。
---

# Minecraft 範囲埋め (fill) Skill

このスキルは、Minecraft の世界において矩形範囲を指定し、その範囲を単一のブロック種類（およびオプションの状態）で一括して埋めるためのものです。

## 使用方法

`mc-cli` ツールを使用して、以下のコマンドを実行します。

```bash
mc-cli fill <x1> <y1> <z1> <x2> <y2> <z2> <block> [--state '<JSON>']
```

### 引数

- `x1, y1, z1`: 範囲の始点座標 (整数)
- `x2, y2, z2`: 範囲の終点座標 (整数)
- `block`: 配置するブロックの ID (例: `minecraft:stone`, `minecraft:air`)
- `--state`: (任意) ブロックの状態を指定する JSON 文字列。配置する全てのブロックに適用されます。

## 使用例

### 指定範囲を石で埋める
```bash
mc-cli fill 100 60 100 110 65 110 minecraft:stone
```

### 指定範囲を空気で消去する (整地)
```bash
mc-cli fill 100 60 100 120 80 120 minecraft:air
```

### 階段を特定の向きで一括配置する
```bash
mc-cli fill 100 64 100 100 64 110 minecraft:oak_stairs --state '{"facing":"north"}'
```

## 注意事項

- 指定された 2 点を対角線とする矩形範囲の全てのブロックが置き換えられます。
- 範囲が非常に大きい場合、API のリクエスト制限やサーバーの負荷に注意してください。
- 内部的には座標ごとにブロックデータを生成して一括送信しているため、標準の `/fill` コマンドと同様の感覚で利用できます。
