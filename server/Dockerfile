FROM golang
RUN go get -u github.com/gorilla/mux
COPY . .
EXPOSE 8888
CMD go run *.go
