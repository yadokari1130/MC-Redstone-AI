package com.yadokari1130.mccli;

import java.util.Map;

/**
 * ブロックデータのDTO。
 * JSON形式: { "x": ..., "y": ..., "z": ..., "block": "minecraft:...", "state": {...}, "nbt": "..." }
 */
public class BlockData {
    /** X座標 */
    public int x;
    /** Y座標 */
    public int y;
    /** Z座標 */
    public int z;
    /** ブロック識別子 (例: "minecraft:redstone_repeater") */
    public String block;
    /** BlockStateプロパティのKey-Valueマップ */
    public Map<String, String> state;
    /**
     * BlockEntityのNBTデータ (SNBT形式の文字列)。
     * BlockEntityが存在しない場合はnull。
     */
    public String nbt;

    public BlockData() {}

    public BlockData(int x, int y, int z, String block, Map<String, String> state, String nbt) {
        this.x = x;
        this.y = y;
        this.z = z;
        this.block = block;
        this.state = state;
        this.nbt = nbt;
    }
}
