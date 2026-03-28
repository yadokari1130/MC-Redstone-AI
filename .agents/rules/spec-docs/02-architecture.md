---
trigger: model_decision
description: Redstone MCP Serverのシステム構成図、技術スタック(Python, Java/Fabric)、動作環境
---

# 2. システムアーキテクチャ・技術スタック

## 2.1. システム構成図
[ AI (LLM) ] <--(MCPプロトコル)--> [ MCPサーバー (Python) ] <--(REST API)--> [ Minecraftプラグイン/Mod (Java/Kotlin) ] <--> [ Minecraft World ]

## 2.2. 技術スタック
- **MCPサーバー側**
  - 言語: Python
  - パッケージ管理: uv
  - ライブラリ: `mcp[cli]`
- **Minecraft側**
  - 言語: Java
  - プラットフォーム: Fabric環境
  - HTTPサーバー: 組み込みWebサーバー (JavalinやNettyなどを想定)
- **動作環境**
  - Minecraftバージョン: 26.1
  - 実行環境: ローカル環境を想定
