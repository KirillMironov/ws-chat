const GREEN_CIRCLE_EMOJI = "&#128994;";

let socket;

function connectToRoom() {
    if (socket != null) {
        socket.close();
        document.getElementById("messages_area").value = "";
    }

    socket = new WebSocket("ws://localhost:8080/connectToRoom?username="
        + document.getElementById("username").value);
    console.log("Attempting Connection...");

    socket.onopen = () => {
        console.log("Successfully Connected");
        document.getElementById("connection_status").innerHTML = GREEN_CIRCLE_EMOJI;
    };

    socket.onmessage = message => {
        document.getElementById("messages_area").value += message.data + "\r\n";
    };

    socket.onclose = event => {
        console.log("Socket Closed Connection: ", event);
        document.getElementById("connection_status").innerHTML = event.type.toString();
    };

    socket.onerror = error => {
        console.log("Socket Error: ", error);
        document.getElementById("connection_status").innerHTML = error.type.toString();
    };
}
