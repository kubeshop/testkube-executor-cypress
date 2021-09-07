# syntax=docker/dockerfile:1
FROM golang:1.17
WORKDIR /build
COPY . .
ENV GONOSUMDB=github.com/kubeshop/* 
ENV CGO_ENABLED=0 
ENV GOOS=linux
RUN cd cmd/executor;go build -o /app -mod mod -a .


FROM cypress/included:8.3.0
# git for checking out partial git repositories (need to be recent one 2.33)
RUN echo "deb http://ftp.de.debian.org/debian sid main" >> /etc/apt/sources.list 
RUN apt-get update && apt-get install -y git
WORKDIR /root/
COPY --from=0 /app /bin/app
EXPOSE 8084
ENTRYPOINT ["/bin/app"]
