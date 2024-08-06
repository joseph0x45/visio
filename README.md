![logo](https://github.com/TheWisePigeon/visio/assets/95161388/de1bc44d-d238-4742-903e-c744c9931d5c)

## Introduction
Visio is an opensource cloud based service that provides algorithms for face detection and recognition. 

## Development

### Without Docker
If you can get [go-face](https://github.com/Kagami/go-face) to compile on your system you are good. Checkout the dependencies required to do that 
in the [library's README](https://github.com/Kagami/go-face#requirements).
Once you installed the required dependencies, all you will need is to add the URL to a Postgres and Redis database in your .env.

```
PG_URL="postgres://...."
REDIS_URL="redis://...."
PORT="8080"
```

### Using Docker
If you are struggling to install the dependencies or use a system that doesn't quite support them (like Manjaro in my case), you have to use Docker.

```
PG_URL="postgres://...."
REDIS_URL="redis://...."
PORT="8080"
CODE_DIR="/path/to/your/code"
```

After setting the environment variables in the `.env` file, run `docker compose up -d` or `docker-compose up -d` to launch the services. The Go service comes
with [air](https://github.com/cosmtrek/air) and the dependencies required to compile all preinstalled. cd into the root project folder and run `air` to compile and 
start the server.

### Last steps
You can find the database schema in the `schema.sql` at the root of the project folder. If you are using Docker to develop you can run `make migrate` to quickly
apply the schema to your database.

The project uses tailwindcss so you will also need to run `npm install` and start the tailwind compilation with `make start-tailwind-compilation`.
If you do not have `make` you can always run it with `npx tailwindcss -i ./assets/app.css -o ./public/output.css --minify --watch`

If you have any issues while setting it up locally join the [Discord](https://discord.gg/9vDumSjK3F) and I will be happy to help :)
