package com.yadokari1130.mccli;

import io.javalin.Javalin;
import net.minecraft.server.MinecraftServer;
import org.slf4j.Logger;

/**
 * JavalinベースのHTTPサーバー管理クラス。
 * エンドポイントのルーティング登録と、サーバーのライフサイクル管理を担う。
 */
public class HttpServer {

    private static final int PORT = 8080;

    private final MinecraftServer minecraftServer;
    private final Logger logger;
    private Javalin app;

    public HttpServer(MinecraftServer minecraftServer, Logger logger) {
        this.minecraftServer = minecraftServer;
        this.logger = logger;
    }

    /**
     * HTTPサーバーを起動し、エンドポイントを登録する。
     */
    public void start() {
        // Javalinが正しいクラスローダーを使用するよう設定
        Thread.currentThread().setContextClassLoader(HttpServer.class.getClassLoader());

        BlockApiHandler blockApiHandler = new BlockApiHandler(minecraftServer);
        InteractApiHandler interactApiHandler = new InteractApiHandler(minecraftServer);
        InventoryApiHandler inventoryApiHandler = new InventoryApiHandler(minecraftServer);

        app = Javalin.create().start(PORT);

        // == エンドポイント登録 ==

        /** 指定範囲のブロック情報を取得する */
        app.get("/api/blocks", blockApiHandler::getBlocks);

        /** ブロックを一括配置する */
        app.post("/api/blocks", blockApiHandler::placeBlocks);

        /** FakePlayerによるインタラクトを実行する */
        app.post("/api/interact", interactApiHandler::interact);

        /** アイテムをドロップする */
        app.post("/api/drop-items", interactApiHandler::dropItems);

        /** インベントリの内容を設定する */
        app.post("/api/inventory", inventoryApiHandler::setInventory);

        // グローバル例外ハンドラ
        app.exception(Exception.class, (e, ctx) -> {
            logger.error("HTTPリクエスト処理中に予期せぬエラーが発生しました: {}", e.getMessage(), e);
            ctx.status(500).result("サーバー内部エラー: " + e.getMessage());
        });

        logger.info("HTTPサーバーをポート {} で起動しました。", PORT);
    }

    /**
     * HTTPサーバーを停止する。
     */
    public void stop() {
        if (app != null) {
            app.stop();
            logger.info("HTTPサーバーを停止しました。");
        }
    }
}
