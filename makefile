clean:
	@git clean -f -x -d

build:
	@rm -rf dist && \
		go mod tidy -v && \
		cd ui && \
		npm install && \
		npm run build && \
		cd .. && \
		go build -v -o dist/picow-led-server ./cmd/picow-led-server

dev:
	#DEBUG=nodemon:*,nodemon nodemon -L --signal SIGTERM --exec 'go run ./cmd/picow-led-server -d -c .api.dev.json' --ext '' --delay 3
	nodemon -L --signal SIGTERM --exec 'go run ./cmd/picow-led-server -d -c .api.dev.json' --ext 'go,mod,sum' --delay 3 --ignore ./ui

