# wsa-entry-task

Entry task for new hires in the WSA team

## Deployment
Run the following batch script to deploy the docker containers
```
$ ./setup.sh
```

## Documentation
View the documentations in the /Documents directory
1. Design Document 
2. API Document
3. Performance Test Report
4. Concluding Report

## Services
1. app
    - Use as a router to route request to frontend and api-gateway
    - Built with Nginx
    - Frontend built with React
    - Frontend accessed with http://localhost
    - Api-gateway accessed with http://localhost/api
    - Use as a load balancer for api-gateway to balance request between 2 api-gateway servers
2. Api-gateway
    - Use as a gateway for all API endpoints
    - Built in Golang & Gin
    - 2 instance of api-gateway running for load to be balanced between
    - Passes each request to their respective services
    - Connected to respective services via gRPC acting as a gRPC client
3. Redis-user-sessions
    - Use as a cache for api-gateway
    - Store user sessions 
4. User-service
    - Service for user functionalities
    - Built with Golang
    - Connected to user-db via TCP
    - Connected to api-gateway via gRPC acting as a gRPC server
5. Item-service
    - Service for item functionalities
    - Built with Python
    - Connected to items-db via TCP
    - There is a load balancer using HaProxy to balance load between 2 instances of item-service
    - Connected to api-gateway via gRPC acting as a gRPC server
6. Price-service
    - Service to retrieve data from Shopee API
    - Built with Python
    - Connected to items-db via TCP
7. User-db
    - Database for user data
    - Using MySQL as DBMS
    - User table to store information of all user login credentials 
    - UserItems table to store information of all user’s watchlist items
8. Items-db
    - Database for item data
    - Using MySQL as DBMS
    - Item table to store information of item details
    - ItemPrice table to store price changelog of items
9. Prom
    - Prometheus gateway for metrics tracking
10. Grafana
    - Grafana dashboard for metrics tracking
11. Cadvisor
    - Monitoring for docker containers


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

## Perfermance Requirement
    1. Do stress test and try to support as much as QPS as you can.
    2. You need to find out the bottleneck and try to solve it until nothing could be improved.
    3. Metrics may include requests count, response latency, response status code, response size and etc for different endpoint and instances. CPU, memory usage, network usage, IO usage of the process and system during testing.
