const API = 'http://localhost:8000';

// --- Estado ---
let token = localStorage.getItem('token');
let userID = localStorage.getItem('userID');
let username = localStorage.getItem('username');
let ws = null;
let currentRoomID = null;

// --- Elementos ---
const authSection = document.getElementById('auth-section');
const chatSection = document.getElementById('chat-section');
const loginView = document.getElementById('login-view');
const registerView = document.getElementById('register-view');
const forgotView = document.getElementById('forgot-view');
const msgInput = document.getElementById('msg-input');
const sendBtn = document.getElementById('send-btn');
const noRoomMsg = document.getElementById('no-room-msg');
const messagesEl = document.getElementById('messages');

// --- Navegación entre vistas ---
function showView(view) {
    [loginView, registerView, forgotView].forEach(v => v.classList.add('hidden'));
    view.classList.remove('hidden');
}

function showAuth() {
    authSection.classList.remove('hidden');
    chatSection.classList.add('hidden');
    showView(loginView);
}

function showChat() {
    authSection.classList.add('hidden');
    chatSection.classList.remove('hidden');
    document.getElementById('user-display').textContent = username;
}

// Inicialización: si hay sesión activa, mostrar el chat
if (token && userID) {
    showChat();
} else {
    showAuth();
}

document.getElementById('go-register').addEventListener('click', e => { e.preventDefault(); showView(registerView); });
document.getElementById('go-login').addEventListener('click', e => { e.preventDefault(); showView(loginView); });
document.getElementById('go-forgot').addEventListener('click', e => { e.preventDefault(); showView(forgotView); });
document.getElementById('go-login-2').addEventListener('click', e => { e.preventDefault(); showView(loginView); });

// --- Login ---
document.getElementById('login-form').addEventListener('submit', async e => {
    e.preventDefault();
    const usernameVal = document.getElementById('login-username').value.trim();
    const password = document.getElementById('login-password').value;
    const errorEl = document.getElementById('login-error');

    try {
        const res = await fetch(`${API}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username: usernameVal, password })
        });

        if (!res.ok) {
            errorEl.textContent = await res.text();
            errorEl.classList.remove('hidden');
            return;
        }

        const data = await res.json();
        token = data.token;

        // Decodificar JWT para obtener el userID (sub)
        const payload = JSON.parse(atob(token.split('.')[1]));
        userID = payload.sub || usernameVal;
        username = usernameVal;

        localStorage.setItem('token', token);
        localStorage.setItem('userID', userID);
        localStorage.setItem('username', username);

        errorEl.classList.add('hidden');
        showChat();
    } catch {
        errorEl.textContent = 'Error de conexión con el servidor';
        errorEl.classList.remove('hidden');
    }
});

// --- Registro ---
document.getElementById('register-form').addEventListener('submit', async e => {
    e.preventDefault();
    const usernameVal = document.getElementById('register-username').value.trim();
    const email = document.getElementById('register-email').value.trim();
    const password = document.getElementById('register-password').value;
    const errorEl = document.getElementById('register-error');
    const successEl = document.getElementById('register-success');

    try {
        const res = await fetch(`${API}/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username: usernameVal, email, password })
        });

        if (!res.ok) {
            errorEl.textContent = await res.text();
            errorEl.classList.remove('hidden');
            successEl.classList.add('hidden');
            return;
        }

        errorEl.classList.add('hidden');
        successEl.textContent = '¡Cuenta creada! Ahora inicia sesión.';
        successEl.classList.remove('hidden');
        setTimeout(() => showView(loginView), 1500);
    } catch {
        errorEl.textContent = 'Error de conexión con el servidor';
        errorEl.classList.remove('hidden');
    }
});

// --- Recuperar Contraseña ---
document.getElementById('forgot-form').addEventListener('submit', async e => {
    e.preventDefault();
    const email = document.getElementById('forgot-email').value.trim();
    const msgEl = document.getElementById('forgot-msg');

    try {
        await fetch(`${API}/forgot-password`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email })
        });
        msgEl.textContent = 'Si el email existe, recibirás un enlace de recuperación.';
        msgEl.classList.remove('hidden');
    } catch {
        msgEl.textContent = 'Error de conexión con el servidor';
        msgEl.classList.remove('hidden');
    }
});

// --- Logout ---
document.getElementById('logout-btn').addEventListener('click', () => {
    localStorage.clear();
    token = null; userID = null; username = null;
    if (ws) ws.close();
    ws = null;
    currentRoomID = null;
    showAuth();
});

// --- Unirse a sala ---
document.getElementById('join-room-btn').addEventListener('click', () => {
    const roomID = parseInt(document.getElementById('room-input').value);
    if (!roomID || roomID < 1) return;
    joinRoom(roomID);
});

document.getElementById('room-input').addEventListener('keydown', e => {
    if (e.key === 'Enter') {
        e.preventDefault();
        const roomID = parseInt(document.getElementById('room-input').value);
        if (roomID && roomID >= 1) joinRoom(roomID);
    }
});

function joinRoom(roomID) {
    if (ws) ws.close();

    currentRoomID = roomID;
    messagesEl.innerHTML = '';
    noRoomMsg.classList.add('hidden');

    document.getElementById('room-info').classList.remove('hidden');
    document.getElementById('current-room-label').textContent = `Sala #${roomID}`;

    const wsURL = `ws://localhost:8000/ws?room_id=${roomID}&user_id=${encodeURIComponent(userID)}`;
    ws = new WebSocket(wsURL);

    ws.onopen = () => {
        addSystemMessage(`Conectado a la sala #${roomID}`);
        msgInput.disabled = false;
        sendBtn.disabled = false;
        msgInput.focus();
    };

    ws.onmessage = event => {
        try {
            const msg = JSON.parse(event.data);
            const isMine = msg.sender_id === userID;
            addMessage(msg.sender_id, msg.content, isMine);
        } catch {
            addMessage('', event.data, false);
        }
    };

    ws.onclose = () => {
        addSystemMessage('Desconectado de la sala');
        msgInput.disabled = true;
        sendBtn.disabled = true;
    };

    ws.onerror = () => addSystemMessage('Error en la conexión WebSocket');
}

// --- Enviar mensaje ---
document.getElementById('msg-form').addEventListener('submit', e => {
    e.preventDefault();
    const msg = msgInput.value.trim();
    if (!msg || !ws || ws.readyState !== WebSocket.OPEN) return;
    ws.send(msg);
    msgInput.value = '';
});

// --- Helpers UI ---
function addMessage(sender, content, isMine) {
    const div = document.createElement('div');
    div.classList.add('message', isMine ? 'message-mine' : 'message-other');

    if (sender && !isMine) {
        const senderEl = document.createElement('span');
        senderEl.classList.add('sender');
        senderEl.textContent = sender;
        div.appendChild(senderEl);
    }

    const contentEl = document.createElement('span');
    contentEl.classList.add('content');
    contentEl.textContent = content;
    div.appendChild(contentEl);

    messagesEl.appendChild(div);
    messagesEl.scrollTop = messagesEl.scrollHeight;
}

function addSystemMessage(text) {
    const div = document.createElement('div');
    div.classList.add('message-system');
    div.textContent = text;
    messagesEl.appendChild(div);
    messagesEl.scrollTop = messagesEl.scrollHeight;
}
