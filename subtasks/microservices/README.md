
# Execution:

docker-compose up --build


# Task1: Develop a Microservice

## Task #1: Develop a Microservice for Storing Short Hyperlinks

The microservice should provide an HTTP API with the following methods:

1. **Adding a link to the database** - The method accepts a link, adds it to the storage, and returns a short link.
2. **Getting the original link** - The method accepts a short link as input and returns the original.

### Implementation Details

- The link storage should be implemented as a slice in memory.
- Persistent storage is not required!

---

## Task #2: Run the Developed Microservice in a Container

1. **Write a Dockerfile.**
2. **Run the microservice on your computer** and connect to the application using a browser or another HTTP client (Postman, etc.).




# Task2

## Add a Microservice to the URL Shortening System

### Task #1

1. Add an analytics microservice to the system.
2. The new microservice should provide statistics on the number of links: total count, average length, etc.
3. When saving a new hyperlink, the existing microservice should send information about it to a queue (Kafka or RabbitMQ).
4. The analytics microservice should subscribe to this queue and update the statistics upon receiving a new message.

