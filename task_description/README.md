# Lynks Project

## Task

**Final Project: Short Hyperlink Storage Service** Welcome to the Lynks project, a service designed to store short hyperlinks efficiently. **Task Details:** The task description can be found on GitHub: [Final Task](https://github.com/thinknetica/go_course_4/blob/master/final.md).

### Overview

Lynks is a URL Shortener service that enables users to create short links associating them with long URLs. The service functionalities are summarized below:

### Service Capabilities

- **Generate Short Links:** Users can submit a long URL, and the service will generate a short URL based on its domain, store the URL pair in the database, and return the short URL.

- **Redirect to Original URL:** When a short URL is accessed, the service retrieves the associated long URL from the database and redirects the user to the original site.

### Example Usage

Imagine our service is running on the domain lynks.org. A user can create a short link using the following HTTP request:
bash

curl --location 'https://lynks.org' \
--header 'Content-Type: application/json' \
--data '{ "destination": "https://github.com/thinknetica/go_course_4" }'

Response from the service:
json

{ "shortUrl": "https://lynks.org/qz6d7", "destination": "https://github.com/thinknetica/go_course_4" }

Following the short link https://lynks.org/qz6d7 will redirect the user to https://github.com/thinknetica/go_course_4.
Component Composition and Architecture

The Lynks service consists of several interconnected components:

    Short Hyperlink Microservice: Provides API methods to generate and resolve short links.

    Short Link Microservice Database: A PostgreSQL database managing URL pairs, optimized with an index for rapid lookup by short links.

    Caching Microservice: Stores recently created URL pairs in memory for expedited access, with a Redis database backend.

    Caching Microservice Database: Utilizes Redis to quickly retrieve URL pairs.

Application Structure

Each microservice follows the course's directory structure guidelines. For more information, refer to the project directory structure. All microservices can be housed within a unified directory in the course module, negating the necessity for separate repositories.
Additional Requirements

    API Requirements: The service offers an HTTP API following REST principles. Use the GET method for retrieving data and POST for submissions. Communicate using JSON-formatted request and response bodies.

    Deployment Details: Deploy databases and microservices as containers. Provide a Dockerfile for each service and a final Docker Compose file to launch the application.

    Monitoring & Logging: Monitor API method requests and processing time using Prometheus metrics. Implement structured logging via zerolog, logging all API requests and errors to stdout.

    Caching Policy: Store hyperlink pairs in the caching database for a maximum of 24 hours, with specified expiration settings.
