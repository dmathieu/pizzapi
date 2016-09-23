# PizzAPI

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

This API provides a way of viewing pizzas and ordering them.

| Verb | Path        | Description                          |
|------|-------------|--------------------------------------|
| GET  | /doc        | Get the API documentation            |
| GET  | /pizzas     | List all pizzas which can be ordered |
| POST | /orders     | Order a new pizza                    |
| GET  | /orders     | List all orders                      |
| GET  | /orders/:id | View an order pizza's status         |

There is no database. Everything is stored in memory, except for the pizzas list which is in `pizza.json`.

### API Documentation

We're using Swagger/OpenAPI v2.0 for the API documentation. The schema is located in `doc/schema.yaml`.

### Ordering a pizza

> curl -H "Authorization: team_name" --data '{"id":1}' http://localhost:5000/orders

## Constraints

I'm using this API to teach resiliance in a master's degree. While it would work nicely by default, we can add constraints:

* `maintenance` - All requests will respond with a 503.
* `slow` - All requests will wait between 30 and 60 seconds before being answered.
* `cancel` - 7 requests out of 10 will be closed without any response.

Constraints are scoped by the `Authorization` header, so they can be added only for some students.

## Adding a constraint

> curl --data '{"name":"maintenance","token":"team_name"}' http://localhost:5000/upgrade

A constraint can be added by setting an empty token.

### Upgrade authentication

By setting the `UPGRADE_KEY` environment variable, you can limit who can perform upgrades.

> curl -H "Authorization: my_token" --data '{"name":"maintenance","token":"team_name"}' http://localhost:5000/upgrade

## Contributing

I'm open to contributions on this project.  
However, please note that I'm using it for the lessons I'm giving. New constraints will not necessarily be adopted (they can always be suggested though).
