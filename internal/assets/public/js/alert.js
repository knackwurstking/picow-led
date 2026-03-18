(function() {
	var alertContainer = document.querySelectorAll(`[data-alert-container]`)

	alertContainer.forEach(container => {
		// Watch this container if anything new is added to it
		new MutationObserver(mutations => {
			mutations.forEach(mutation => {
				mutation.addedNodes.forEach(node => {
					if (node.nodeType === Node.ELEMENT_NODE && node.matches(`[data-alert]`)) {
						// If the new node is an alert, set a timeout to remove it after 5 seconds
						setTimeout(() => {
							node.remove()
						}, 5000)
					}
				})
			})
		}).observe(container, { childList: true })
	})

	// Global add alert function
	window.addAlert = function(containerID, message, type = "info") {
		// Create a new alert element
		const alert = document.createElement("div")
		alert.setAttribute("data-alert", "")
		alert.setAttribute("data-alert-type", type)
		alert.textContent = message

		// Add click handler to remove the alert on click
		alert.addEventListener("click", () => {
			alert.remove()
		})

		// Append the alert to the first alert container
		if (alertContainer.length > 0) {
			for (var container of alertContainer) {
				if (container.id === containerID) {
					container.appendChild(alert)
				}
			}
		} else {
			console.warn("No alert container found. Please add an element with [data-alert-container] attribute.")
		}
	}
})();
