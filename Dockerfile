FROM registry.access.redhat.com/ubi9/go-toolset AS builder

COPY . /src
WORKDIR /src

RUN \
    make /tmp/image-clone-controller OUTPUT_DIR=/tmp

FROM registry.access.redhat.com/ubi9/ubi-minimal

COPY --from=builder /tmp/image-clone-controller /usr/bin/image-clone-controller

CMD /usr/bin/image-clone-controller
