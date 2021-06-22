FROM debian:buster-slim
# Official Debian and Ubuntu images automatically run apt-get clean, so explicit invocation is not required.
RUN set -xe && apt-get update && apt-get install -y curl
COPY build/nxpx /usr/local/bin/nxpx
RUN chmod +x /usr/local/bin/nxpx
CMD ["/usr/local/bin/nxpx"]