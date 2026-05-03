package com.yadokari1130.mccli;

import java.util.Map;

/**
 * エンティティデータのDTO。
 * JSON形式: { "uuid": "...", "type": "minecraft:...", "x": ..., "y": ..., "z": ..., "yaw": ..., "pitch": ..., "nbt": {...} }
 */
public class EntityData {
    /** エンティティのUUID */
    public String uuid;
    /** エンティティタイプ識別子 (例: "minecraft:boat") */
    public String type;
    /** X座標（小数可） */
    public double x;
    /** Y座標（小数可） */
    public double y;
    /** Z座標（小数可） */
    public double z;
    /** 水平方向の回転角度 */
    public float yaw;
    /** 垂直方向の回転角度 */
    public float pitch;
    /**
     * エンティティのNBTデータ（JSON形式のマップ）。
     * クライアントとの送受信にはJSONを使用し、サーバー内部ではCompoundTagに変換する。
     */
    public Map<String, Object> nbt;

    public EntityData() {}

    public EntityData(String uuid, String type, double x, double y, double z, float yaw, float pitch, Map<String, Object> nbt) {
        this.uuid = uuid;
        this.type = type;
        this.x = x;
        this.y = y;
        this.z = z;
        this.yaw = yaw;
        this.pitch = pitch;
        this.nbt = nbt;
    }
}
