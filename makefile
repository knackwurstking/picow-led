clean:
	git clean -f -x -d

build:
	@go mod -v tidy && \
		cd frontend && \
		npm install && \
		npm run build && \
		cd .. && \
		go build -v -o build/picow-led-server ./cmd/picow-led-server
