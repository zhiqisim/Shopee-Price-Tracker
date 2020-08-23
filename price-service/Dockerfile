FROM python:2.7-alpine


RUN apk add --no-cache tzdata
ENV TZ Asia/Singapore
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /price-service

COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt

COPY . /price-service

CMD [ "python", "./scheduler.py" ]