package org.ezcampus.firepit.system;

import java.nio.file.Paths;

public class GlobalSettings
{
	public static final String BRAND = "firepit";
    public static final String BRAND_LONG = "SchedulePlatform-" + BRAND;
    
    public static boolean IS_DEBUG = false;


    public static String Log_Dir = Paths.get(".", "logs").toString();
    public static String Log_File = "";


}
