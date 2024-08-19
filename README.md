# Odin Leather Store

An application for Leather Crafting!

More details soon...

## Develop

Run this command to start the server and watching for change go and templ
changes

```bash
docker compose --profile dev up --watch
```

And to watch CSS, **inside the container** you might need to run `make css-dev` outside the container too

```bash
# Enter the container
docker exec -it odinls-dev bash

# Run the tailwind css generator
make css-dev
```

### Seed the database

You might need to seed the database with the roles and such. To do so run:

```bash
# Enter the container
docker exec -it odinls-dev bash

# Seed the database
go run cmd/seeder/main.go
```
