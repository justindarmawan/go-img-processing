FROM debian:bullseye

RUN apt-get update && apt-get install -y \
    git \
    cmake \
    make \
    pkg-config \
    libjpeg-dev \
    libpng-dev \
    libtiff-dev \
    libavcodec-dev \
    libavformat-dev \
    libswscale-dev \
    libv4l-dev \
    libgtk-3-dev \
    libxvidcore-dev \
    libx264-dev \
    libatlas-base-dev \
    gfortran \
    wget \
    unzip \
    sudo \
    && rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.21.1
RUN wget -qO- "https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz" | tar -C /usr/local -xzf -
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN git clone --branch v0.35.0 https://github.com/hybridgroup/gocv.git $GOPATH/src/gocv.io/x/gocv
WORKDIR $GOPATH/src/gocv.io/x/gocv
RUN make install

WORKDIR $GOPATH/src/app
COPY . .
RUN go mod download
RUN go mod tidy

RUN go test ./internal/test -v
RUN go build -o app
CMD ["./app", "rest"]
