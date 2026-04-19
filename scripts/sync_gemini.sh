#!/bin/bash

# AntigravityのSkillsとWorkflowsをGemini CLIに同期するスクリプト

# エラーが発生した時点でスクリプトを中断
set -e

# プロジェクトのルートディレクトリに移動（スクリプトの場所に関わらず）
cd "$(dirname "$0")/.."

echo "--- Workflowsの同期を開始します ---"
if [ -f "scripts/sync_workflows.py" ]; then
    python3 scripts/sync_workflows.py
else
    echo "エラー: scripts/sync_workflows.py が見つかりません。"
    exit 1
fi

echo -e "\n--- Skillsの同期を開始します ---"
# gemini skills link は指定したディレクトリ内の各スキルディレクトリ（SKILL.mdを含むもの）を認識します
# --scope workspace を指定することで、プロジェクトローカルな .gemini/skills/ にリンクを作成します

CATEGORIES=("blocks" "concepts" "mc_cli")

for CATEGORY in "${CATEGORIES[@]}"; do
    TARGET_DIR=".agents/skills/$CATEGORY"
    if [ -d "$TARGET_DIR" ]; then
        echo "カテゴリー '$CATEGORY' を同期中..."
        gemini skills link "$TARGET_DIR" --scope workspace
    else
        echo "警告: ディレクトリ $TARGET_DIR が見つかりません。スキップします。"
    fi
done

echo -e "\n同期が完了しました。Geminiセッション内で /skills reload および /commands reload を実行して変更を反映してください。"
