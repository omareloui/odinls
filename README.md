# Odin Leather Store

An application for Leather Crafting!

More details soon...

## Develop

Start the server and watching for change go and templ changes

```bash
docker compose --profile dev up --watch --build
```

To generate mock files for development, you can run:

```bash
mockery
```

Install node dependencies then generate the CSS files and watch for changes

```bash
pnpm install
pnpm css:dev
```

### Seed the database

You might need to seed the database with the roles and such. To do so run:

```bash
# Enter the container
docker exec -it odinls-dev bash

# Seed the database
go run cmd/seeder/main.go
```
