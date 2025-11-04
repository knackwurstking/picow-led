function powerToggleBeforeRequest(event) {
	const target = event.currentTarget;
	target.disabled = true;

	const spinner = document.createElement("span");
	spinner.classList.add("spinner");

	target.append(spinner);

	console.debug("Power toggle before request", target);
}

function powerToggleAfterRequest(event) {
	const target = event.currentTarget;
	target.disabled = false;

	const allSpinners = target.querySelectorAll(".spinner");
	if (allSpinners.length > 0) {
		allSpinners.forEach(function (spinner) {
			spinner.remove();
		});
	}

	console.debug("Power toggle after request", target);
}

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
});
