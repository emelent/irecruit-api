# prepare build
FROM golang:1.10.3-stretch as build
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go get -d -v ./...
RUN make production-build

# prepare final image
# ADD create env files
# TODO export ENV=production
FROM scratch
COPY --from=build /app/main /main
CMD ["/main"]
EXPOSE 9999

