# syntax=docker/dockerfile:1
FROM golang:1.16
WORKDIR /build
COPY . .
ENV GONOSUMDB=github.com/kubeshop/* 
ENV CGO_ENABLED=0 
ENV GOOS=linux
RUN cd cmd/executor;go build -o /app -mod mod -a .


FROM cypress/included:8.3.0
WORKDIR /root/
COPY --from=0 /app /bin/app
EXPOSE 8083
ENTRYPOINT ["/bin/app"]
