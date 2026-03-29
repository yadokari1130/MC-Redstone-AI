---
trigger: model_decision
description: Redstone CLIのシステム構成図、技術スタック(Python, Java/Fabric)、動作環境
---

# 2. システムアーキテクチャ・技術スタック

## 2.1. システム構成図
[ AI (LLM) ] <--(CLIコマンド実行)--> [ CLIツール (Python) ] <--(REST API)--> [ Minecraftプラグイン/Mod (Java/Kotlin) ] <--> [ Minecraft World ]

## 2.2. 技術スタック
- **CLIツール側**
  - 言語: Python
  - パッケージ管理: uv
  - ライブラリ: `argparse` または `Typer` 等
- **Minecraft側**
  - 言語: Java
  - プラットフォーム: Fabric環境
  - HTTPサーバー: 組み込みWebサーバー (Javalin)
- **動作環境**
  - Minecraftバージョン: 26.1
  - 実行環境: ローカル環境を想定
