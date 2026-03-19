FROM mambaorg/micromamba:latest

RUN micromamba install -c conda-forge -c bioconda -y -n base antismash go git gcc

ARG MAMBA_DOCKERFILE_ACTIVATE=1

RUN download-antismash-databases

WORKDIR /home/mambauser

RUN git clone https://github.com/adraismawur/mibig-submission.git

WORKDIR /home/mambauser/mibig-submission/api

RUN go build mibig-submission.go

RUN chmod a+x ./mibig-submission

CMD ["./mibig-submission", "start-api"]
