SRC = main.go
CC = go
BIN = main
GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux 
run:
	${CC} run ${SRC}
mock:
	${CC} run ${SRC}  -mock
build:
	${CC} build -o ${BIN} ${SRC}
static-build:
	${GO_BUILD_ENV} ${CC} build -a -installsuffix cgo -o ${BIN} .
test:
	${CC} test