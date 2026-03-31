package com.yadokari1130.mccli;

import net.fabricmc.api.ModInitializer;
import net.fabricmc.fabric.api.event.lifecycle.v1.ServerLifecycleEvents;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class MCCLI implements ModInitializer {
	public static final String MOD_ID = "mc-cli";
	public static final Logger LOGGER = LoggerFactory.getLogger(MOD_ID);
	private HttpServer httpServer;

	@Override
	public void onInitialize() {
		LOGGER.info("MC-CLI Mod を初期化しています...");

		ServerLifecycleEvents.SERVER_STARTED.register(server -> {
			LOGGER.info("HTTPサーバーを起動します...");
			httpServer = new HttpServer(server, LOGGER);
			httpServer.start();
		});

		ServerLifecycleEvents.SERVER_STOPPING.register(server -> {
			if (httpServer != null) {
				httpServer.stop();
			}
		});
	}
}