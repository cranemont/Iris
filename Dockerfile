FROM ubuntu:20.04 AS base

COPY sources.list /etc/apt/
ENV TZ=Asia/Seoul
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone


#######################
## Install libjudger ##
#######################
FROM base AS build-judger

WORKDIR /build
COPY ./libjudger ./

RUN apt-get update && apt-get -y install libseccomp-dev cmake
RUN cmake CMakeLists.txt && make && make install


#####################
## for development ##
#####################
FROM base AS development

WORKDIR /app
RUN mkdir -p sandbox/policy \
  && mkdir sandbox/results\
  && mkdir -p sandbox/logs/run \
  && mkdir -p sandbox/logs/compile

COPY libjudger/java_policy sandbox/policy/

RUN buildDeps='software-properties-common curl' \
  && apt-get update && apt-get install -y $buildDeps \
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

COPY --from=build-judger /build/output/libjudger.so ./sandbox/

RUN apt-get update && apt-get install -y wget git vim
RUN wget -P /tmp https://go.dev/dl/go1.19.linux-amd64.tar.gz

RUN tar -C /usr/local -xzf /tmp/go1.19.linux-amd64.tar.gz
RUN rm /tmp/go1.19.linux-amd64.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src/github.com/cranemont" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

WORKDIR $GOPATH/src/github.com/cranemont
# RUN git clone https://github.com/cranemont/judge-manager
ENV APP_ENV=dev