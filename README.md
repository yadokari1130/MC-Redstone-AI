# Redstone AI

このプロジェクトは、AIを活用してMinecraft内で完全に自動化されたレッドストーン回路を設計・生成・動作検証するための統合リポジトリです。

[![](https://dcbadge.limes.pink/api/server/https://discord.com/invite/DSWjEq3NN2)](https://discord.com/invite/DSWjEq3NN2)

## 概要

レッドストーン回路の構築には複雑な論理設計と立体的なブロック配置が求められます。本プロジェクトは、AIにそれらの回路設計や配置を理解させ、Minecraft上で直接実行させるためのサーバー群、連携ツール、および各種ワークフローを提供するプロジェクトです。
最終的に、回路の要件定義から具体的なシステムの設計、ワールド上への配置、さらには動作確認に至るまでの一連のプロセスをAIが自律的かつ正確に行うことを目指しています。
情報収集・要件整理から回路設計・実装・検証まで、サブエージェントによるフェーズ分け実行を行うワークフローも備えています。

## 実行環境・対応エージェント

本プロジェクトは **Gemini CLI** を前提として設計・動作確認を行っていますが、もし Gemini CLI 以外のエージェントの方が性能が良さそうであれば、ぜひ教えてください。

現在、以下のエージェントで動作するよう構成しています。
- **Gemini CLI**（主な実行環境）
- **Claude Code**
- **Codex**
- **OpenCode**
- **Antigravity**

## リポジトリ構成

本リポジトリはモノレポ構成を採用しており、AIによる回路自動生成を実現するための複数のコンポーネントが統合されています。

### 1. AIリソース群
AIエージェントに自律的な思考とレッドストーン設計の専門知識を与えるためのリソースです。
- **`.gemini/agents/`**: カスタムエージェント定義。情報収集・要件整理を行う `circuit_researcher`、回路の設計・実装・検証を一貫して行う `circuit_builder` などが定義されています。
- **`.gemini/commands/`**: カスタムコマンド定義。特定の回路生成タスクをワンコマンドで実行できるようにまとめられています。
- **`.gemini/skills/`**: Agent Skills 定義。AIが `mc-cli` を通じてMinecraftの世界を操作するためのツール定義や、座標の扱い、ブロック状態(BlockState)の取得・配置、論理ゲートやクロック回路などのレッドストーン専門知識が含まれます。

### 2. Minecraftサーバー側 (Fabric Mod / HTTPサーバー)
Minecraftサーバー(Fabric)上で動作するHTTPサーバーです。Javalinを使用してREST APIを提供し、ワールドのブロック情報を取得・配置・操作します。サーバーのメインスレッドで安全に処理を行い、FakePlayerを使用してバニラ同様のブロック操作（インタラクト）をエミュレートします。

### 3. CLIツール (Go / Cobra)
AIエージェントとMinecraftサーバーの間に立ち、中継および制御を行うCLIツールです。Go言語とCobraライブラリで実装されており、AIは提供される操作用コマンドを駆使して、安全かつ正確に環境を操作します。

提供されているコマンド:
- `get-blocks`: 指定範囲のブロック情報を取得
- `place-blocks`: レッドストーン回路などの構造物を配置
- `interact-block`: レバーやボタンなどのブロックを操作
- `fill`: 指定範囲を特定のブロックで埋める
- `drop-items`: アイテムをエンティティとしてドロップ
- `set-inventory`: コンテナブロックのインベントリにアイテムをセット
- `test`: 回路の動作をYAMLテストファイルで検証

### 4. その他のディレクトリ
- **`scripts/`**: 各種スクリプトを置くディレクトリです。

## 詳細仕様書

システムの具体的な要件やアーキテクチャについては、以下の仕様書ドキュメント群にまとめられています。

- [総合仕様書・目次](.agents/rules/specification.md)
  1. [はじめに](.agents/rules/spec-docs/01-introduction.md)
  2. [システムアーキテクチャ・技術スタック](.agents/rules/spec-docs/02-architecture.md)
  3. [機能要件詳細](.agents/rules/spec-docs/03-functional-requirements.md)
  4. [Minecraftプラグイン APIインターフェース仕様](.agents/rules/spec-docs/04-api-specification.md)
  5. [CLIコマンド (AI用) 仕様](.agents/rules/spec-docs/05-cli-commands.md)
  6. [データモデル・JSON構造](.agents/rules/spec-docs/06-data-model.md)
  7. [非機能要件・その他](.agents/rules/spec-docs/07-non-functional-requirements.md)
