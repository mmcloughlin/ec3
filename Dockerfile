FROM ubuntu:bionic

ARG GO_VERSION=1.13.1
ARG Z3_VERSION=4.8.6

# Install required packages.
RUN apt-get update
RUN apt-get install --yes \
        build-essential \
        ca-certificates \
        git \
        gzip \
        python \
        ssh \
        tar \
        wget \
    ;

# Install go.
WORKDIR /tmp/go
RUN wget --quiet \
        --output-document go.tar.gz \
        https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz \
    && \
    tar -C /usr/local -xzf go.tar.gz \
    ;

# Install z3.
WORKDIR /tmp/z3
RUN wget --quiet \
        --output-document z3.tar.gz \
        https://github.com/Z3Prover/z3/archive/z3-${Z3_VERSION}.tar.gz \
    && \
    tar -xzf z3.tar.gz --strip-components=1 \
    && \
    python scripts/mk_make.py --prefix /usr/local --staticlib \
    && \
    make -C build -j 8 install \
    ;

# Setup environment.
ENV GOPATH /go
ENV PATH ${GOPATH}/bin:/usr/local/go/bin:${PATH}

RUN mkdir -p "${GOPATH}/src" "${GOPATH}/bin" \
    && \
    chmod -R 777 "${GOPATH}" \
    ;
WORKDIR ${GOPATH}
