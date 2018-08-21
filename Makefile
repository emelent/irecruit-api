SRC = main.go
CC = go
BIN = main
GO_PRODUCTION_ENV := CGO_ENABLED=0 GOOS=linux 
GO_BUILD_ENV := CGO_ENABLED=0 

# run in development env
run:export ENV := development
run:
	${CC} run ${SRC}
	
# run in development env with mock database
mock:export ENV := development
mock:
	${CC} run ${SRC}  -mock

# standard local build
build:
	${CC} build -o ${BIN} ${SRC}
	
# create static build for local sys
static-build:
	${GO_BUILD_ENV} ${CC} build -a -installsuffix cgo -o ${BIN} .

# create production build
production-build:
	${GO_PROD_ENV} ${CC} build -a -installsuffix cgo -o ${BIN} .

# run all tests
test:export ENV := test
test:
	${CC} test ./tests/...

# run functional tests
test_functional:export ENV := test
test_functional:
	${CC} test ./tests/functional

# run unit tests
test_unit:export ENV := test
test_unit:
	${CC} test ./tests/unit