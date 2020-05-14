FROM python:2.7

WORKDIR /price-service

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

CMD [ "python", "./price_service.py" ]