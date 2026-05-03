package com.yadokari1130.mccli;

import com.google.gson.Gson;
import io.javalin.http.Context;
import net.minecraft.core.registries.BuiltInRegistries;
import net.minecraft.resources.Identifier;
import net.minecraft.server.MinecraftServer;
import net.minecraft.world.entity.Entity;
import net.minecraft.world.level.Level;
import net.minecraft.world.phys.AABB;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

/**
 * POST /api/kill-entities のハンドラ。
 * 指定範囲内のエンティティを削除する。
 * すべてのワールド操作はCompletableFutureを通じてメインスレッドで実行する。
 */
public class EntityApiHandler {

    private static final Gson GSON = new Gson();
    /** タイムアウト秒数 */
    private static final int TIMEOUT_SECONDS = 10;

    private final MinecraftServer server;

    public EntityApiHandler(MinecraftServer server) {
        this.server = server;
    }

    /**
     * POST /api/kill-entities
     * リクエストボディ: { "x1": ..., "y1": ..., "z1": ..., "x2": ..., "y2": ..., "z2": ..., "type": "..." }
     * type はオプション。未指定時は範囲内のすべてのエンティティを削除。
     */
    public void killEntities(Context ctx) throws Exception {
        String body = ctx.body();
        if (body == null || body.isBlank()) {
            ctx.status(400).result("リクエストボディが空です。{\"x1\":...,\"y1\":...,\"z1\":...,\"x2\":...,\"y2\":...,\"z2\":...} 形式で送信してください。");
            return;
        }

        KillEntitiesRequest request;
        try {
            request = GSON.fromJson(body, KillEntitiesRequest.class);
        } catch (Exception e) {
            ctx.status(400).result("リクエストJSONの解析に失敗しました: " + e.getMessage());
            return;
        }

        // 必須パラメータのバリデーション
        if (request.x1 == 0 && request.y1 == 0 && request.z1 == 0 &&
            request.x2 == 0 && request.y2 == 0 && request.z2 == 0) {
            // デフォルト値のままかもしれないので、明示的なエラーにはしない
        }

        // 座標範囲を正規化
        final int minX = Math.min(request.x1, request.x2);
        final int minY = Math.min(request.y1, request.y2);
        final int minZ = Math.min(request.z1, request.z2);
        final int maxX = Math.max(request.x1, request.x2);
        final int maxY = Math.max(request.y1, request.y2);
        final int maxZ = Math.max(request.z1, request.z2);

        // typeのフィルタがあれば準備
        final String typeFilter = request.type != null && !request.type.isBlank() ? request.type : null;
        final Identifier typeId;
        if (typeFilter != null) {
            typeId = Identifier.tryParse(typeFilter);
            if (typeId == null) {
                ctx.status(400).result("無効なエンティティタイプ識別子です: '" + typeFilter + "'");
                return;
            }
        } else {
            typeId = null;
        }

        // メインスレッドでエンティティ削除を実行
        CompletableFuture<KillEntitiesResult> future = new CompletableFuture<>();
        server.execute(() -> {
            try {
                Level world = server.overworld();
                List<String> removedEntities = new ArrayList<>();
                int removedCount = 0;
                int skippedCount = 0;

                AABB aabb = new AABB(minX, minY, minZ, maxX + 1, maxY + 1, maxZ + 1);

                for (Entity entity : world.getEntities(null, aabb)) {
                    // typeフィルタがあれば一致しないものはスキップ
                    if (typeId != null) {
                        Identifier entityTypeId = BuiltInRegistries.ENTITY_TYPE.getKey(entity.getType());
                        if (entityTypeId == null || !entityTypeId.equals(typeId)) {
                            skippedCount++;
                            continue;
                        }
                    }

                    entity.discard();
                    removedCount++;
                    removedEntities.add(entity.getUUID().toString());
                }

                future.complete(new KillEntitiesResult(removedCount, skippedCount, removedEntities));
            } catch (Exception e) {
                future.completeExceptionally(e);
            }
        });

        KillEntitiesResult result = future.get(TIMEOUT_SECONDS, TimeUnit.SECONDS);
        ctx.contentType("application/json").result(GSON.toJson(result));
    }

    /**
     * POST /api/kill-entities のリクエストボディ用内部クラス。
     */
    private static class KillEntitiesRequest {
        public int x1;
        public int y1;
        public int z1;
        public int x2;
        public int y2;
        public int z2;
        /** 削除対象のエンティティタイプ（オプション） */
        public String type;
    }

    /**
     * POST /api/kill-entities のレスポンス用内部クラス。
     */
    private static class KillEntitiesResult {
        /** 削除されたエンティティの数 */
        public int removed;
        /** タイプフィルタによりスキップされたエンティティの数 */
        public int skipped;
        /** 削除されたエンティティのUUIDリスト */
        public List<String> entities;

        public KillEntitiesResult(int removed, int skipped, List<String> entities) {
            this.removed = removed;
            this.skipped = skipped;
            this.entities = entities;
        }
    }
}
