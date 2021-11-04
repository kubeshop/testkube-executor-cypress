# syntax=docker/dockerfile:1
FROM golang:1.17
WORKDIR /build
COPY . .
ENV CGO_ENABLED=0 
ENV GOOS=linux

RUN cd cmd/agent;go build -o /runner -mod mod -a .

FROM node:17-alpine3.12 
RUN apk --no-cache add ca-certificates
RUN npm install -g cypress
WORKDIR /root/
COPY --from=0 /runner /bin/runner
EXPOSE 8082
ENTRYPOINT ["/bin/runner"]
