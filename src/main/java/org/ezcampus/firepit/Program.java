package org.ezcampus.firepit;

import java.io.IOException;
import java.lang.Thread.UncaughtExceptionHandler;

import org.ezcampus.firepit.system.GlobalSettings;
import org.ezcampus.firepit.system.ResourceLoader;
import org.glassfish.jersey.server.ResourceConfig;
import org.tinylog.Logger;

import jakarta.ws.rs.ApplicationPath;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.container.ContainerRequestContext;
import jakarta.ws.rs.container.ContainerResponseContext;
import jakarta.ws.rs.container.ContainerResponseFilter;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;
import jakarta.ws.rs.ext.Provider;

@ApplicationPath("/")
public class Program extends ResourceConfig
{
	private static class TinylogHandler implements UncaughtExceptionHandler
	{
		@Override
		public void uncaughtException(Thread thread, Throwable ex)
		{
			Logger.error("! === Unhandled Exception === !");
			Logger.error(ex);
		}
	}


	private void onShutdown()
	{
		Logger.info("Shut down hook triggered...");

	}


	public Program() throws IOException
	{
		GlobalSettings.IS_DEBUG = true;
		
		ResourceLoader.loadEnv();

		ResourceLoader.loadTinyLogConfig();
		Thread.setDefaultUncaughtExceptionHandler(new TinylogHandler());

		Logger.info("{} starting...", GlobalSettings.BRAND_LONG);
		Logger.info("Running as debug: {}", GlobalSettings.IS_DEBUG);
		
		this.packages("org.ezcampus.firepit.api");
		this.packages("org.ezcampus.firepit.websocket");
		
		// Enable CORS
		this.register(CORSFilter.class);
		
		Runtime.getRuntime().addShutdownHook(new Thread(this::onShutdown));
	}

	@Provider
	public static class CORSFilter implements ContainerResponseFilter
	{
		@Override
		public void filter(ContainerRequestContext requestContext, ContainerResponseContext responseContext)
				throws IOException
		{

			// Please note that setting Access-Control-Allow-Origin to "*" allows requests
			// from any origin.
			responseContext.getHeaders().add("Access-Control-Allow-Origin", "*");

			responseContext.getHeaders().add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE");
			responseContext.getHeaders().add("Access-Control-Allow-Headers", "Content-Type");
		}
	}

}