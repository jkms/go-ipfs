FROM golang:1.9.2-alpine
MAINTAINER John Stafford <john@jkms.me>

# There is a copy of this Dockerfile called Dockerfile.fast,
# which is optimized for build time, instead of image size.
#
# Please keep these two Dockerfiles in sync.
ENV GX_IPFS ""
ENV SRC_DIR /go/src/github.com/ipfs/go-ipfs

COPY . $SRC_DIR

# Get su-exec, a very minimal tool for dropping privileges,
# and tini, a very minimal init daemon for containers
RUN apk update
RUN apk add build-base git binutils binutils-gold libc6-compat
RUN apk add su-exec
RUN apk add tini

# Get the TLS CA certificates, they're not provided by busybox.
RUN apk add ca-certificates openssl

# Build the thing.
# Also: fix getting HEAD commit hash via git rev-parse.
# Also: allow using a custom IPFS API endpoint.
RUN cd $SRC_DIR \
  && mkdir .git/objects \
  && ([ -z "$GX_IPFS" ] || echo $GX_IPFS > /root/.ipfs/api) \
  && make build

# Get the ipfs binary, and entrypoint script
RUN cp $SRC_DIR/cmd/ipfs/ipfs /usr/local/bin/ipfs
RUN cp $SRC_DIR/bin/container_daemon /usr/local/bin/start_ipfs

# This shared lib (part of glibc) doesn't seem to be included with busybox.
# COPY --from=0 /lib/x86_64-linux-gnu/libdl-2.24.so /lib/libdl.so.2

# Ports for Swarm TCP, Swarm uTP, API, Gateway, Swarm Websockets
EXPOSE 4001
EXPOSE 4002/udp
EXPOSE 5001
EXPOSE 8080
EXPOSE 8081

# Create the fs-repo directory and switch to a non-privileged user.
ENV IPFS_PATH /data/ipfs
RUN mkdir -p $IPFS_PATH \
  && adduser -D -h $IPFS_PATH -u 1000 -G users ipfs \
  && chown ipfs:users $IPFS_PATH

# Expose the fs-repo as a volume.
# start_ipfs initializes an fs-repo if none is mounted.
# Important this happens after the USER directive so permission are correct.
VOLUME $IPFS_PATH

# The default logging level
ENV IPFS_LOGGING ""

# This just makes sure that:
# 1. There's an fs-repo, and initializes one if there isn't.
# 2. The API and Gateway are accessible from outside the container.
ENTRYPOINT ["/sbin/tini", "--", "/usr/local/bin/start_ipfs"]

# Execute the daemon subcommand by default
CMD ["daemon", "--migrate=true"]
