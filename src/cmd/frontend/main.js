

let ws;
const HOST = "localhost:3000";


function setSpeaker(speakerId) {

    const grid = document.getElementById("peopleGrid");
    const personElements = grid.getElementsByClassName("person");
    
        console.log(speakerId)
     for (const personElement of personElements) {
        console.log(personElement)
         if (personElement.dataset.clientId === speakerId) {
            const speakerH = document.getElementById("currentSpeaker");
            
            speakerH.textContent = personElement.textContent;
             break;
         }
     }
}

 function removePerson(clientId) {
    const grid = document.getElementById("peopleGrid");
    const personElements = grid.getElementsByClassName("person");
    
     for (const personElement of personElements) {
         if (personElement.dataset.clientId === clientId) {
             console.log(`person with id ${clientId} was removed`);
             grid.removeChild(personElement);
             break;
         }
     }
}

function addPersonToGrid(person) {
    
    const grid = document.getElementById("peopleGrid");
    const personElement = document.createElement("div");
    personElement.classList.add("person");
    personElement.textContent = person.client_name;
    personElement.dataset.clientId = person.client_id;
    grid.appendChild(personElement);
}

function populatePeopleGrid(data) {
    data.forEach(person => {
        addPersonToGrid(person);
    });
}


function getOKMessage(){
    return JSON.stringify({
        messageType: 200
    });
}

function joinRoom(){
   var a = document.getElementById("roomId"); 
    
   socket_connect(a.value);
}

function newRoom() {
    // Define your payload here
    let payload = {
        room_name: "YourRoomName", // Replace with the desired room name
        room_capacity: 10 // Replace with the desired room capacity
    };

    fetch(`http://${HOST}/room/new`, {
        method: 'GET',
    })
    .then(response => response.text())
    .then(roomId => {
        socket_connect(roomId);
    })
    .catch(error => console.error('Error:', error));
}

function socket_connect(roomId){

    if(ws && ws.readyState === WebSocket.OPEN){
        console.log("already in a room")
        return;
    }

    if(ws) {
        ws.close();
    }
    const d = document.getElementById("clientName");
    
    console.log("Got room id from server " + roomId);

    var url = "ws://localhost:3000/ws";
    url += "?client_name=";
    url += d.value;
    url += "&client_id=nyah";
    // url += "&client_occupation=test";
    url += "&rid=";
    url += roomId;

    ws = new WebSocket(url);

    ws.onopen = function (event) {
        console.log("websocket open");
        ws.send(getOKMessage());
    }

    // parse messages received from the server and update the UI accordingly
    ws.onmessage = function (event) {
        console.log("room data: " + event.data);
        

        const json = JSON.parse(event.data);

        if(!json) {
            return;
        }

        console.log(json);

        
        // room info message
        if(json.messageType === 60) {
            
            populatePeopleGrid(json.payload.room.room_members);
            
            setSpeaker(json.payload.room.room_speaker.client_id);
        }

        
        // client joins room
        if(json.messageType === 50) {
            
            addPersonToGrid(json.payload.client);
        }

        // client leaves room
        if(json.messageType === 40) {
            
            removePerson(json.payload.client.client_id);
        }
        
        // set speaker 
        if(json.messageType === 30) {
            
            console.log("Setting speaker to " + json.payload.speaker_id);
            setSpeaker(json.payload.speaker_id);
        }
    }
}

function sendMessage(){

    if(ws && ws.readyState === WebSocket.OPEN){
        ws.send("hello server!");
    }
}

function sendSetSpeakerMessage(){
    
    const d = document.getElementById("newSpeakerId");
    console.log("setting speaker to " + d.value);
    
    if(ws && ws.readyState === WebSocket.OPEN){
        console.log("Sending message");
        ws.send(JSON.stringify({
            messageType : 30,
            speaker_id : d.value
        }));
    }
    else {
        console.log("socket is closed");
    }
}

function sendOKMessage(){

    if(ws && ws.readyState === WebSocket.OPEN){
        console.log("Sending message");
        ws.send(getOKMessage());
    }
    else {
        console.log("socket is closed");
    }
}
function sendBadMessage(){

    if(ws && ws.readyState === WebSocket.OPEN){
        console.log("Sending message");
        ws.send("hello");
    }
    else {
        console.log("socket is closed");
    }
}


(function(){
    
    // fetch(`http://${HOST}`, {
    //     method: 'GET',
    // })
    //     .then(x => x.json())
    //     .then(json => {
        
    //         console.log(json);

    // });
    

    
    console.log("Client Open");
}())