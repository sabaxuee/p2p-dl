# This is a local-build docker image for p2p-dl test
FROM ubuntu:14.04
MAINTAINER Zonesan <san@zoobook.org>
EXPOSE 9090
RUN echo "This is a local-build docker image for test" > /home/readme
CMD ["mkdir", "/home/p2p-dl"]
ADD p2p-dl create_testfile.sh /home/p2p-dl/
ENTRYPOINT  /home/p2p-dl/p2p-dl
