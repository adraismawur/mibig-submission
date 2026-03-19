FROM golang:latest

WORKDIR /

RUN git clone https://github.com/adraismawur/mibig-submission.git

WORKDIR /mibig-submission/api

RUN go build mibig-submission.go

RUN chmod a+x ./mibig-submission

CMD ["./mibig-submission", "start-api"]
