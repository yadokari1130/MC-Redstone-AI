import argparse
import os
import re
from pathlib import Path

# 定数定義
LOCAL_WORKFLOW_DIR = Path(".agents/workflows")
LOCAL_COMMAND_DIR = Path(".gemini/commands")
GLOBAL_WORKFLOW_DIR = Path.home() / ".gemini/antigravity/global_workflows"
GLOBAL_COMMAND_DIR = Path.home() / ".gemini/commands"

def parse_markdown(content):
    """
    MarkdownファイルからFrontmatterとコンテンツを抽出する。
    """
    frontmatter = {}
    remaining_content = content
    
    if content.startswith("---"):
        match = re.match(r"^---\s*\n(.*?)\n---\s*\n(.*)", content, re.DOTALL)
        if match:
            fm_text = match.group(1)
            remaining_content = match.group(2)
            
            # 簡易的なYAMLパース（descriptionのみ抽出）
            desc_match = re.search(r"description:\s*(.*)", fm_text)
            if desc_match:
                frontmatter["description"] = desc_match.group(1).strip()
    
    return frontmatter, remaining_content

def sync_directory(src_dir, dest_dir):
    """
    ディレクトリ内のMarkdownファイルをTOMLに変換して同期する。
    """
    if not src_dir.exists():
        print(f"警告: ソースディレクトリが見つかりません: {src_dir}")
        return

    if not dest_dir.exists():
        dest_dir.mkdir(parents=True, exist_ok=True)
        print(f"ディレクトリを作成しました: {dest_dir}")

    for md_file in src_dir.glob("*.md"):
        cmd_name = md_file.stem
        toml_file = dest_dir / f"{cmd_name}.toml"
        
        with open(md_file, "r", encoding="utf-8") as f:
            content = f.read()
        
        frontmatter, prompt_body = parse_markdown(content)
        description = frontmatter.get("description", f"{cmd_name} ワークフローを実行します")
        
        # TOMLの生成
        # 改行やダブルクォートを適切に処理するために、トリプルクォートを使用する
        # また、末尾に {{args}} を追加して、ユーザー入力を受け取れるようにする
        toml_content = f'description = "{description}"\n'
        toml_content += 'prompt = """\n'
        toml_content += prompt_body.strip()
        toml_content += '\n\nRequest: {{args}}\n"""\n'
        
        with open(toml_file, "w", encoding="utf-8") as f:
            f.write(toml_content)
        
        print(f"同期完了: {md_file} -> {toml_file}")

def main():
    parser = argparse.ArgumentParser(description="AntigravityワークフローをGemini CLIカスタムコマンドに同期します。")
    parser.add_argument("--local", action=argparse.BooleanOptionalAction, default=True, help="ローカルワークフローの同期 (デフォルト: True)")
    parser.add_argument("--global", action=argparse.BooleanOptionalAction, default=False, dest="sync_global", help="グローバルワークフローの同期 (デフォルト: False)")
    
    args = parser.parse_args()

    if args.local:
        # ローカルワークフローの同期
        print("ローカルワークフローを同期中...")
        sync_directory(LOCAL_WORKFLOW_DIR, LOCAL_COMMAND_DIR)
    
    if args.sync_global:
        # グローバルワークフローの同期
        print("\nグローバルワークフローを同期中...")
        sync_directory(GLOBAL_WORKFLOW_DIR, GLOBAL_COMMAND_DIR)

if __name__ == "__main__":
    main()
