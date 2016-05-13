# This is a local-build docker image for p2p-dl test
FROM alpine
MAINTAINER Zonesan <san@zoobook.org>
EXPOSE 9090
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && \
    echo "This is a local-build docker image for test" > /home/readme && \
    mkdir /home/p2p-dl
ADD p2p-dl create_testfile.sh /home/p2p-dl/
ENTRYPOINT  /home/p2p-dl/p2p-dl
