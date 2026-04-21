---
name: redstone_concept_block_update_order
description: ブロック更新順序（Block Update Order）、座標依存（Locationality）、方角依存（Directionality）の論理、構成方法、および使用に関するガイドラインです。
---

# ブロック更新順序 (Block Update Order)

## 概要
ブロック更新順序とは、Minecraft内で複数のイベントが同じゲームティック（Tick）内に発生した際、それらがどの順番で処理されるかを決定する仕組みです。特にレッドストーン回路において、同時に動くはずのピストンやレッドストーンダストが予期せぬ動作をする原因（座標依存や方角依存）を理解し、制御するために不可欠な知識です。

## 仕様・論理
Minecraft Java Edition 1.21.2 以降では、レッドストーンワイヤー（ダスト）の更新ロジックが大幅に改善され、より予測可能になっています。

### 1. 標準的な近傍更新（Neighbor Update / NC）
ブロックの状態が変化したとき、周囲の6ブロックに送られる更新（Neighbor Changed）の標準的な順序は以下の通りです。
1.  **西 (West / -X)**
2.  **東 (East / +X)**
3.  **下 (Down / -Y)**
4.  **上 (Up / +Y)**
5.  **北 (North / -Z)**
6.  **南 (South / +Z)**
これを「**W-E-D-U-N-S**」順と呼びます。レバーやトーチなどの動力源はこの順序で周囲を更新するため、回路の向きによって動作が変わる「方角依存（Directionality）」の原因となります。

### 2. レッドストーンワイヤーの更新（1.21.2 以降の実験的機能）
「レッドストーンの実験」を有効にした場合、以下の優先順位で更新が行われます。
-   **距離優先 (Distance Priority)**: 動力源に近いワイヤーから順に更新されます。
-   **タイブレーク (Left-First)**: 信号の進行方向に対して「左側」にあるワイヤーが優先されます。
-   **予測可能性**: 方角依存 (W-E-D-U-N-S) や座標依存の多くが解消されます。

### 3. 座標依存 (Locationality)
同じティック内に発生する複数のイベントが、座標に基づいて処理順序が決まる現象です。
-   **対策**: リピーターによる遅延を 1 ティック以上挟むことで、同一ティック内の競合を避けるのが確実です。

---

## 使用ガイドライン
- **方角依存の回避**: 非ワイヤー系動力源（トーチ、レバー、オブザーバー等）は W-E-D-U-N-S 順の影響を受けます。
- **決定論的な設計**: 座標や方角に依存しない回路を作るには、**「リピーターによる明示的な遅延」** を活用してください。
- **26.1 での推奨**: 複雑な回路設計では、ワールド設定で「レッドストーンの実験」を有効にすることが推奨されます。

## ソース (References)
- [Minecraft Wiki - Block update](https://minecraft.wiki/w/Block_update)
- [Minecraft Wiki - Redstone wire#Redstone_Experiments](https://minecraft.wiki/w/Redstone_wire#Redstone_Experiments)
- [Mojang - Redstone Experiments (Snapshot 24w33a+)](https://www.minecraft.net/en-us/article/redstone-experiments-minecraft)
