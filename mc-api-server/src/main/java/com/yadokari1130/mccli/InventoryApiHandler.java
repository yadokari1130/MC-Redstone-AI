package com.yadokari1130.mccli;

import com.google.gson.Gson;
import io.javalin.http.Context;
import net.minecraft.core.BlockPos;
import net.minecraft.core.registries.BuiltInRegistries;
import net.minecraft.resources.Identifier;
import net.minecraft.server.MinecraftServer;
import net.minecraft.server.level.ServerLevel;
import net.minecraft.world.Container;
import net.minecraft.world.item.Item;
import net.minecraft.world.item.ItemStack;
import net.minecraft.world.level.block.entity.BlockEntity;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

/**
 * /api/inventory エンドポイントのハンドラ。
 * 指定された座標のブロックインベントリを操作する。
 */
public class InventoryApiHandler {

    private static final Gson GSON = new Gson();
    private static final int TIMEOUT_SECONDS = 10;

    private final MinecraftServer server;

    public InventoryApiHandler(MinecraftServer server) {
        this.server = server;
    }

    /**
     * POST /api/inventory
     * リクエストボディ: { "x": int, "y": int, "z": int, "items": [{ "id": "...", "amount": int }, ...] }
     */
    public void setInventory(Context ctx) throws Exception {
        String body = ctx.body();
        if (body == null || body.isBlank()) {
            ctx.status(400).result("リクエストボディが空です。");
            return;
        }

        InventoryRequest request;
        try {
            request = GSON.fromJson(body, InventoryRequest.class);
        } catch (Exception e) {
            ctx.status(400).result("リクエストJSONの解析に失敗しました: " + e.getMessage());
            return;
        }

        final InventoryRequest finalRequest = request;

        CompletableFuture<String> future = new CompletableFuture<>();
        server.execute(() -> {
            try {
                ServerLevel world = server.overworld();
                BlockPos pos = new BlockPos(finalRequest.x, finalRequest.y, finalRequest.z);
                BlockEntity blockEntity = world.getBlockEntity(pos);

                if (!(blockEntity instanceof Container container)) {
                    future.complete("ERROR:指定された座標のブロック（" + world.getBlockState(pos).getBlock().toString() + "）はインベントリを持っていません。");
                    return;
                }

                // インベントリをクリア
                container.clearContent();

                if (finalRequest.items != null) {
                    List<String> errors = new ArrayList<>();
                    int slot = 0;
                    int maxSlots = container.getContainerSize();

                    for (InventoryItem itemData : finalRequest.items) {
                        if (slot >= maxSlots) {
                            break; // インベントリ容量を超えたら終了
                        }

                        Identifier itemId = Identifier.tryParse(itemData.id);
                        if (itemId == null || !BuiltInRegistries.ITEM.containsKey(itemId)) {
                            errors.add("不明なアイテムID: " + itemData.id);
                            continue;
                        }

                        Optional<Item> itemOpt = BuiltInRegistries.ITEM.get(itemId).map(h -> h.value());
                        if (itemOpt.isEmpty()) {
                            errors.add("アイテムの取得に失敗: " + itemData.id);
                            continue;
                        }

                        ItemStack stack = new ItemStack(itemOpt.get(), itemData.amount);
                        container.setItem(slot, stack);
                        slot++;
                    }

                    container.setChanged();

                    if (errors.isEmpty()) {
                        future.complete("OK:インベントリを更新しました。");
                    } else {
                        future.complete("PARTIAL:一部のアイテム設定に失敗しました: " + String.join(", ", errors));
                    }
                } else {
                    container.setChanged();
                    future.complete("OK:インベントリを空にしました。");
                }
            } catch (Exception e) {
                future.completeExceptionally(e);
            }
        });

        String result = future.get(TIMEOUT_SECONDS, TimeUnit.SECONDS);
        if (result.startsWith("ERROR:")) {
            ctx.status(400).result(result.substring(6));
        } else if (result.startsWith("PARTIAL:")) {
            ctx.status(207).result(result.substring(8));
        } else {
            ctx.status(200).result(result.substring(3));
        }
    }

    private static class InventoryRequest {
        public int x;
        public int y;
        public int z;
        public List<InventoryItem> items;
    }

    private static class InventoryItem {
        public String id;
        public int amount;
    }
}
