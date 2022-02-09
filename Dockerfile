FROM golang:latest
RUN git clone https://github.com/wf001/rouph
WORKDIR /go/rouph
RUN make
ENV PATH=$PATH:/go/rouph/bin
