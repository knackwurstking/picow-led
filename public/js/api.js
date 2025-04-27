var __async = (__this, __arguments, generator) => {
  return new Promise((resolve, reject) => {
    var fulfilled = (value) => {
      try {
        step(generator.next(value));
      } catch (e) {
        reject(e);
      }
    };
    var rejected = (value) => {
      try {
        step(generator.throw(value));
      } catch (e) {
        reject(e);
      }
    };
    var step = (x) => x.done ? resolve(x.value) : Promise.resolve(x.value).then(fulfilled, rejected);
    step((generator = generator.apply(__this, __arguments)).next());
  });
};
function getUrl() {
  return ``;
}
function color(color2, ...devices) {
  return __async(this, null, function* () {
    if (!color2) {
      color2 = [255, 255, 255, 255];
    }
    const url = getUrl() + "/api/devices/color";
    const data = { devices, color: color2 };
    console.debug(`POST "${url}":`, data);
    let status;
    fetch(url, {
      method: "POST",
      header: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(data)
    }).then((resp) => {
      status = resp.status;
      return resp.json();
    }).then((data2) => {
      console.debug(`fetch "${url}":`, data2);
      if ("message" in data2) {
        throw new Error(`${status}: ${data2.message}`);
      }
      console.warn(data2);
    }).catch((err) => {
      console.error(err);
    });
  });
}
window.api = {
  color
};
