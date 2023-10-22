package org.ezcampus.firepit.data;

import org.tinylog.Logger;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;

import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;
import jakarta.enterprise.context.ApplicationScoped;

@ApplicationScoped
public class JsonService
{
	private final ObjectMapper mapper  = new ObjectMapper();

	@PostConstruct
	void init()
	{

	}

	@PreDestroy
	void destroy()
	{
	}

	public String toJson(Object o)
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

	public <T> T fromJson(String json, Class<T> valueType)
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
