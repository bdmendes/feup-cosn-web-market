# HaaS Web Market

<!-- SET: @architectMaster = @bdmendes -->
<!-- SET: @bloatLover = @sirze -->
<!-- SET: @spaghettiLover = @fernandorego -->

## Architectural system design <!-- ARCHITECTURE COPIED FROM PRELIMINARY REPORT  -->

<!-- @architectMaster - NEED REVIEW AND POSSIBLE UPDATES -->

### Domain model

Our subsystem is responsible for orders and everything related to them. This includes the products that are ordered, the payment methods used to pay for the order, the delivery methods used to deliver the order, and the consumers who place the orders.

![Domain model](./assets/domain.png)

While we are not responsible for *Supplier*, *Product* and *Category*, these are a direct requirement for the functioning of the subsystem and are therefore included here. This means that while other subsystems may model them differently, there should be an interface between the subsystems that consistently serializes them.

### Services architecture

We are interested in the interactions between the services that we are responsible for, and the interactions between our services and the services of other subsystems (in this diagram, the connection to the *Stock Service*).

![Services architecture](./assets/arch.png)

The (direct) connections may be intercepted, in the final system, by a gateway that will be responsible for routing the requests to the correct service and/or taking care of authentication, authorization and observability. This will be skipped in the early development stages and will be added later when the group responsible for the gateway delivers a working prototype that we can use.

## Services description and their operations

All services have been developed using Go, and each service is equipped with its own MongoDB database. Kafka is also used to enable communication through its subscribe/publish model, facilitating efficient exchange between different components of a system.

### Consumers

<!-- @architectMaster + @spaghettiLover(kafka) -->

Outlined below are key operations managed by this microservice:

### Orders

The Order Service, a pivotal microservice, is responsible for overseeing the comprehensive order processing workflow within the system. Beyond its primary function of receiving new orders, this service plays a fundamental role in the communication with the Payment Service for order validation and with the Delivery Service to facilitate the shipment of purchases. 
To automate order processing, interactions with the payment and delivery services are triggered upon order creation and successful payment processing, respectively.
Finally, upon successful payment processing, is published a message in the broker to notify other services about the ordered products. This is particularly useful to update product stock.

To facilitate the management of each order, there are 5 different states to control the process flow:

- PENDING
  - Represents the initial state when an order is created but payment processing has not occurred.
- AUTHORIZED 
  - Indicates that the payment has been successfully processed, but the order has not yet been shipped.
- SHIPPED    
  - Marks the state when the delivery service has accepted the order for shipment.
- DELIVERED
  - Confirm that the order has been successfully delivered to the customer.
- CANCELLED  
  - Allows users to cancel the order, but only if the order is PENDING.

![Order State Diagram](./assets/order_state_diagram.png)

Outlined below are key operations managed by this microservice:

- Create a new order.
- Update an existing order.
- Publish a message in the broker to alert other services about the ordered products.
- Retrieve a specific order by its unique identifier.
- Retrieve a list of all orders associated with a specific client.
- Retrieve a list of all orders.
- Communication between services to ensure the correct processing of orders.

### Delivery

<!-- @architectMaster -->

### Payments

<!-- @spaghettiLover -->

## Resilience Patterns
- specification and implementation (minimum of 2)

<!-- @bloatLover or @architectMaster -->

## Observability patterns
- specification and implementation (minimum of 2)

<!-- @bloatLover or @architectMaster -->

## Security implementation

<!-- @bloatLover or @spaghettiLover -->

<!-- SUPER SPAGHETTI CODE LEADS TO HIGH SECURITY DUE TO OBVIOUS REASONS -->

## Link to Service APIs specification in OpenAPI

<!-- @bloatLover -->
