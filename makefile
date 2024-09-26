SERVER_ADDR=localhost:50833
TEST_SETUP_COLOR_WHITE='[255, 255, 255, 255]'

build:
	cd frontend && \
		npm install && \
		npm run build && \
		cd .. && \
		go build -v -o ./build/ ./cmd/picow-led-server/ && \
		cp ./cmd/picow-led-server/picow-led-server.service ./build/

run:
	cd frontend && npm run build
	go run -v ./cmd/picow-led-server -debug

setup:
	# Add "Kitchen" device
	curl -X POST \
		--header 'content-type: application/json' \
		--data '{ "server": { "name": "Kitchen", "addr": "192.168.178.58:3000" }, "pins": [0, 1, 2, 3] }' \
		${SERVER_ADDR}/api/device
	# Add "PS Room" device
	curl -X POST \
		--header 'content-type: application/json' \
		--data '{ "server": { "name": "PC Room", "addr": "192.168.178.68:3000" }, "pins": [0, 1, 2, 3] }' \
		${SERVER_ADDR}/api/device
	# Add "Sleep Room" device
	curl -X POST \
		--header 'content-type: application/json' \
		--data '{ "server": { "name": "Sleep Room", "addr": "192.168.178.67:3000" }, "pins": [0, 1, 2, 3] }' \
		${SERVER_ADDR}/api/device
	# Set "white" color
	curl -X POST \
		--header 'content-type: application/json' \
		--data ${TEST_SETUP_COLOR_WHITE} \
		${SERVER_ADDR}/api/colors/white
	# Get ap data
	curl -X GET \
		${SERVER_ADDR}/api
