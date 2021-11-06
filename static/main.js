const GREEN_CIRCLE_EMOJI = '&#128994;';
const CHAT_MESSAGE = 'chat-message';

let socket;

function connectToRoom() {
    if (socket != null) {
        socket.close();
        document.getElementById('messages_area').value = '';
    }

    socket = new WebSocket('ws://localhost:8080/connectToRoom' +
        '?username='
        + document.getElementById('username').value
        + '&roomId='
        + document.getElementById('room_id').value);

    console.log('Attempting Connection...');

    socket.onopen = () => {
        console.log('Successfully Connected');
        document.getElementById('connection_status').innerHTML = GREEN_CIRCLE_EMOJI;
    };

    socket.onmessage = message => {
        let json = JSON.parse(message.data);
        if (json.event === CHAT_MESSAGE) {
            document.getElementById('messages_area').value += `${json.payload.username}: ${String(json.payload.text)}\r\n`
        }
    };

    socket.onclose = event => {
        console.log('Socket Closed Connection: ', event);
        document.getElementById('connection_status').innerHTML = event.type.toString();
    };

    socket.onerror = error => {
        console.log('Socket Error: ', error);
        document.getElementById('connection_status').innerHTML = error.type.toString();
    };
}

function sendMessage() {
    let message = document.getElementById('message_input').value.toString();

    if (socket == null || message === '') {
        return;
    }

    socket.send(message);
    document.getElementById('message_input').value = '';
}
