

let ws;
const HOST = "localhost:8085/firepit";


function socket_connect(){

    if(ws) {
        ws.close();
    }
    
    fetch(`http://${HOST}/room/new`)
        .then(x => x.text())
        .then(roomId => {

            ws = new WebSocket(`ws://${HOST}/websocket/?rid=${roomId}`);

            ws.onopen = function (event) {
                console.log("websocket open");
                ws.send("hello server!");
            }

            // parse messages received from the server and update the UI accordingly
            ws.onmessage = function (event) {
                console.log(event.data);
            }
    });

    
}

function sendMessage(){
    if(ws && ws.readyState === WebSocket.OPEN){
        ws.send("hello server!");
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
    
    socket_connect();

    
    console.log("Client Open");
}())
