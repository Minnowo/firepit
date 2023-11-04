

let ws;
const HOST = "localhost:8085/firepit";


function setSpeaker(speakerId) {

    const grid = document.getElementById("peopleGrid");
    const personElements = grid.getElementsByClassName("person");
    
     for (const personElement of personElements) {
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

function newRoom(){
    
    fetch(`http://${HOST}/room/new`)
        .then(x => x.text())
        .then(roomId => {
            socket_connect(roomId);
    });
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

    ws = new WebSocket(`ws://${HOST}/websocket/?rid=${roomId}&name=${d.value}`);

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
            
            populatePeopleGrid(json.room.room_members);
            
            setSpeaker(json.room.room_speaker);
        }

        
        // client joins room
        if(json.messageType === 50) {
            
            addPersonToGrid(json.client);
        }

        // client leaves room
        if(json.messageType === 40) {
            
            removePerson(json.client.client_id);
        }
        
        // set speaker 
        if(json.messageType === 30) {
            
            console.log("Setting speaker to " + json.speaker_id);
            setSpeaker(json.speaker_id);
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
    
    fetch(`http://${HOST}`, {
        method: 'GET',
    })
        .then(x => x.json())
        .then(json => {
        
            console.log(json);

    });
    

    
    console.log("Client Open");
}())
