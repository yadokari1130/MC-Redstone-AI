package com.yadokari1130.mccli;

import net.fabricmc.api.ModInitializer;
import net.fabricmc.fabric.api.event.lifecycle.v1.ServerLifecycleEvents;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import io.javalin.Javalin;

public class MCCLI implements ModInitializer {
	public static final String MOD_ID = "mc-cli";
	public static final Logger LOGGER = LoggerFactory.getLogger(MOD_ID);
	private Javalin app;

	@Override
	public void onInitialize() {
		LOGGER.info("Initializing HTTP Server Mod!");

		ServerLifecycleEvents.SERVER_STARTED.register(server -> {
			LOGGER.info("Starting Javalin HTTP server on port 8080...");
			
			Thread.currentThread().setContextClassLoader(MCCLI.class.getClassLoader());

			app = Javalin.create().start(8080);
			
			app.get("/hello", (io.javalin.http.Context ctx) -> ctx.result("Hello, World!"));
		});

		ServerLifecycleEvents.SERVER_STOPPING.register(server -> {
			if (app != null) {
				LOGGER.info("Stopping Javalin HTTP server...");
				app.stop();
			}
		});
	}
}