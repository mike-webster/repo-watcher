FROM golang:latest
WORKDIR /repo-watcher
COPY . .
RUN go build -o repo-watcher

CMD ["/repo-watcher/repo-watcher"]