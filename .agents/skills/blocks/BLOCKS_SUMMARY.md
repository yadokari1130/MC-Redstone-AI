# レッドストーン・ブロック・要約リファレンス

このドキュメントは、Minecraft のレッドストーン回路で使用される主要なブロックの ID、重要なプロパティ、および動作特性を要約したものです。詳細な情報は各ディレクトリの `SKILL.md` を参照してください。

## 共通プロパティ
- `facing`: ブロックの向き (`north`, `south`, `east`, `west`, `up`, `down`)
- `powered`: レッドストーン信号を受けているか (`true`, `false`)
- `lit`: 発光しているか (ランプなど, `true`, `false`)

---

## 主要コンポーネント

### 1. 動力源・入力
| ブロック | ID | 主要プロパティ | 特徴 |
| :--- | :--- | :--- | :--- |
| **レバー** | `minecraft:lever` | `face` (wall, floor, ceiling), `facing` | インタラクトで ON/OFF 切り替え。 |
| **ボタン** | `minecraft:stone_button` 等 | `face`, `facing`, `powered` | 一定時間信号を出力。 |
| **RSブロック** | `minecraft:redstone_block` | - | 常に信号強度15を出力する不透明ブロック。 |

### 2. 伝達・制御
| ブロック | ID | 主要プロパティ | 特徴 |
| :--- | :--- | :--- | :--- |
| **RSダスト** | `minecraft:redstone_wire` | `power` (0-15), `north/south/east/west` (none, side, up) | 信号を伝達。1ブロックごとに強度が1減少。 |
| **リピーター** | `minecraft:repeater` | `delay` (1-4 RSチック), `facing`, `locked` | 信号の増幅（強度15にリセット）、遅延、一方向伝達。 |
| **コンパレータ** | `minecraft:comparator` | `mode` (compare, subtract), `facing`, `powered` | 信号の比較、減算、インベントリ検知。 |
| **RSトーチ** | `minecraft:redstone_torch` | `lit` | 入力信号を反転（NOT回路）。上および横に強度15を出力。 |

### 3. 出力・アクション
| ブロック | ID | 主要プロパティ | 特徴 |
| :--- | :--- | :--- | :--- |
| **RSランプ** | `minecraft:redstone_lamp` | `lit` | 信号を受けると点灯。 |
| **ピストン** | `minecraft:piston` | `extended`, `facing` | ブロックを押し出す。通常/粘着 (`sticky_piston`) がある。 |
| **オブザーバー** | `minecraft:observer` | `facing`, `powered` | 正面のブロック更新を検知し、背面からパルスを出力。 |

### 4. その他
| ブロック | ID | 主要プロパティ | 特徴 |
| :--- | :--- | :--- | :--- |
| **ホッパー** | `minecraft:hopper` | `enabled`, `facing` | アイテムの移動。信号を受けると停止 (`enabled=false`)。 |
| **土台ブロック** | `minecraft:smooth_stone` 等 | - | 導電性のあるフルブロックを土台として使用推奨。 |

---

## 設計時の注意点
1. **出力形式**: CLI のデフォルト出力は `[BlockID, [X, Y, Z], {Properties}]` の配列形式です。
2. **座標管理**: 回路を配置する際は、土台となるブロックの上に各コンポーネントを配置してください。
3. **信号強度**: 長距離の配線にはリピーターが必要です。
