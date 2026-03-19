FROM python:3.11

RUN git clone https://github.com/adraismawur/mibig-submission.git

WORKDIR /mibig-submission/web

RUN pip install --upgrade pip
RUN pip install .

CMD ["gunicorn", "--bind", "0.0.0.0:8991", "submission:create_app()"]
