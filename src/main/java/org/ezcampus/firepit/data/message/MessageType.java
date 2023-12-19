package org.ezcampus.firepit.data.message;

public interface MessageType
{

	public static final int SET_CLIENT_NAME = 10;
	public static final int SET_CLIENT_POSITION = 20;
	public static final int CLIENT_SET_SPEAKER = 30;
	public static final int CLIENT_LEAVE_ROOM = 40;
	public static final int CLIENT_JOIN_ROOM = 50;
	public static final int CLIENT_WHO_AM_I_MESSAGE = 100;
	
	public static final int ROOM_INFO = 60;
	

	public static final int SERVER_OK_MESSAGE = 200;
	public static final int SERVER_BAD_MESSAGE = 400;
	
}