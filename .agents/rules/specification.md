---
trigger: model_decision
description: Redstone CLI全体の仕様書の目次と概要（各章の詳細はspec-docsディレクトリを参照）
---

# Redstone CLI 詳細仕様書

本仕様書は、Redstone CLIの開発に必要な詳細要件を各章ごとに定義したものです。
各章の詳細は `spec-docs` ディレクトリ内のファイルを参照してください。

## 目次

1. [はじめに](./spec-docs/01-introduction.md)
   - 背景・目的・プロジェクトスコープ
2. [システムアーキテクチャ・技術スタック](./spec-docs/02-architecture.md)
   - システム全体の構成図、使用言語、動作環境
3. [機能要件詳細](./spec-docs/03-functional-requirements.md)
   - 情報取得、ブロック配置、インタラクト操作の実現要件
4. [Minecraftプラグイン APIインターフェース仕様](./spec-docs/04-api-specification.md)
   - プラグイン側のHTTPエンドポイント定義 (`GET /api/blocks`, `POST /api/blocks`, `POST /api/interact`)
5. [CLIコマンド (AI用) 仕様](./spec-docs/05-cli-commands.md)
   - AIが利用するコマンドの定義 (`get_blocks`, `place_blocks`, `interact_block`)
6. [データモデル・JSON構造](./spec-docs/06-data-model.md)
   - 送受信されるブロックデータスキーマの定義 (`BlockData`)
7. [非機能要件・その他](./spec-docs/07-non-functional-requirements.md)
   - パフォーマンス、エラーハンドリング、セキュリティ事項
