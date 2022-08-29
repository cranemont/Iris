FROM ubuntu:20.04

COPY sandbox /tmp
COPY sources.list /etc/apt/

ENV TZ=Asia/Seoul
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN buildDeps='software-properties-common git libtool cmake python-dev python3-pip libseccomp-dev curl' && \
  apt-get update && apt-get install -y build-essential python python3 python-pkg-resources python3-pkg-resources $buildDeps && \
  add-apt-repository ppa:ubuntu-toolchain-r/test && \
  add-apt-repository ppa:openjdk-r/ppa && \
  apt-get update && apt-get install -y openjdk-11-jdk gcc-9 g++-9 && \
  update-alternatives --install  /usr/bin/gcc gcc /usr/bin/gcc-9 40 && \
  update-alternatives --install  /usr/bin/g++ g++ /usr/bin/g++-9 40 && \
  cd /tmp && cmake CMakeLists.txt && make && make install && \
  apt-get purge -y --auto-remove $buildDeps && \
  apt-get clean && rm -rf /var/lib/apt/lists/* && \
  useradd -u 12001 compiler && useradd -u 12002 code && useradd -u 12003 spj && usermod -a -G code spj

RUN apt-get update && apt-get install -y wget git vim
RUN wget -P /tmp https://go.dev/dl/go1.19.linux-amd64.tar.gz

RUN tar -C /usr/local -xzf /tmp/go1.19.linux-amd64.tar.gz
RUN rm /tmp/go1.19.linux-amd64.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src/github.com/cranemont" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

WORKDIR $GOPATH/src/github.com/cranemont
RUN git clone https://github.com/cranemont/judge-manager