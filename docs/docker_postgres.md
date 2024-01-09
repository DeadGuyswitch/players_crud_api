# Using PostgreSQL with Docker

## Installing PostgreSQL as a Docker Container 

### Method 1: Without Using `docker-compose`

**Step 1:** Pull the latest official PostgreSQL image from Docker Hub using the following command:

```shell
docker pull postgres

```

**Step 2:** Run a new Docker container with the following command:
```shell
docker run --name players_crud_postgres -e POSTGRES_PASSWORD=******* -d postgres
```

**Step 3:** Verify if the container is running with the following command:
```shell
docker ps
```

## Creating a Database inside the PostgreSQL Container


**Step 1**: Connect to the PostgreSQL container using the following command:
  ```shell
  docker exec -it players_crud_postgres psql -U postgres 
  ```

**Step 2:** Create a database in psql using the following command:
```
createdb players
```