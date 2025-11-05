// Keep details tag state (open/closed)
document.addEventListener("DOMContentLoaded", function () {
	const allDetails = document.querySelectorAll("details");
	allDetails.forEach(function (detail) {
		detail.addEventListener("toggle", function () {
			if (!detail.id) return;
			localStorage.setItem(detail.id, detail.open.toString());
		});

		const savedState = localStorage.getItem(detail.id);
		if (savedState !== null) {
			detail.open = savedState === "true";
		}
	});

	function triggerReloads() {
		document.body.dispatchEvent(new CustomEvent("reloadDevices"));
		document.body.dispatchEvent(new CustomEvent("reloadGroups"));
	}
	document.body.addEventListener("visible", triggerReloads);
	triggerReloads();
});
