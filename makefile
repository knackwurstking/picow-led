build:
	@go mod tidy && \
		cd frontend && \
		npm install && \
		npm run build && \
		cd .. && \
		go build -o build/picow-led-server ./cmd/picow-led-server
