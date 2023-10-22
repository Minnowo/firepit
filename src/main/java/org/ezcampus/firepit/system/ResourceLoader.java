package org.ezcampus.firepit.system;

import java.nio.file.Path;

import org.tinylog.Logger;
import org.tinylog.configuration.Configuration;

public class ResourceLoader
{
    public static void loadTinyLogConfig()
    {
        try
        {
            Configuration.set("writer1", "file");
            Configuration.set("writer1.file", Path.of(GlobalSettings.Log_Dir, GlobalSettings.Log_File).toString());
            Configuration.set("writer1.format", "[{date: yyyy-MM-dd HH:mm:ss.SSS}] [{level}] {message}");
            Configuration.set("writer1.append", "true");
            Configuration.set("writer1.level", "trace");

            Configuration.set("writer2", "console");
            Configuration.set("writer2.format", "[{date: yyyy-MM-dd HH:mm:ss.SSS}] [{level}] {message}");

            if(GlobalSettings.IS_DEBUG)
            {
                Configuration.set("writer2.level", "trace");
            }
            else
            {
                Configuration.set("writer2.level", "info");
            }

            Logger.info("TinyLog has initialized!");
        }
        catch (UnsupportedOperationException e)
        {
            Logger.warn("Tried to update tinylog config, it was already set");
        }
    }
    
    
    public static void loadEnv() {
    	
    	GlobalSettings.Log_File = String.format("%d.log", System.currentTimeMillis());
    	
    	String log_file = System.getenv("LOG_FILE");
    	String log_dir = System.getenv("LOG_DIR");

    	String debug = System.getenv("IS_DEBUG");
    	
    	if(debug != null && !debug.isBlank()) {
    		GlobalSettings.IS_DEBUG = debug.toLowerCase().equals("true");
        	System.out.println(String.format("IS_DEBUG Loaded: %s", GlobalSettings.IS_DEBUG ));
    	}
    	
    	if(log_dir != null && !log_dir.isBlank()) {
    		GlobalSettings.Log_Dir = Path.of(log_dir).toString();
        	System.out.println(String.format("Log_Dir Loaded: %s", GlobalSettings.Log_Dir ));
    	}
    	
    	if(log_file != null && !log_file.isBlank()) {
    		GlobalSettings.Log_File = log_file;
        	System.out.println(String.format("Log_File Loaded: %s", GlobalSettings.Log_File ));
    	}
    }
}