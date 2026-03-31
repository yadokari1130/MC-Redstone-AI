package com.yadokari1130.mccli;

import com.google.gson.Gson;
import com.google.gson.reflect.TypeToken;
import io.javalin.http.Context;
import net.minecraft.core.BlockPos;
import net.minecraft.core.registries.BuiltInRegistries;
import net.minecraft.nbt.CompoundTag;
import net.minecraft.resources.Identifier;
import net.minecraft.server.MinecraftServer;
import net.minecraft.world.level.Level;
import net.minecraft.world.level.block.Block;
import net.minecraft.world.level.block.entity.BlockEntity;
import net.minecraft.world.level.block.state.BlockState;
import net.minecraft.world.level.block.state.properties.Property;

import java.lang.reflect.Type;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

/**
 * GET /api/blocks および POST /api/blocks のハンドラ。
 * すべてのワールド操作はCompletableFutureを通じてメインスレッドで実行する。
 */
public class BlockApiHandler {

    private static final Gson GSON = new Gson();
    /** タイムアウト秒数 */
    private static final int TIMEOUT_SECONDS = 10;

    private final MinecraftServer server;

    public BlockApiHandler(MinecraftServer server) {
        this.server = server;
    }

    /**
     * GET /api/blocks
     * クエリパラメータ: x1, y1, z1, x2, y2, z2 (すべて整数)
     * 空気ブロック (minecraft:air) は除外して返す。
     */
    public void getBlocks(Context ctx) throws Exception {
        // クエリパラメータのバリデーション
        int x1, y1, z1, x2, y2, z2;
        try {
            x1 = Integer.parseInt(ctx.queryParam("x1"));
            y1 = Integer.parseInt(ctx.queryParam("y1"));
            z1 = Integer.parseInt(ctx.queryParam("z1"));
            x2 = Integer.parseInt(ctx.queryParam("x2"));
            y2 = Integer.parseInt(ctx.queryParam("y2"));
            z2 = Integer.parseInt(ctx.queryParam("z2"));
        } catch (NumberFormatException e) {
            ctx.status(400).result("クエリパラメータ x1, y1, z1, x2, y2, z2 はすべて整数で指定してください。");
            return;
        }

        // 座標範囲を正規化
        final int minX = Math.min(x1, x2);
        final int minY = Math.min(y1, y2);
        final int minZ = Math.min(z1, z2);
        final int maxX = Math.max(x1, x2);
        final int maxY = Math.max(y1, y2);
        final int maxZ = Math.max(z1, z2);

        // メインスレッドでワールド操作を実行
        CompletableFuture<List<BlockData>> future = new CompletableFuture<>();
        server.execute(() -> {
            try {
                Level world = server.overworld();
                List<BlockData> result = new ArrayList<>();

                for (int x = minX; x <= maxX; x++) {
                    for (int y = minY; y <= maxY; y++) {
                        for (int z = minZ; z <= maxZ; z++) {
                            BlockPos pos = new BlockPos(x, y, z);
                            BlockState state = world.getBlockState(pos);

                            // 空気ブロックは除外
                            if (state.isAir()) {
                                continue;
                            }

                            // ブロック識別子を取得 (Identifierをそのままtostring)
                            Identifier blockId = BuiltInRegistries.BLOCK.getKey(state.getBlock());
                            String blockIdStr = blockId != null ? blockId.toString() : "minecraft:unknown";

                            // BlockStateのプロパティをMapに変換
                            Map<String, String> stateMap = new HashMap<>();
                            state.getValues().forEach(propValue -> {
                                // ワイルドカードキャプチャを回避するためにヘルパーメソッドを経由する
                                String[] nameAndValue = getPropertyNameAndValue(propValue);
                                stateMap.put(nameAndValue[0], nameAndValue[1]);
                            });

                            // BlockEntityのNBTデータを取得
                            String nbtString = null;
                            BlockEntity blockEntity = world.getBlockEntity(pos);
                            if (blockEntity != null) {
                                CompoundTag nbt = blockEntity.saveWithoutMetadata(server.registryAccess());
                                nbtString = nbt.toString();
                            }

                            result.add(new BlockData(x, y, z, blockIdStr, stateMap, nbtString));
                        }
                    }
                }

                future.complete(result);
            } catch (Exception e) {
                future.completeExceptionally(e);
            }
        });

        List<BlockData> blocks = future.get(TIMEOUT_SECONDS, TimeUnit.SECONDS);
        ctx.contentType("application/json").result(GSON.toJson(blocks));
    }

