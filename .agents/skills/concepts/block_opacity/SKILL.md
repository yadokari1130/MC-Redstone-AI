---
name: redstone_concept_block_opacity
description: ブロックの透過性（不透過・透過）の違いによるレッドストーン信号の伝達、遮断、および垂直配線の知識ガイド。
---

# ブロックの透過性 (Block Opacity / Transparency)

## 概要
レッドストーン回路において、ブロックの「透過性」は単なる見た目の問題ではなく、信号を「通すか」「遮断するか」「蓄えるか」を決定する重要な物理特性です。ブロックは大きく分けて **不透過ブロック（不透過）** と **透過ブロック（透過）** の2種類に分類され、それぞれ異なるルールに従います。

## 仕様・論理

### 1. 不透過ブロック (Opaque / Solid Blocks)
フルブロック（石、土、木材、羊毛など）の多くがこれに該当します。
- **動力の保持**: レッドストーン信号を受けて「強動力」または「弱動力」状態になることができます。
- **信号の遮断 (ピンチ)**: レッドストーンダストが1段上がる際の通り道に不透過ブロックを置くと、ワイヤーが切断されます。これを利用して意図的に信号を遮断（断線）させることが可能です。
- **隣接への伝達**: 動力を持った不透過ブロックは、隣接するレッドストーンランプやピストン等を起動できます。

### 2. 透過ブロック (Transparent Blocks)
ガラス、ハーフブロック、階段、リーフ、チェスト、**スライムブロック、ハチミツブロック**などが該当します。
- **動力を保持しない**: 透過ブロック自体は動力を持つことができません。
- **信号を遮断しない**: レッドストーンダストが斜めに接続されている場所に透過ブロックを置いても、接続は維持されます。
- **Java版特有の垂直伝達**: グロウストーンや上ハーフブロックなどは、**下から上へは信号を伝えるが、上から下へは信号を伝えません**。これは「弱動力化」されない性質により、上の段のダストが下の段のダストを（ブロック経由で）起動できないためです。

## 推奨レイアウト

### 階段状の配線比較
#### 不透過ブロックの場合（信号を遮断する「ピンチ」）
下段ダストの真上に不透過ブロックを置くと、上段への接続が物理的に切断されます。
```text
(横から見た断面図)
[P][D2]   P : 遮断ブロック (不透過 / ここで線が切れる)
[D1][S2]  D1: 下段ダスト   / D2: 上段ダスト
[S1][  ]  S1: 下段の土台   / S2: 上段の土台
```

#### 透過ブロックの場合（信号を遮断しない）
透過ブロック（ガラス等）は、ダストの真上にあっても接続を妨げません。
```text
(横から見た断面図)
[G][D2]   G : 透過ブロック (信号は遮断されず接続される)
[D1][S2]  D1: 下段ダスト   / D2: 上段ダスト
[S1][  ]  S1: 下段の土台   / S2: 上段の土台
```

### ブロック分類の例
- **不透過**: `minecraft:stone`, `minecraft:oak_planks`, `minecraft:wool`, `minecraft:target`, `minecraft:crafter`, `minecraft:redstone_block`, `minecraft:dispenser`, `minecraft:dropper`
- **透過**: `minecraft:glass`, `minecraft:glowstone`, `minecraft:oak_slab`, `minecraft:oak_stairs`, `minecraft:honey_block`, `minecraft:slime_block`, `minecraft:copper_bulb`, `minecraft:hopper`, `minecraft:chest`, `minecraft:observer`, `minecraft:repeater`, `minecraft:comparator`

## 26.1 における特記事項
- **銅の電球 (Copper Bulb)**: 外見はフルブロックですが、**透過ブロック** です。
- **クラフター (Crafter)**: **不透過ブロック** です。
- **レッドストーンブロック**: 常に動力を出力し続けますが、レッドストーン的には **不透過ブロック** です。
- **最新の透過設定**: リソースパック等で見た目を透過にしても、ブロック内部の「レッドストーン的な透過性」は不変です。回路設計時はブロックID自体の特性に依存してください。

## 使用ガイドライン
- **断線の活用**: 2つの並行するレッドストーンラインが隣接している場合、間に不透過ブロックを置くことで混線を防ぎつつ、高さを変えて信号を制御できます。
- **コンパクト化**: ガラスやハーフブロックを使うことで、信号を遮断せずに信号線を「重ねる」ことができ、回路をより多層的でコンパクトに設計できます。
- **上方向への高速伝達**: レッドストーントーチのタワー（遅延あり）の代わりに、グロウストーンや上ハーフブロックの階段（遅延なし）を使うことで、信号を垂直に即座に伝えることができます。

## ソース (References)
- [Minecraft Wiki - Opacity](https://minecraft.wiki/w/Opacity)
- [Minecraft Wiki - Redstone circuits#Vertical_transmission](https://minecraft.wiki/w/Redstone_circuits#Vertical_transmission)
- [Minecraft Wiki - Redstone components#Copper_Bulb](https://minecraft.wiki/w/Redstone_components#Copper_Bulb)
- [Minecraft Wiki - Redstone components#Crafter](https://minecraft.wiki/w/Redstone_components#Crafter)
