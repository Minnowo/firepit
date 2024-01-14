package org.ezcampus.firepit.data;

import org.tinylog.Logger;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;


public class JsonService
{
	private static final ObjectMapper mapper  = new ObjectMapper();

	public static String toJson(Object o)
	{
		try
		{
			return mapper.writeValueAsString(o);
		}
		catch (JsonProcessingException e)
		{
			Logger.error(e);
			return null;
		}
	}

	public static <T> T fromJson(String json, Class<T> valueType)
	{
		try
		{
			return mapper.readValue(json, valueType);
		}
		catch (JsonProcessingException e)
		{
			Logger.error(e);
			return null;
		}
	}
}
