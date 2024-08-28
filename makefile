ADDR="192.168.178.58"
PIN1="0"
PIN2="1"
PIN3="2"
PIN4="3"

build:
	go build -o build/ ./cmd/picow-led/

test-setup:
	go run ./cmd/picow-led -debug -addr ${ADDR} \
		set -id 1 -full-response -pretty-print config led ${PIN1} ${PIN2} ${PIN3} ${PIN4} \
		get -id 2 -full-response -pretty-print config led

test-on:
	go run ./cmd/picow-led -debug -addr ${ADDR} \
	    set -id 1 -full-response -pretty-print led duty 100 100 100 100
	go run ./cmd/picow-led -debug -addr ${ADDR} \
	    get -id 2 -full-response -pretty-print led duty

test-off:
	go run ./cmd/picow-led -debug -addr ${ADDR} \
	    set -id 1 -full-response -pretty-print led duty 0 0 0 0
	go run ./cmd/picow-led -debug -addr ${ADDR} \
	    get -id 2 -full-response -pretty-print led duty
