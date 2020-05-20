FROM python:2.7

ENV TZ Asia/Singapore
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /item-service


COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt

COPY . /item-service

CMD [ "python", "./service.py" ]