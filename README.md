# Information about the project GYMONDO

Here is my solution to the task for Gymondo.

How to Start the Service in Docker?

To start the service and databases in Docker, follow these steps:

1. Clone the repository:
    ```
    git clone git@github.com:GermanLepin/gymondo.git
    ```

2. Start the service and database with a single command:
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

In this case, the service will be reachable on port 80 (http://localhost:80).

Note: When the service is launched, migrations are automatically performed, and the database is populated with initial data.

Open question: Should a user be allowed to have only one active (non-cancelled) subscription at a time? 
This rule was not explicitly mentioned in the task. In my current implementation, it is possible for a user to have multiple active subscriptions. However, I am open to modifying this if needed.


# API

The list of all active handles:
1. `/api/products/"`
   
   Retrieves a list of all available products. This endpoint provides information about the products that users can subscribe to.


2. `/api/products/:voucher_code`

   Fetches details of a specific product associated with a given voucher code. The voucher code is used to apply discounts or offers to the product.


3. `/api/product/:product_id`

   Retrieves detailed information about a specific product using the unique product_id. This includes pricing, description, and other attributes.


4. `/api/product/subscribe/`

   Allows users to subscribe to a product. This endpoint creates a new subscription for a user, including selecting a product and setting the subscription parameters (e.g., trial period, voucher code).


5. `/api/subscription/:subscription_id`
   Provides details of an active subscription. The subscription_id is used to fetch information about a specific subscription, such as its status, start date, end date, and other relevant information.


6. `/api/subscription/:subscription_id/manage`
   Manages an existing subscription. This endpoint allows users to update or modify their subscription, such as pausing, canceling, or changing other settings related to the subscription.

