

let ws;
let jwt=function(){};
const HOST = "localhost:3000";

function loginAccount(){

    let pass = document.getElementById("passwordLoginId").value;
    let user = document.getElementById("usernameLoginId").value;

    fetch(`http://${HOST}/auth/token`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            username: user,
            password: pass
        }),
    }).then(r => r.json()).
        then(json => {

            console.log(json);
            jwt.token = json.token;
        });
}

function createAccount(){

    let pass = document.getElementById("passwordId").value;
    let user = document.getElementById("usernameId").value;

    fetch(`http://${HOST}/auth/create`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            username: user,
            password: pass
        }),
    }).then(r => r.json()).
        then(json => {

            console.log(json);
            jwt.token = json.token;
        });
}

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
    const btn = document.createElement("button");
    btn.innerHTML = "make speaker";
    btn.onclick = function(){setSpeakerById(person.client_id)}
    personElement.classList.add("person");
    personElement.textContent = person.client_name;
    personElement.dataset.clientId = person.client_id;
    personElement.appendChild(btn);
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
    var b = document.getElementById("reconnectToken"); 

    socket_connect(a.value, b.value);
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
            socket_connect(roomId, 0);
        })
        .catch(error => console.error('Error:', error));
}

function socket_connect(roomId, reconnectToken){

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
    url += "?name=";
    url += d.value;

    url += "&rid=";
    url += roomId;

    url += "&rtoken=";
    url += reconnectToken;

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

        if(json.messageType === 100) {
            document.getElementById("reconnectToken").value = json.payload.reconnection_token
        }

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

function postNewQuote(){
    quote = document.getElementById("newQuoteId").value;

    console.log(quote);

    fetch(`http://${HOST}/authed/quote`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${jwt.token}`,
        },
        body: JSON.stringify({
            quote: quote,
        }),
    }).catch(e => console.log(e));
}

function setSpeakerById(id){

    if(ws && ws.readyState === WebSocket.OPEN){
        console.log("Sending message");
        ws.send(JSON.stringify({
            messageType : 30,
            payload: {
                speaker_id : id
            }
        }));
    }
    else {
        console.log("socket is closed");
    }
}
function sendSetSpeakerMessage(){

    const d = document.getElementById("newSpeakerId");
    console.log("setting speaker to " + d.value);
    setSpeakerById(d.value);

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
