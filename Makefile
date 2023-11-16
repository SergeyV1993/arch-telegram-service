APP?=telegram

clean:
	rm -f bin/${APP}

init:
	go get all

run: clean init
	go build -o bin/${APP} telegram.go && bin/${APP}

compose:
	docker build -t telegram .

docker-run: compose
	docker run --env-file .env -d telegram

docker-compose-run:
	docker-compose up -d