(function() {
	document.querySelectorAll(`[data-alert-container]`).forEach(container => {
		// Watch this container if anything new is added to it
		new MutationObserver(mutations => {
			mutations.forEach(mutation => {
				mutation.addedNodes.forEach(node => {
					if (node.nodeType === Node.ELEMENT_NODE && node.matches(`[data-alert]`)) {
						// If the new node is an alert, set a timeout to remove it on timeout
						setTimeout(() => {
							node.remove()
						}, 60000) // 60000 milliseconds = 60 seconds
					}
				})
			})
		}).observe(container, { childList: true })
	})
})();
