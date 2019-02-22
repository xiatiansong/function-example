FROM dev.docker.pt.xiaomi.com/miot/golang:1.10

ENV HTTP_PROXY http://10.231.45.239:1080
ENV HTTPS_PROXY http://10.231.45.239:1080
RUN export http_proxy=$HTTP_PROXY
RUN export https_proxy=$HTTPS_PROXY

# Install dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN go get golang.org/x/net/context
RUN mkdir -p /go/pkg/dep && chown 1000:1000 /go/pkg/dep

# Prepare function environment
RUN ln -s /kubeless $GOPATH/src/kubeless

RUN mkdir /.cache && chown 1000:1000 /.cache

ADD . /go/src/controller/

# Install controller
USER 1000
