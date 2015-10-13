# PizzAPI

This API provides a way of viewing pizzas and ordering them.

| Verb | Path        | Description                          |
|------|-------------|--------------------------------------|
| GET  | /pizzas     | List all pizzas which can be ordered |
| POST | /orders     | Order a new pizza                    |
| GET  | /orders/:id | View an order pizza's status         |

There is no database. Everything is stored in memory, except for the pizzas list which is in `pizza.json`.

### Ordering a pizza

> curl --data '{"id":1}' http://localhost:5000/orders

## Constraints

I'm using this API to teach resiliance in a master's degree. While it would work nicely by default, we can add constraints:

* `maintenance` - All requests will respond with a 503.
* `slow` - All requests will wait between 30 and 60 seconds before being answered.
* `erroring` - 7 requests out of 10 will respond with a 500. All others will answer properly.

Constraints are scoped by IP address, so they can be added only for some students.

## Adding a constraint

> curl --data '{"name":"maintenance","ip":"[::1]:51192"}' http://localhost:5000/upgrade

## Contributing

I'm open to contributions on this project.  
However, please note that I'm using it for the lessons I'm giving. New constraints will not necessarily be adopted (they can always be suggested though).
