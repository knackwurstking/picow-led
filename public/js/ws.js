let socket = null;
let timeout = null;
let timeoutDuration = 1e3;
const onClose = function() {
  if (timeout !== null) {
    clearTimeout(timeout);
    timeout = null;
  }
  timeout = setTimeout(() => {
    connect();
  }, timeoutDuration);
};
const onOpen = () => {
};
const onMessage = () => {
};
function addr() {
  return ``;
}
function isOpen() {
  if (!socket) return false;
  return socket.readyState === socket.OPEN;
}
function connect() {
  if (socket) close();
  const wsAddr = addr();
  console.debug(`Try to connect WebSocket to ${wsAddr}`);
  socket = new WebSocket(wsAddr);
  socket.addEventListener("close", onClose);
  socket.addEventListener("open", onOpen);
  socket.addEventListener("message", onMessage);
}
function close() {
  if (timeout) {
    clearTimeout(timeout);
    timeout = null;
  }
  if (socket) {
    socket.removeEventListener("close", onClose);
    if (isOpen()) socket.close();
    socket = null;
  }
}
window.ws = {
  addr,
  isOpen,
  connect,
  close
};
