FROM golang:1.12.5 AS builder

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

COPY instantclient-basiclite-linux.x64-12.2.0.1.0.zip /tmp

# RUN mkdir /opt/oracle && unzip /tmp/instantclient-basiclite-linux.x64-12.2.0.1.0.zip -d /opt/oracle &&\
#    ln -s /opt/oracle/instantclient_12_2/libclntsh.so.12.1 /opt/oracle/instantclient_12_2/libclntsh.so && \
#    ln -s /opt/oracle/instantclient_12_2/libocci.so.12.1 /opt/oracle/instantclient_12_2/libocci.so


# Copy the code from the host and compile it
WORKDIR $GOPATH/src/waitfordb
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix nocgo -o /waitfordb .

CGO_ENABLED=1 GOOS=linux go build -o /waitfordb .