SRC = main.go
CC = go

run:
	${CC} run ${SRC}

build:
	${CC} build -o server ${SRC}

test:
	${CC} test