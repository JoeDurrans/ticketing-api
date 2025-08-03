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
