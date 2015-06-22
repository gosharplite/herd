FROM google/debian:wheezy
ADD herd herd
EXPOSE 8090
ENTRYPOINT ["/herd"]