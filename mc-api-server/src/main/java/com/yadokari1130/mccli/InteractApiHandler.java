package com.yadokari1130.mccli;

import com.google.gson.Gson;
import com.mojang.authlib.GameProfile;
import io.javalin.http.Context;
import net.minecraft.core.BlockPos;
import net.minecraft.core.Direction;
import net.minecraft.server.MinecraftServer;
import net.minecraft.server.level.ClientInformation;
import net.minecraft.server.level.ServerLevel;
import net.minecraft.server.level.ServerPlayer;
import net.minecraft.world.InteractionResult;
import net.minecraft.world.level.block.state.BlockState;
import net.minecraft.world.phys.BlockHitResult;
import net.minecraft.world.phys.Vec3;

import java.util.UUID;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

/**
 * POST /api/interact のハンドラ。
 * FakePlayerパターンを使用して、バニラのインタラクト操作を完全にエミュレートする。
 * BlockStateの直接書き換えは行わず、ブロックの useWithoutItem メソッドを呼び出すことで
 * サウンド再生やボタン自動復帰などのバニラメカニクスを正確に再現する。
 */
public class InteractApiHandler {

    private static final Gson GSON = new Gson();
    /** タイムアウト秒数 */
    private static final int TIMEOUT_SECONDS = 10;

    /** FakePlayerに使用する固定UUID */
    private static final UUID FAKE_PLAYER_UUID = UUID.fromString("a1b2c3d4-e5f6-7890-abcd-ef1234567890");
    /** FakePlayerのゲーム内名称 */
    private static final String FAKE_PLAYER_NAME = "MC-CLI-FakePlayer";

    private final MinecraftServer server;

    public InteractApiHandler(MinecraftServer server) {
        this.server = server;
    }

    /**
     * POST /api/interact
     * リクエストボディ: { "x": x座標, "y": y座標, "z": z座標 }
     */
    public void interact(Context ctx) throws Exception {
        String body = ctx.body();
        if (body == null || body.isBlank()) {
            ctx.status(400).result("リクエストボディが空です。{\"x\":...,\"y\":...,\"z\":...} 形式で送信してください。");
            return;
        }

        InteractRequest request;
        try {
            request = GSON.fromJson(body, InteractRequest.class);
        } catch (Exception e) {
            ctx.status(400).result("リクエストJSONの解析に失敗しました: " + e.getMessage());
            return;
        }

        final InteractRequest finalRequest = request;

        // メインスレッドでFakePlayerによるインタラクトを実行
        CompletableFuture<InteractionResult> future = new CompletableFuture<>();
        server.execute(() -> {
            try {
                ServerLevel world = server.overworld();
                BlockPos pos = new BlockPos(finalRequest.x, finalRequest.y, finalRequest.z);
                BlockState blockState = world.getBlockState(pos);

                // FakePlayerを生成
                ServerPlayer fakePlayer = createFakePlayer(world);

                // ブロック中心から少し上にヒット位置を設定し、BlockHitResultを構築
                // Direction.UP = ブロック上面からのクリックをエミュレート
                Vec3 hitPos = Vec3.atCenterOf(pos).add(0, 0.5, 0);
                BlockHitResult hitResult = new BlockHitResult(hitPos, Direction.UP, pos, false);

                // バニラのuseWithoutItemメソッドを呼び出してインタラクトをエミュレート
                // これはプレイヤーが素手（MainHand）でブロックを右クリックした際のロジックと同等
                InteractionResult result = blockState.useWithoutItem(world, fakePlayer, hitResult);

                future.complete(result);
            } catch (Exception e) {
                future.completeExceptionally(e);
            }
        });

        InteractionResult result;
        try {
            result = future.get(TIMEOUT_SECONDS, TimeUnit.SECONDS);
        } catch (Exception e) {
            ctx.status(500).result("インタラクト処理中にエラーが発生しました: " + e.getMessage());
            return;
        }

        // consumesAction() == true ならインタラクト成功（SUCCESSまたはCONSUME）
        if (result.consumesAction()) {
            ctx.status(200).result("インタラクト成功: " + result.getClass().getSimpleName());
        } else {
            // PASS または FAIL の場合は対象ブロックがインタラクト不能と判断
            ctx.status(400).result("インタラクト失敗: 対象ブロックはインタラクト不能です。(InteractionResult=" + result.getClass().getSimpleName() + ")");
        }
    }

    /**
     * FakePlayerを生成する。
     * Carpet Modで実績のある手法に基づき、固定UUID/GameProfileで ServerPlayer を生成する。
     * このFakePlayerはMinecraftのコアから「本物のプレイヤー」として認識されるため、
     * サウンド再生・ブロック更新・Tick遅延処理などのバニラメカニクスがすべて正常に動作する。
     */
    private ServerPlayer createFakePlayer(ServerLevel world) {
        GameProfile profile = new GameProfile(FAKE_PLAYER_UUID, FAKE_PLAYER_NAME);
        return new ServerPlayer(
                server,
                world,
                profile,
                ClientInformation.createDefault()
        );
    }

    /** POST /api/interact のリクエストボディ用内部クラス */
    private static class InteractRequest {
        /** X座標 */
        public int x;
        /** Y座標 */
        public int y;
        /** Z座標 */
        public int z;
    }
}
