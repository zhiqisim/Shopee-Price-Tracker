# wsa-entry-task

Entry task for new hires in the WSA team

## Deployment
Run the following batch script to deploy the docker containers
```
$ ./setup.sh
```
## Services
1. app - Nginx & React () 
    - Nginx load balances and reverse proxy to api-gateway
    - Access frontend on localhost
    - Access API on localhost/api
2. api-gateway - gateway to all API endpoints
    - api-gateway1 - server 1 (Golang/Gin)
    - api-gateway2 - server 2 (Golang/Gin)
4. user-service - service for user functionality (Golang)
5. item-service - service for item functionality
    - load balancer for item-service (HaProxy)
    - item-service1 - server 1 (Python)
    - item-service2 - server 2 (Python)
6. price-service - service to retrieve data from Shopee API (Python)
7. user-db - database for user data (MySQL)
8. items-db - database for item data (MySQL)
9. prom - prometheus gateway for metrics tracking
10. grafana - grafana dashboard for metrics tracking
11. cadvisor - monitoring for docker containers

## Logs
View logs at each service directories in the /log folder. 

## Business Requirement
    1. Our system needs to support register and login. Use could only access other API after login.
    2. Different users might be interested in different products. User will submit the product by its item detail page URL.
    3. By right our system should track all items in Shopee, however the volume is too huge to be handled in this entry task. To make it
        simple, we will only automatically put the flashsale items into our tracking system since they are more popular. Note: the items list will
        be different for different flashsale sessions.
    4. The system will support two kinds of queries
        a. List down all the items added by a buyer
        b. Display the price changelog of an item by itemid

## Technical Requirement
    1. You could either use firebase or mysql as our database.
    2. The system needs to be implemented using the microservice structure, which consists of:
        a. API gateway service in GoLang 1.12 + Gin 1.4.0
        b. User service in GoLang 1.12
        c. Price Tracking service in Python 2.7
        d. Pure Web Frontend in Django 1.6.11 or React. (No need to have a fancy UI. It is enough as long it could submit the form and
        display the result)
    3. The services will communicate with each other using protobuf via gRPC.
    4. Nginx is used to route traffic to gateway and frontend.
    5. Grafana + Prometheus need to be used to monitor the API(HTTP/RPC) traffic
    6. Redis is used as cache.
    7. The whole system should be deployed in single docker container or VirtualDev environment (VirtualDev / VM)
