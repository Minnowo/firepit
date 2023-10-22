

let ws;
const HOST = "localhost:8085/firepit";

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
    
    console.log("Got room id from server " + roomId);

    ws = new WebSocket(`ws://${HOST}/websocket/?rid=${roomId}`);

    ws.onopen = function (event) {
        console.log("websocket open");
        ws.send(getOKMessage());
    }

    // parse messages received from the server and update the UI accordingly
    ws.onmessage = function (event) {
        console.log(event.data);
    }
}

function sendMessage(){

    if(ws && ws.readyState === WebSocket.OPEN){
        ws.send("hello server!");
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
