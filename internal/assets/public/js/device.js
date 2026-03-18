// The SERVER_PATH_PREFIX is defined inside the "layout.templ" components script section

function changeDeviceColor(addr, deviceID, hex) {
	console.debug(`Changing color for device with ${addr} to`, hex);

	query = new URLSearchParams({
		color: hex,
	});
	url = (SERVER_PATH_PREFIX || "") + `/api/devices/${deviceID}/color?${query.toString()}`;

	fetch(url, { method: "POST" })
		.then((response) => {
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			return response.json();
		})
		.then((data) => {
			console.debug("Device color changed successfully:", data);
		})
		.catch((error) => {
			console.error("Error changing device color:", error);
		});
}

function changeDeviceWhite(addr, value) {
	// TODO: Change the device color via API call
	console.debug(`Changing white for device with ${addr} to`, value);
}

function changeDeviceWhite2(addr, value) {
	// TODO: Change the device color via API call
	console.debug(`Changing white2 for device with ${addr} to`, value);
}
