#Events

## Event received
RouteCreated
- id
- distance
- directions
- - lat
- - lng

### Execute and return other event

FreightCalculated
- route_id
- amount

---
DeliveryStarted
- route_id

### Execute driver movimentation

DriverMoved
- route_id
- lat
- lng