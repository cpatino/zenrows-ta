# ZenRows Technical Assessment

The repository contains the GO application that is part of the technical assessment for ZenRows.

## Libraries used in the application

- [Gin] The most popular web framework for Go, used to build the RESTful API endpoints of the application.
- [MongoDB Go Driver] The official MongoDB driver for Go, used to interact with the MongoDB database for storing and
  retrieving data.
- [Logrus] A structured logger for Go, used for logging application events and errors in a structured format.
- [Godotenv] A library for loading environment variables from a .env file, used to manage configuration settings for the
  application.

## Why use mongodb as the database for this application?

One of the constraints set on tech assessment was that the database schema should support future changes gracefully.
MongoDB, being a document-oriented/NoSQL database, allows for flexible schema design, which means that we can easily add
or modify fields in our documents without having to worry about schema migrations or downtime.

## How to run the application:

### Pre-requisites:

- Docker installed on your machine.
- Run ``docker compose up -d`` from misc folder to start the mongodb container and create the network for the
  application.
- Wait until the mongodb container is up and running before starting the application container.

### Already prepared data:

There is a file named seed.js in the misc/mongo-seed folder that contains already baked data for the application. Also,
it creates the indexes for the collection.
When starting the mongodb container, the data is automatically seeded into the database using a docker volume that
mounts the seed.js file to the mongodb container and runs it as part of the container initialization process.

#### Basic auth users:

| username | password |
|----------|----------|
| test1    | test1    |
| test2    | test2    |

### Run locally

1. Clone the repository and navigate to the project directory.
2. Run the application using the command:
   ```go run .```

### Run in a container

1. Build the image:
   ```docker build -t zenrows-ta .```

2. This creates an image named zenrows-ta from the Dockerfile in the current directory.
   Run the container:
   ```docker run -p 8080:8080 --network misc_zenrows-ta-network zenrows-ta```

## How to manually test the application

Use the postman collection provided in the misc folder to test the API endpoints. The collection contains pre-configured
requests for all the API endpoints, along with the necessary authentication details.

**Important**: Create an environment in Postman to allow the pre-build scripts to set the variables correctly.

## How to run the tests:

1. Run the tests using the command:
   ```go test -v ./...```

## Enhancements and future improvements:
- Implement pagination for the find all endpoint to handle large datasets efficiently.
- Add more comprehensive error handling and logging to improve the application's robustness and maintainability.
- Better authentication and authorization mechanism, to enhance security.
- Edit endpoint overwrites the entire document, the application must expose a patch endpoint that allows partial updates to the document.
- Create allow to duplicate entries (same data with different ids).
- Better testing, I did not put the effort to create a better testing suite due to time constraints.
