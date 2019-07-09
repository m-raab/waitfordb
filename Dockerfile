FROM golang:1.12.5 AS builder

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN apt-get update && apt-get install unzip libaio1 && chmod +x /usr/bin/dep

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/waitfordb
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./

RUN unzip /go/src/waitfordb/instantclient-basiclite-linux.x64-12.2.0.1.0.zip -d /opt/oracle && \
    ln -s /opt/oracle/instantclient_12_2/libclntsh.so.12.1 /opt/oracle/instantclient_12_2/libclntsh.so && \
    ln -s /opt/oracle/instantclient_12_2/libocci.so.12.1 /opt/oracle/instantclient_12_2/libocci.so

# RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix nocgo -o /waitfordb .

RUN CGO_ENABLED=1 GOOS=linux go build -a -o /waitfordb .

ENV LD_LIBRARY_PATH=/opt/oracle/instantclient_12_2:/lib:/lib64

RUN /waitfordb --jdbcurl=jdbc:oracle:thin:@jmraabmac.dhcp.j.ad.intershop.net:1521:XE --user=intershop --password=intershop && \
    exit 0