    /**
     * POST /api/blocks
     * リクエストボディ: BlockDataの配列 (JSON)、またはPlaceBlocksRequest形式のJSON
     * オプション: "flags" フィールドで更新フラグを指定 (デフォルト: 3)
     */
    public void placeBlocks(Context ctx) throws Exception {
        String body = ctx.body();
        if (body == null || body.isBlank()) {
            ctx.status(400).result("リクエストボディが空です。BlockDataの配列をJSON形式で送信してください。");
            return;
        }

        // リクエストをパース（配列形式またはオブジェクト形式の両方に対応）
        List<BlockData> blocks;
        int flags = 3; // デフォルト: バニラ挙動 (ブロック更新+クライアント同期)
        String trimmed = body.trim();
        if (trimmed.startsWith("[")) {
            // 配列形式
            Type listType = new TypeToken<List<BlockData>>() {}.getType();
            try {
                blocks = GSON.fromJson(body, listType);
            } catch (Exception e) {
                ctx.status(400).result("リクエストJSONの解析に失敗しました: " + e.getMessage());
                return;
            }
        } else {
            // オブジェクト形式 { "blocks": [...], "flags": 2 }
            PlaceBlocksRequest request;
            try {
                request = GSON.fromJson(body, PlaceBlocksRequest.class);
            } catch (Exception e) {
                ctx.status(400).result("リクエストJSONの解析に失敗しました: " + e.getMessage());
                return;
            }
            blocks = request.blocks;
            if (request.flags != 0) {
                flags = request.flags;
            }
        }

        if (blocks == null || blocks.isEmpty()) {
            ctx.status(400).result("blocksが空です。");
            return;
        }

        final List<BlockData> finalBlocks = blocks;
        final int finalFlags = flags;

        // メインスレッドでブロック配置を実行
        CompletableFuture<String> future = new CompletableFuture<>();
        server.execute(() -> {
            try {
                Level world = server.overworld();
                List<String> errors = new ArrayList<>();
                int successCount = 0;

                for (BlockData blockData : finalBlocks) {
                    if (blockData.block == null || blockData.block.isBlank()) {
                        errors.add(String.format("(%d, %d, %d): block識別子がnullまたは空です。", blockData.x, blockData.y, blockData.z));
                        continue;
                    }

                    // ブロック識別子からブロックを解決
                    Identifier blockId = Identifier.tryParse(blockData.block);
                    if (blockId == null || !BuiltInRegistries.BLOCK.containsKey(blockId)) {
                        errors.add(String.format("(%d, %d, %d): 不明なブロック識別子 '%s'", blockData.x, blockData.y, blockData.z, blockData.block));
                        continue;
                    }

                    Optional<Block> blockOpt = BuiltInRegistries.BLOCK.get(blockId).map(h -> h.value());
                    if (blockOpt.isEmpty()) {
                        errors.add(String.format("(%d, %d, %d): ブロックの取得に失敗しました: '%s'", blockData.x, blockData.y, blockData.z, blockData.block));
                        continue;
                    }

                    Block block = blockOpt.get();
                    BlockState state = block.defaultBlockState();

                    // BlockStateのプロパティを適用
                    if (blockData.state != null) {
                        for (Map.Entry<String, String> entry : blockData.state.entrySet()) {
                            state = applyProperty(state, entry.getKey(), entry.getValue());
                        }
                    }

                    BlockPos pos = new BlockPos(blockData.x, blockData.y, blockData.z);
                    world.setBlock(pos, state, finalFlags);
                    successCount++;
                }

                if (errors.isEmpty()) {
                    future.complete("OK: " + successCount + "個のブロックを配置しました。");
                } else {
                    future.complete("PARTIAL: " + successCount + "個配置成功、エラー: " + String.join("; ", errors));
                }
            } catch (Exception e) {
                future.completeExceptionally(e);
            }
        });

        String result = future.get(TIMEOUT_SECONDS, TimeUnit.SECONDS);
        if (result.startsWith("PARTIAL")) {
            ctx.status(207).result(result);
        } else {
            ctx.status(200).result(result);
        }
    }

    /**
     * BlockStateに指定プロパティを適用する。
     * getPossibleValues() を使って文字列名を照合し、合致した値を setValue する。
     * 無効なプロパティ名/値の場合はそのままのstateを返す。
     */
    private static BlockState applyProperty(BlockState state, String key, String value) {
        for (Property<?> property : state.getProperties()) {
            if (!property.getName().equals(key)) {
                continue;
            }
            // ワイルドカードキャプチャを回避するためにヘルパーメソッドを経由する
            BlockState applied = trySetPropertyValue(state, property, value);
            if (applied != null) {
                return applied;
            }
        }
        return state;
    }

    /**
     * ジェネリクスパラメータを確定させてワイルドカードキャプチャを回避するヘルパー。
     * 値の文字列名が一致した場合は setValue した新しい BlockState を返し、
     * 一致しなかった場合は null を返す。
     */
    @SuppressWarnings("unchecked")
    private static <T extends Comparable<T>> BlockState trySetPropertyValue(
            BlockState state, Property<T> property, String valueName) {
        for (T possible : property.getPossibleValues()) {
            if (property.getName(possible).equals(valueName)) {
                return state.setValue(property, possible);
            }
        }
        return null;
    }

    /**
     * Property.Value<?> からプロパティ名と値名を取り出すヘルパー（ワイルドカードキャプチャ回避用）。
     * @return [0] = プロパティ名, [1] = 値名
     */
    private static <T extends Comparable<T>> String[] getPropertyNameAndValue(
            Property.Value<T> propValue) {
        Property<T> prop = propValue.property();
        T val = propValue.value();
        return new String[]{prop.getName(), prop.getName(val)};
    }

    /**
     * POST /api/blocks のリクエストボディ用内部クラス（オブジェクト形式）。
     */
    private static class PlaceBlocksRequest {
        /** 配置するブロックのリスト */
        public List<BlockData> blocks;
        /**
         * 更新フラグ。
         * 3 = ブロック更新 + クライアント同期（バニラ通常挙動）
         * 2 = クライアント同期のみ（高速一括配置用）
         */
        public int flags = 3;
    }
}
