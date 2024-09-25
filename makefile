SERVER_ADDR=localhost:50833
TEST_SETUP_DEVICE='{ "server": { "name": "Kitchen", "addr": "192.168.178.58:3000" }, "pins": [0, 1, 2, 3] }'
TEST_SETUP_COLOR_WHITE='[255, 255, 255, 255]'

build:
	cd frontend && \
		npm install && \
		npm run build && \
		cd .. && \
		go build -v -o ./build/ ./cmd/picow-led-server/ && \
		cp ./cmd/picow-led-server/picow-led-server.service ./build/

run:
	go run -v ./cmd/picow-led-server -debug

setup:
	curl -X POST \
		--header 'content-type: application/json' \
		--data ${TEST_SETUP_DEVICE} \
		${SERVER_ADDR}/api/device | jq
	curl -X POST \
		--header 'content-type: application/json' \
		--data ${TEST_SETUP_COLOR_WHITE} \
		${SERVER_ADDR}/api/colors/white | jq
	curl -X GET \
		${SERVER_ADDR}/api | jq
