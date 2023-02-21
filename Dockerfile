### SERVER ###
FROM golang:1.20-alpine as build

WORKDIR /build
COPY . .

RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main ./main.go

### PRODUCTION ###
FROM ubuntu:20.04

COPY sources.list /etc/apt/
ENV TZ=Asia/Seoul
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app
RUN buildDeps='software-properties-common curl' \
  && apt-get update && apt-get install -y $buildDeps netcat \
  && add-apt-repository ppa:deadsnakes/ppa \
  && curl -sL https://deb.nodesource.com/setup_16.x | bash -E - \
  && apt-get update && apt-get install -y \
  gcc \
  g++ \
  nodejs \
  python3.10 \
  pypy3 \
  openjdk-17-jdk \
  && apt-get purge -y --auto-remove $buildDeps \
  && apt-get clean && rm -rf /var/lib/apt/lists/*


RUN mkdir -p sandbox/policy \
  && mkdir sandbox/results\
  && mkdir -p sandbox/logs/run \
  && mkdir -p sandbox/logs/compile

COPY lib/judger/policy/java_policy sandbox/policy/
COPY lib/judger/libjudger.so ./sandbox/
COPY --from=build /build/main .

ENV APP_ENV=production
COPY ./scripts/entrypoint.sh /app/
ENTRYPOINT [ "sh", "./entrypoint.sh" ]
