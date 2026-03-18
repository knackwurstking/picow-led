// The SERVER_PATH_PREFIX is defined inside the "layout.templ" components script section

function changeDeviceColor(deviceID, hex) {
	console.debug(`Changing color for device with ID ${deviceID} to`, hex);

	query = new URLSearchParams({
		color: hex,
	});
	url = (SERVER_PATH_PREFIX || "") + `/api/devices/${deviceID}/color?${query.toString()}`;

	fetch(url, { method: "POST" })
		.then((response) => {
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			if (response.status === 204) {
				// No content, return an empty object
				return {};
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

function changeDeviceWhite(deviceID, value) {
	console.debug(`Changing white for device with ID ${deviceID} to`, value);

	query = new URLSearchParams({
		white: value,
	});
	url = (SERVER_PATH_PREFIX || "") + `/api/devices/${deviceID}/white?${query.toString()}`;

	fetch(url, { method: "POST" })
		.then((response) => {
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			if (response.status === 204) {
				// No content, return an empty object
				return {};
			}
			return response.json();
		})
		.then((data) => {
			console.debug("Device white level changed successfully:", data);
		})
		.catch((error) => {
			console.error("Error changing device white level:", error);
		});
}

function changeDeviceWhite2(deviceID, value) {
	// TODO: Change the device color via API call
	console.debug(`Changing white2 for device with ID ${deviceID} to`, value);
}
