function powerToggleBeforeRequest(event) {
	const target = event.currentTarget;
	target.disabled = true;
	target.classList.add("spinner");
}

function powerToggleAfterRequest(event) {
	const target = event.currentTarget;
	target.disabled = false;
	target.classList.remove("spinner");
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
