FROM python:2.7-alpine

WORKDIR /price-service

COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt

COPY . /price-service

CMD [ "python", "./scheduler.py" ]