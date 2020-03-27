FROM centos:7

# install prerequisites
RUN yum install -y git zip gcc && yum clean all

# Setup golang environment
ENV APPDIR="/opt"
ENV GOROOT="$APPDIR/go" \
    GOPATH="/gopath"
ENV PATH=$PATH:$GOROOT/bin

# Install golang engine
RUN curl -Ls https://dl.google.com/go/go1.14.1.linux-amd64.tar.gz | tar -xzC "$APPDIR/"

# Create a copy of the repository inside the container.
RUN mkdir -p "$GOPATH/src" && \
    cd $GOPATH
ENV GOPROJ="$GOPATH/src/github.com/iftsoft/device"
COPY ./. $GOPROJ/.

# Build all projects
WORKDIR $GOPROJ
RUN ./build.sh
RUN ./package.sh
