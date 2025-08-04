# Working with Microservices in Go
----
This project consists of a number of loosely coupled microservices, all written in Go:
- broker-service: an optional single entry point to connect to all services from one place (accepts JSON; sends JSON, makes calls via gRPC, and pushes to RabbitMQ)
- authentication-service: authenticates users against a Postgres database (accepts JSON)
- logger-service: logs important events to a MongoDB database (accepts RPC, gRPC, and JSON)
- queue-listener-service: consumes messages from amqp (RabbitMQ) and initiates actions based on payload (sends via RPC)
- mail-service: sends email (accepts JSON)

In addition to the microservices, the included docker-compose.yml at the project directory of the project starts the following services:

Postgresql - used by the authentication service to store user accounts
MongoDB - used by the logger service to save logs from all services
mailhog - used as a fake mail server to work with the mail service

## Running the Project
----
From the Project Directory of the project, execute this command (this assumes that you have GNU make and a recent version of Docker installed on your machine):
`make up_build`
If the code has not changed, subsequent runs can just be `make up`

Then Start the front end:
`make start`
Hit the front end with your web browser at http://localhost:80

To stop Everything:
`
make stop
makw down
`
