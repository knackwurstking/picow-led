clean:
	git clean -f -x -d

build:
	@go mod tidy -v && \
		cd frontend && \
		npm install && \
		npm run build && \
		cd .. && \
		go build -v -o build/picow-led-server ./cmd/picow-led-server
