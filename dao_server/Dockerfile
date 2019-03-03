FROM golang
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/hashicorp/raft
RUN go get -u github.com/otoolep/hraftd/store
RUN go get -u github.com/lib/pq
RUN go get -u github.com/hashicorp/raft-boltdb
COPY . .
EXPOSE 8888
EXPOSE 12000
CMD go run *.go