package org.ezcampus.firepit.websocket;

import java.io.IOException;

import org.ezcampus.firepit.data.Client;
import org.ezcampus.firepit.data.JsonService;
import org.ezcampus.firepit.data.Room;
import org.ezcampus.firepit.data.RoomSessionController;
import org.ezcampus.firepit.data.message.BadMessage;
import org.ezcampus.firepit.data.message.ClientJoinRoomMessage;
import org.ezcampus.firepit.data.message.ClientLeaveRoomMessage;
import org.ezcampus.firepit.data.message.MessageType;
import org.ezcampus.firepit.data.message.OkMessage;
import org.ezcampus.firepit.data.message.RenameMessage;
import org.ezcampus.firepit.data.message.RoomInfoMessage;
import org.ezcampus.firepit.data.message.SetSpeakerMessage;
import org.ezcampus.firepit.data.message.SocketMessage;
import org.tinylog.Logger;

import jakarta.inject.Inject;
import jakarta.websocket.CloseReason;
import jakarta.websocket.CloseReason.CloseCodes;
import jakarta.websocket.OnClose;
import jakarta.websocket.OnMessage;
import jakarta.websocket.OnOpen;
import jakarta.websocket.Session;
import jakarta.websocket.server.PathParam;
import jakarta.websocket.server.ServerEndpoint;

@ServerEndpoint("/websocket/")
public class WSRoom {

	@Inject
	RoomSessionController roomSessionController;

	@Inject
	JsonService jsonService;

	Client client;

	Room clientRoom;

	@OnOpen
	public void onOpen(@PathParam("rid") String roomid,
			@PathParam("disp_name") String displayName,
			@PathParam("disp_occup") String displayOccupation,
			@PathParam("disp_avatar") int avatarIndex,
			Session session) throws IOException {

				System.out.println("Incoming Connection Details:\n" +
                   "------------------------------\n" +
                   "Room ID: " + roomid + "\n" +
                   "Display Name: " + displayName + "\n" +
                   "Display Occupation: " + displayOccupation + "\n" +
                   "Avatar Index: " + avatarIndex + "\n" +
                   "Session ID: " + session.getId() + "\n" +
                   "------------------------------");


		clientRoom = roomSessionController.getRoom(roomid);

		// ! --- ASSERT: Existence of Room ---
		if (clientRoom == null) {
			Logger.info("Someone tried to connect to a room which does not exist!");

			session.close(new CloseReason(CloseCodes.UNEXPECTED_CONDITION, "Room does not exist"));
			return;
		}

		// ! --- ASSERT: If Room has enough space for incoming client ---
		int spaceDiff = clientRoom.roomCapacity - clientRoom.roomMembers.size();
		if (spaceDiff == 0) {
			session.close(new CloseReason(CloseCodes.UNEXPECTED_CONDITION, "Room doesn't have enough space!"));
			return;
		}

		// * All Assertions Passed! Instanciate client & do assertions on their inputted data now

		client = new Client(session, displayName, displayOccupation);

		clientRoom.addClient(client);

		Logger.info("New client connection {}", session.getId());

		session.getBasicRemote().sendText(jsonService.toJson(new OkMessage("Room has been joined")));

		session.getBasicRemote().sendText(jsonService.toJson(new RoomInfoMessage(clientRoom)));

		ClientJoinRoomMessage m = new ClientJoinRoomMessage(client);
		clientRoom.broadCast(jsonService.toJson(m), session);
	}

	@OnClose
	public void onClose(CloseReason reason, Session session) {

		Logger.info("Client {} has disconnected", session.getId());

		if (clientRoom == null) {
			Logger.warn("Cannot broadcast because room is null");
			return;
		}

		ClientLeaveRoomMessage m = new ClientLeaveRoomMessage(client);

		clientRoom.broadCast(jsonService.toJson(m), session);

		clientRoom.removeClient(session.getId());
	}

	@OnMessage
	public String onMessage(String json, Session session) throws IOException {
		if (clientRoom == null) {

			session.close(new CloseReason(CloseCodes.UNEXPECTED_CONDITION, "Room does not exist"));

			return null;
		}

		SocketMessage m = jsonService.fromJson(json, SocketMessage.class);

		if (m == null) {
			return jsonService.toJson(new BadMessage("Invalid message request"));
		}

		switch (m.getMessageType()) {
			case MessageType.SERVER_BAD_MESSAGE:
			case MessageType.SERVER_OK_MESSAGE:
			case MessageType.CLIENT_JOIN_ROOM:
			case MessageType.CLIENT_LEAVE_ROOM:
				return jsonService.toJson(new BadMessage("Message does not make sense!"));

			case MessageType.CLIENT_SET_SPEAKER:

				SetSpeakerMessage ssm = (SetSpeakerMessage) m;

				if (!clientRoom.speaker.clientId.equals(session.getId())) {
					return jsonService.toJson(new BadMessage("You cannot set the speaker unless you are the speaker!"));
				}

				if (clientRoom.setSpeakerFromPublicId(ssm.newSpeakerId)) {

					clientRoom.broadCast(json, null);

					return jsonService.toJson(new OkMessage("The speaker has been changed"));
				}

				return jsonService.toJson(new BadMessage("The speaker id did not exist in the room"));

			case MessageType.SET_CLIENT_NAME:
			case MessageType.SET_CLIENT_POSITION:
				break;
		}

		Logger.info("Client {} said {}", session.getId(), m.getMessageType());

		return jsonService.toJson(new OkMessage());
	}

}
