FROM python:2.7

WORKDIR /item-service


COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt

COPY . /item-service

CMD [ "python", "./service.py" ]