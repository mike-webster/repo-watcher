FROM golang:latest
WORKDIR /repo-watcher
COPY . .
RUN GO_ENV=production go build -o repo-watcher

EXPOSE 3100
CMD ["/repo-watcher/repo-watcher"]