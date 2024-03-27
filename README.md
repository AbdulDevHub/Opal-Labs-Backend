# backend

## Setup

### .env

`.env.example` contains a template for your own setup. Make a copy of it with `cp .env.example .env`. 

### CompileDaemon installation

Follow the instalation below to enable Hot-reload for development.

<https://pkg.go.dev/github.com/githubnemo/compiledaemon>

### Postgres Setup

#### Local Installation

 - Download and setup PostgresSQL (<https://www.enterprisedb.com/downloads/postgres-postgresql-downloads>)
 - When setting up the account make password=password and go default for everything else.
 - Then search for _pgadmin4_ from windows search bar and open _pgadmin4_ and create a new database called `Website`

#### ElephantSQL

If you don't want to install and run a database locally, use ElephantSQL.

 - Create an account <https://customer.elephantsql.com/login>
 - Press `Create New Instance`
 - In the setup, ensure that Plan is set to _Tiny Turtle (Free)_

After setting up a Database in ElephantSQL, in your `.env`, change the following fields so that they _match_ your ElephantSQL Instance:

```bash
DB_USER="<User & Default Database>"
DB_PASSWORD="<Password>"
DB_NAME="<User & Default Database>"
DB_HOST="<Server>"
DB_PORT=5432
DB_SSL_MODE="disable"
```
### Redis Setup

#### Local Installation

 - Download and setup the latest stable release of Redis (<https://redis.io/download>)
 - Ensure that the server is running on `localhost:6379` (this is the default and has been set in the `.env.example` template file)
 - Default password is `""` (empty string), and the default database is `0`. These are also set in the `.env.example` file.
 - If you have a different setup, change the `REDIS_ADDR`, `REDIS_PASSWORD` and `REDIS_DB` fields in your `.env` file to match your setup.

#### Hosted Redis

 - If you don't want to install and run a Redis server locally, you can use a hosted service such as RedisLabs.
 - Create an account at <https://redislabs.com/>, and create a new _free_ database.
 - In the `.env` file, change the following fields so that they _match_ your RedisLabs setup:
 ```
 REDIS_ADDR="<Public Endpoint>"
 REDIS_PASSWORD="<Default User Password>"
 REDIS_DB=0
 ```
 - The information can be found in the _Configuration_ tab under your specific database in the RedisLabs website.
 - If you have a different setup, change the `REDIS_ADDR`, `REDIS_PASSWORD` and `REDIS_DB` fields in your `.env` file to match your setup.

#### Redis Usage
- The backend will run if the Redis server is not running, but the caching functionality will not work.
- If you are using a _local_ redis instance, the server can be started by running `redis-server` in the terminal, and you can stop it by running `redis-cli shutdown`. 
- If you are using a _hosted_ redis instance, you can manage the cache through RedisInsight or a similar tool. Our suggestion is to use the `redis-cli` tool, which can be downloaded from the official Redis website or through a package manager such as `apt` or `pacman`. RedisLabs provides the command to connect to the database, which can be found by clicking on the _Connect_ button in the RedisLabs' Databases dashboard.
- These are some of the commands that you can use to manage the cache:
    - `SET key value` - Set a key-value pair in the cache
    - `GET key` - Get the value of a key
    - `DEL key` - Delete a key-value pair
    - `FLUSHALL` - Delete all key-value pairs in all databases

#### Google OAuth 2.0

Follow the instructions outlined in _Setting up Google OAuth 2.0_ (located in the team drive).

<https://drive.google.com/drive/folders/1PWzpsJGXIDA_RnRRoEcJe_U5yvGC6s_U?usp=sharing>

## How to run

- `go build` (install dependencies and build project)
- `go run .` (run server)
- `CompileDaemon -command="./backend"` (run with Hot-reload)

## How to tests
- Once the server is running open another terminal and change directories to tests and run "go test -v" to run all tests functions, or you can run "go test -run 'test_function_name'" to run individual test functions.

## How to tests backend Server hosted in cloud
- Send requests as usual except replace "localhost:8000" with "https://backend-h3ks.onrender.com" example:
"https://backend-h3ks.onrender.com/page-get/123"

