FROM python:2.7

ENV TZ Asia/Singapore
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /item-service

# Copy wait-for-it.sh into our image
COPY wait-for ./item-service/wait-for
# Make it executable, in Linux
RUN chmod +x ./item-service/wait-for


COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt

COPY . /item-service

CMD ["python", "./service.py" ]
