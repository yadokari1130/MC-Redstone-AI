# Redstone AI

このプロジェクトは、AIを活用してMinecraft内で完全に自動化されたレッドストーン回路を設計・生成・動作検証するための統合リポジトリです。

## 概要

レッドストーン回路の構築には複雑な論理設計と立体的なブロック配置が求められます。本プロジェクトは、AIにそれらの回路設計や配置を理解させ、Minecraft上で直接実行させるためのサーバー群、連携ツール、および各種ワークフローを提供するプロジェクトです。
最終的に、回路の要件定義から具体的なシステムの設計、ワールド上への配置、さらには動作確認に至るまでの一連のプロセスをAIが自律的かつ正確に行うことを目指しています。

## リポジトリ構成

本リポジトリはモノレポ構成を採用しており、AIによる回路自動生成を実現するための複数のコンポーネントが統合されています。

### 1. AIリソース群 (ワークフロー・スキル)
AIエージェントに自律的な思考とレッドストーン設計の専門知識を与えるためのリソースです。
- **workflows**: 回路の要件定義から設計・配置までのプロセスを、AIが順を追って実行するための手順書です。
- **skills**: Minecraftの座標の扱い方やレッドストーン特有の複雑なブロック状態(BlockState)、特殊な設置手順などタスクを処理するための知識（プロンプトやルール）をまとめたものです。

### 2. Minecraftプラグイン (Java HTTPサーバー)
Minecraftサーバー内部で軽量な機能特化型HTTPサーバーを立ち上げ、REST APIの形式でワールドのブロック情報を直接取得・配置・操作するためのシステムです。設置順序による回路崩壊などの問題を未然に防ぎ、安全かつ確実な状態取得・一括配置を実現します。

### 3. CLIツール (Python)
AIエージェントとMinecraftプラグイン（HTTPサーバー）の間に立ち、中継および制御を行うCLIツールです。AIはここから提供される操作用コマンド (`get_blocks`, `place_blocks`, `interact_block`) を駆使して、安全かつ正確に環境を操作します。

## 詳細仕様書

システムの具体的な要件やアーキテクチャについては、以下の仕様書ドキュメント群にまとめられています。

- [総合仕様書・目次](.agents/rules/specification.md)
  - 1. [はじめに](.agents/rules/spec-docs/01-introduction.md)
  - 2. [システムアーキテクチャ・技術スタック](.agents/rules/spec-docs/02-architecture.md)
  - 3. [機能要件詳細](.agents/rules/spec-docs/03-functional-requirements.md)
  - 4. [Minecraftプラグイン APIインターフェース仕様](.agents/rules/spec-docs/04-api-specification.md)
  - 5. [CLIコマンド (AI用) 仕様](.agents/rules/spec-docs/05-cli-commands.md)
  - 6. [データモデル・JSON構造](.agents/rules/spec-docs/06-data-model.md)
  - 7. [非機能要件・その他](.agents/rules/spec-docs/07-non-functional-requirements.md)
