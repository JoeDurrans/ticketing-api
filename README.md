# Setup

1. **Environment Setup**
    Copy the .env.example file and create a new file named .env. Update the environment variables in the .env file as per your local setup (this step can be skipped if you already have a .env file).

    ```sh
    cp .env.example .env
    ```

2. **Run the Docker Compose**
    Use Docker Compose to start the application dependencies.

    ```sh
    docker-compose up postgres scylla
    ```

3. **Initialise Databases**
    Create a database and keyspace for PostgreSQL and ScyllaDB respectively, using the database names set in the environments file (for simplicity id recommend creating the scylla keyspace without replication for local development).

    ```sql
    CREATE DATABASE your_database_name;
    ```

    ```sql
    CREATE KEYSPACE your_keyspace_name WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};
    ```

4. **Run Database Migrations** 
    Run the database migrations to set up the database schema:

    ```sh
    make postgres_migrate_up
    make scylla_migrate_up
    ```

5. **Begin the Application** 
    Again, use docker to start the application:
    
     ```sh
    docker-compose up go
    ```


## Notes

### Account Setup
You will need to add an user to utilise any of the routes. 
To do this, please modify the routes file from:

```go
	router.HandleFunc("POST /account", IsAdmin(makeHTTPHandleFunc(s.handleCreateAccount)))
```

To: 

```go
	router.HandleFunc("POST /account", makeHTTPHandleFunc(s.handleCreateAccount))
```

Then, re-run the application and create an account using the endpoint.

Finally, return the application to its original state and re-run it once more.

# Live

https://ticketing-api-production.up.railway.app

(Login details for the live environment are in the .env file which should be in the zip file)

**THE LIVE ENVIRONMENT MAY GIVE YOU AN ERROR INITIALLY. THIS IS BECAUSE THE DATABASE IS SLEEPING DUE TO INACTIVITY TO REDUCE COSTS. PLEASE SEND A REQUEST TO AN ENDPOINT THAT UTILISES THE DATABASE TO WAKE IT UP, WAIT, AND TRY AGAIN** 

**PLEASE DO NOT PERFORMANCE TEST THE LIVE APPLICATION FOR COST REASONS. THIS LIVE VERSION MAY BE DISABLED AT ANY TIME TO PREVENT UNNECESSARY COSTS SO YOU MAY HAVE TO SETUP THE LOCAL ENVIRONMENT ABOVE** 

# Testing

I've included a exported postman collection for all the REST endpoints so they are easy to test. Unfortunately, WebSocket endpoints cannot be exported so please use the demonstration video in the powerpoint to guide testing.

