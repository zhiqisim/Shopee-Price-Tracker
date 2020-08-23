# Shopee Price Tracker

Price tracker on Shopee's flash sale items

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
    - Access with localhost:9090
10. Grafana
    - Grafana dashboard for metrics tracking
    - Access with localhost:3000
11. Cadvisor
    - Monitoring for docker containers
    - Access with localhost:8080


## Logs
View logs at each service directories in the /log folder. 
