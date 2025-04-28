function setOnlineIndicator(state) {
  const el = document.querySelector(`.online-indicator`);
  if (state) {
    el.setAttribute(`data-state`, "online");
  } else {
    el.setAttribute(`data-state`, "offline");
  }
}
function powerButtonClickHandler(_ev, device) {
  let color;
  if (!device.color || !device.color.find((c) => c > 0)) {
    color = [255, 255, 255, 255];
  } else {
    color = [0, 0, 0, 0];
  }
  window.api.color(color, device);
}
function registerServiceWorker(serverPathPrefix) {
  if (!("serviceWorker" in navigator)) {
    console.warn("Browser doesn't support service workers");
    return;
  }
  window.addEventListener("load", function() {
    navigator.serviceWorker.register(serverPathPrefix + "/service-worker.js").then(function(reg) {
      console.info("Service worker registered", reg);
    }).catch(function(err) {
      console.error("Service worker registration failed:", err);
    });
  });
}
window.utils = {
  setOnlineIndicator,
  powerButtonClickHandler,
  registerServiceWorker
};
