package com.yadokari1130.mccli;

import com.google.gson.Gson;
import com.google.gson.reflect.TypeToken;
import io.javalin.http.Context;
import net.minecraft.core.BlockPos;
import net.minecraft.core.registries.BuiltInRegistries;
import net.minecraft.nbt.*;
import net.minecraft.resources.Identifier;
import net.minecraft.server.MinecraftServer;
import net.minecraft.util.ProblemReporter;
import net.minecraft.world.entity.Entity;
import net.minecraft.world.entity.EntitySpawnReason;
import net.minecraft.world.entity.EntityType;
import net.minecraft.world.level.Level;
import net.minecraft.world.level.block.Block;
import net.minecraft.world.level.block.entity.BlockEntity;
import net.minecraft.world.level.block.state.BlockState;
import net.minecraft.world.level.block.state.properties.Property;
import net.minecraft.world.level.redstone.Orientation;
import net.minecraft.world.level.storage.TagValueInput;
import net.minecraft.world.level.storage.TagValueOutput;
import net.minecraft.world.level.storage.ValueInput;
import net.minecraft.world.phys.AABB;

import java.lang.reflect.Type;
import java.util.*;
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
     * オプション: include_entities (boolean, デフォルト: false)
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

        String includeEntitiesParam = ctx.queryParam("include_entities");
        boolean includeEntities = "true".equalsIgnoreCase(includeEntitiesParam);

        // 座標範囲を正規化
        final int minX = Math.min(x1, x2);
        final int minY = Math.min(y1, y2);
        final int minZ = Math.min(z1, z2);
        final int maxX = Math.max(x1, x2);
        final int maxY = Math.max(y1, y2);
        final int maxZ = Math.max(z1, z2);

        // メインスレッドでワールド操作を実行
        CompletableFuture<Object> future = new CompletableFuture<>();
        server.execute(() -> {
            try {
                Level world = server.overworld();
                List<BlockData> blockResult = new ArrayList<>();

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

                            blockResult.add(new BlockData(x, y, z, blockIdStr, stateMap, nbtString));
                        }
                    }
                }

                if (includeEntities) {
                    List<EntityData> entityResult = new ArrayList<>();
                    AABB aabb = new AABB(minX, minY, minZ, maxX + 1, maxY + 1, maxZ + 1);
                    for (Entity entity : world.getEntities(null, aabb)) {
                        // プレイヤーエンティティは除外
                        if (entity.getType() == EntityType.PLAYER) {
                            continue;
                        }

                        Identifier entityTypeId = BuiltInRegistries.ENTITY_TYPE.getKey(entity.getType());
                        String typeStr = entityTypeId != null ? entityTypeId.toString() : "minecraft:unknown";

                        // NBTデータを取得してJSON形式に変換
                        ProblemReporter.Collector problems = new ProblemReporter.Collector();
                        TagValueOutput output = TagValueOutput.createWithContext(problems, server.registryAccess());
                        entity.saveWithoutId(output);
                        CompoundTag nbt = output.buildResult();
                        Map<String, Object> nbtMap = (Map<String, Object>) nbtToObject(nbt);
                        // 不要な大きなフィールドを除外して軽量化
                        nbtMap.remove("Pos");
                        nbtMap.remove("Rotation");
                        nbtMap.remove("UUID");

                        entityResult.add(new EntityData(
                                entity.getUUID().toString(),
                                typeStr,
                                entity.getX(),
                                entity.getY(),
                                entity.getZ(),
                                entity.getYRot(),
                                entity.getXRot(),
                                nbtMap
                        ));
                    }
                    future.complete(new BlocksAndEntitiesResponse(blockResult, entityResult));
                } else {
                    future.complete(blockResult);
                }
            } catch (Exception e) {
                future.completeExceptionally(e);
            }
        });

        Object result = future.get(TIMEOUT_SECONDS, TimeUnit.SECONDS);
        ctx.contentType("application/json").result(GSON.toJson(result));
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
        List<EntityData> entities = new ArrayList<>();
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
            // オブジェクト形式 { "blocks": [...], "entities": [...], "flags": 2 }
            PlaceBlocksRequest request;
            try {
                request = GSON.fromJson(body, PlaceBlocksRequest.class);
            } catch (Exception e) {
                ctx.status(400).result("リクエストJSONの解析に失敗しました: " + e.getMessage());
                return;
            }
            blocks = request.blocks;
            if (request.entities != null) {
                entities = request.entities;
            }
            if (request.flags != 0) {
                flags = request.flags;
            }
        }

        if ((blocks == null || blocks.isEmpty()) && (entities == null || entities.isEmpty())) {
            ctx.status(400).result("blocksおよびentitiesが両方とも空です。");
            return;
        }

        final List<BlockData> finalBlocks = blocks != null ? blocks : new ArrayList<>();
        final List<EntityData> finalEntities = entities;
        final int finalFlags = flags;

        // メインスレッドでブロック配置とエンティティスポーンを実行
        CompletableFuture<String> future = new CompletableFuture<>();
        server.execute(() -> {
            try {
                Level world = server.overworld();
                List<String> errors = new ArrayList<>();
                List<BlockPos> placedPositions = new ArrayList<>();
                int successBlockCount = 0;
                int successEntityCount = 0;

                // 1. ブロック配置
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

                    // 設置可能かチェック (プレイヤー設置と同様の生存条件チェック)
                    if (!state.canSurvive(world, pos)) {
                        errors.add(String.format("(%d, %d, %d): ブロック '%s' はこの位置に設置できません。", blockData.x, blockData.y, blockData.z, blockData.block));
                        continue;
                    }

                    world.setBlock(pos, state, finalFlags);
                    placedPositions.add(pos);
                    successBlockCount++;
                }

                // 全ブロック配置後に handleNeighborChanged を発火させて接続状態を自然に矯正する
                for (BlockPos p : placedPositions) {
                    BlockState placedState = world.getBlockState(p);
                    if (!placedState.isAir()) {
                        Block placedBlock = placedState.getBlock();
                        placedState.handleNeighborChanged(world, p, placedBlock, (Orientation) null, false);
                    }
                }

                // 2. エンティティスポーン
                for (EntityData entityData : finalEntities) {
                    if (entityData.type == null || entityData.type.isBlank()) {
                        errors.add("エンティティタイプがnullまたは空です。");
                        continue;
                    }

                    Identifier entityId = Identifier.tryParse(entityData.type);
                    if (entityId == null || !BuiltInRegistries.ENTITY_TYPE.containsKey(entityId)) {
                        errors.add(String.format("不明なエンティティ識別子 '%s'", entityData.type));
                        continue;
                    }

                    Optional<EntityType<?>> entityTypeOpt = BuiltInRegistries.ENTITY_TYPE.get(entityId).map(h -> h.value());
                    if (entityTypeOpt.isEmpty()) {
                        errors.add(String.format("エンティティタイプの取得に失敗しました: '%s'", entityData.type));
                        continue;
                    }

                    EntityType<?> entityType = entityTypeOpt.get();
                    Entity entity = entityType.create(world, EntitySpawnReason.COMMAND);
                    if (entity == null) {
                        errors.add(String.format("エンティティ '%s' の生成に失敗しました。", entityData.type));
                        continue;
                    }

                    entity.setPos(entityData.x, entityData.y, entityData.z);
                    entity.setYRot(entityData.yaw);
                    entity.setXRot(entityData.pitch);

                    // NBTデータがあれば適用
                    if (entityData.nbt != null && !entityData.nbt.isEmpty()) {
                        try {
                            CompoundTag nbt = (CompoundTag) objectToNbt(entityData.nbt);
                            // PosとRotationは直接設定した値を優先するため上書き
                            nbt.remove("Pos");
                            nbt.remove("Rotation");
                            ProblemReporter.Collector loadProblems = new ProblemReporter.Collector();
                            ValueInput input = TagValueInput.create(loadProblems, server.registryAccess(), nbt);
                            entity.load(input);
                            // load後に座標と回転を再設定
                            entity.setPos(entityData.x, entityData.y, entityData.z);
                            entity.setYRot(entityData.yaw);
                            entity.setXRot(entityData.pitch);
                        } catch (Exception e) {
                            errors.add(String.format("エンティティ '%s' のNBT適用に失敗しました: %s", entityData.type, e.getMessage()));
                        }
                    }

                    world.addFreshEntity(entity);
                    successEntityCount++;
                }

                if (errors.isEmpty()) {
                    future.complete(String.format("OK: ブロック%d個、エンティティ%d個を配置しました。", successBlockCount, successEntityCount));
                } else {
                    future.complete(String.format("PARTIAL: ブロック%d個・エンティティ%d個配置成功、エラー: %s", successBlockCount, successEntityCount, String.join("; ", errors)));
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
     * NBTタグをJSON対応のJavaオブジェクト（Map, List, Number, String）に変換する。
     */
    private static Object nbtToObject(Tag tag) {
        if (tag instanceof CompoundTag compound) {
            Map<String, Object> map = new HashMap<>();
            for (String key : compound.keySet()) {
                Tag child = compound.get(key);
                if (child != null) {
                    map.put(key, nbtToObject(child));
                }
            }
            return map;
        } else if (tag instanceof ListTag list) {
            List<Object> result = new ArrayList<>();
            for (int i = 0; i < list.size(); i++) {
                result.add(nbtToObject(list.get(i)));
            }
            return result;
        } else if (tag instanceof ByteArrayTag arr) {
            List<Byte> result = new ArrayList<>();
            for (byte b : arr.getAsByteArray()) {
                result.add(b);
            }
            return result;
        } else if (tag instanceof IntArrayTag arr) {
            List<Integer> result = new ArrayList<>();
            for (int v : arr.getAsIntArray()) {
                result.add(v);
            }
            return result;
        } else if (tag instanceof LongArrayTag arr) {
            List<Long> result = new ArrayList<>();
            for (long v : arr.getAsLongArray()) {
                result.add(v);
            }
            return result;
        } else if (tag instanceof ByteTag t) {
            return t.byteValue();
        } else if (tag instanceof ShortTag t) {
            return t.shortValue();
        } else if (tag instanceof IntTag t) {
            return t.intValue();
        } else if (tag instanceof LongTag t) {
            return t.longValue();
        } else if (tag instanceof FloatTag t) {
            return t.floatValue();
        } else if (tag instanceof DoubleTag t) {
            return t.doubleValue();
        } else if (tag instanceof StringTag t) {
            return t.value(); // StringTagはrecordなのでvalue()で文字列を取得
        }
        return tag.toString();
    }

    /**
     * Javaオブジェクト（Map, List, Number, String）をNBTタグに変換する。
     */
    private static Tag objectToNbt(Object obj) {
        if (obj instanceof Map<?, ?> map) {
            CompoundTag compound = new CompoundTag();
            for (Map.Entry<?, ?> entry : map.entrySet()) {
                String key = entry.getKey().toString();
                Tag value = objectToNbt(entry.getValue());
                compound.put(key, value);
            }
            return compound;
        } else if (obj instanceof List<?> list) {
            ListTag listTag = new ListTag();
            for (Object item : list) {
                listTag.add(objectToNbt(item));
            }
            return listTag;
        } else if (obj instanceof Number num) {
            if (num instanceof Byte) return ByteTag.valueOf(num.byteValue());
            if (num instanceof Short) return ShortTag.valueOf(num.shortValue());
            if (num instanceof Integer) return IntTag.valueOf(num.intValue());
            if (num instanceof Long) return LongTag.valueOf(num.longValue());
            if (num instanceof Float) return FloatTag.valueOf(num.floatValue());
            return DoubleTag.valueOf(num.doubleValue());
        } else if (obj instanceof String str) {
            return StringTag.valueOf(str);
        } else if (obj instanceof Boolean b) {
            return ByteTag.valueOf((byte) (b ? 1 : 0));
        }
        return StringTag.valueOf(obj.toString());
    }

    /**
     * POST /api/blocks のリクエストボディ用内部クラス（オブジェクト形式）。
     */
    private static class PlaceBlocksRequest {
        /** 配置するブロックのリスト */
        public List<BlockData> blocks;
        /** 配置するエンティティのリスト */
        public List<EntityData> entities;
        /**
         * 更新フラグ。
         * 3 = ブロック更新 + クライアント同期（バニラ通常挙動）
         * 2 = クライアント同期のみ（高速一括配置用）
         */
        public int flags = 3;
    }

    /**
     * GET /api/blocks?include_entities=true 時のレスポンス用内部クラス。
     */
    private static class BlocksAndEntitiesResponse {
        public List<BlockData> blocks;
        public List<EntityData> entities;

        public BlocksAndEntitiesResponse(List<BlockData> blocks, List<EntityData> entities) {
            this.blocks = blocks;
            this.entities = entities;
        }
    }
}
