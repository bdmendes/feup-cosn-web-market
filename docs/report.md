# HaaS Web Market

<!-- SET: @architectMaster = @bdmendes -->
<!-- SET: @bloatLover = @sirze -->
<!-- SET: @spaghettiLover = @fernandorego -->

## Architectural system design <!-- ARCHITECTURE COPIED FROM PRELIMINARY REPORT  -->

<!-- @architectMaster - NEED REVIEW AND POSSIBLE UPDATES -->

### Domain model

Our subsystem is responsible for orders and everything related to them. This includes the products that are ordered, the payment methods used to pay for the order, the delivery methods used to deliver the order, and the consumers that place the orders.

![Domain model](./assets/domain.png)

While we are not responsible for *Supplier*, *Product* and *Category*, these are a direct requirement for the functioning of the subsystem and are therefore included here. This means that while other subsystems may model them differently, there should be an interface between the subsystems that serializes them in a consistent way.

### Services architecture

We are interested in the interactions between the services that we are responsible for, and the interactions between our services and the services of other subsystems (in this diagram, the connection to the *Stock Service*).

![Services architecture](./assets/arch.png)

The (direct) connections may be intercepted, in the final system, by a gateway that will be responsible for routing the requests to the correct service and/or take care of authentication, authorization and observability. This will be skipped in early development stages and will be added later when the group responsible for the gateway delivers a working prototype that we are able to use.

## Services description and their operations

### Consumers

<!-- @architectMaster + @spaghettiLover(kafka) -->

### Orders

<!-- @spaghettiLover -->

The Order Service, a pivotal microservice, is responsible for overseeing the comprehensive order processing workflow within the system. Beyond its primary function of receiving new orders, this service plays a fundamental role in the communication with the Payment Service for order validation and with the Delivery Service to facilitate the shipment of purchases.

Outlined below are key operations managed by this microservice:


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
