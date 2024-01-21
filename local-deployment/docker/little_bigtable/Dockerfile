FROM golang:1.20

WORKDIR /app

RUN git clone https://github.com/gobitfly/little_bigtable.git /little_bigtable
RUN cd /little_bigtable && make install
