FROM golang:1.20

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
    && rm -rf /var/lib/apt/lists/*

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p $GOPATH/src/app

WORKDIR $GOPATH/src/app

COPY . .

RUN make install

RUN go mod download

RUN go build -o app

CMD ["./app", "rest"]
