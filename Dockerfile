FROM golang:latest
WORKDIR /repo-watcher
COPY . .

ENV GIN_MODE=release
ENV GO_ENV=production

RUN  go build -o repo-watcher

EXPOSE 3100
CMD ["/repo-watcher/repo-watcher"]