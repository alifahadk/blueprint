# syntax=docker/dockerfile:1

#####
# Auto-generated Dockerfile for process workspace reserv_ctr
#   Dockerfile auto-generated by linuxcontainer plugin
#   Code generation located at linuxcontainer/linuxgen/dockerfilegen.go
#

###
# Step 1: run custom commands provided by processes
###



####### BEGIN
#  custom docker build commands provided by goproc.Process reserv_proc
#

FROM golang:1.23-bookworm AS reserv_proc

COPY ./reserv_proc /src

WORKDIR /src
RUN go mod download

RUN mkdir /reserv_proc
RUN go build -o /reserv_proc ./reserv_proc

#
# custom docker build commands provided by goproc.Process reserv_proc
######## END




###
# Step 2: prepare the final image
###

FROM gcr.io/distroless/base-debian12

# Copy artifacts for processes that didn't have custom build commands


# Copy artifacts for processes with custom build commands
COPY --from=reserv_proc /reserv_proc /reserv_proc


# Get a shell
COPY --from=busybox:1.35.0-uclibc /bin/sh /bin/sh

# Copy the build.sh file and run it
WORKDIR /
COPY ./build.sh /
RUN ["/bin/sh", "./build.sh"]

# Copy the run.sh file and configure the entrypoint
WORKDIR /
COPY ./run.sh /
ENTRYPOINT ["bin/sh", "./run.sh"]
