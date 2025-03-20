# Information about the project GYMONDO

Here is my solution to the task for Gymondo.

How to Start the Service in Docker?

To start the service and databases in Docker, follow these steps:

1. Clone the repository:
    ```
    git clone git@github.com:GermanLepin/gymondo.git
    ```

2. Move to the cloned repository and run the command:
    ```
    cd gymondo
    ```

3. Start the service and database with a single command:
    ```
    make up_build
    ```
    This command will build the necessary Docker images (if not already built) and start the containers.

Accessing the services:
* The main service will be accessible at http://localhost:9000.
* PostgreSQL database will be accessible at localhost:5432.

Alternative Method (from IDE)
If you prefer to start the service directly from your IDE, you can use the following command:

1. Start the databases with the following command:
    ```
    run_postgres
    ```

2. Start the service
    ```
    go run cmd/main.go
    ```

In this case, the service will be reachable on port 80 http://localhost:80.

Note: When the service is launched, migrations are automatically performed, and the database is populated with initial data.

Open question: Should a user be allowed to have only one active (non-cancelled) subscription at a time? 
This rule was not explicitly mentioned in the task. In my current implementation, it is possible for a user to have multiple active subscriptions. However, I am open to modifying this if needed.


# SWAGGER API

The documentation for the service is generated using gin-swagger:
https://github.com/swaggo/gin-swagger

Before testing the service, please check the following file: `gymondo/db/postgres/migrations/0001_init.go`
In this file, I have inserted initial data into the main tables (users, products, vouchers). 
This information will be really helpful for testing the service.

The Swagger service starts automatically, and the documentation is available at:
http://localhost:9000/swagger/index.html

To regenerate the documentation after making changes, run the following command:
```
swag init -g ./cmd/main.go -o cmd/docs
```

I hope you enjoy using the service!
