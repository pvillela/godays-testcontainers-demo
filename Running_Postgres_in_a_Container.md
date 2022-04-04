# References and Commands to Run Postgres in a Container

### References

- Postgres Docker image
    https://hub.docker.com/_/postgres
    https://github.com/docker-library/docs/blob/master/postgres/README.md
- Accessing PostgreSQL databases in Go
    https://eli.thegreenplace.net/2021/accessing-postgresql-databases-in-go/
    https://github.com/eliben/code-for-blog/tree/master/2021/go-postgresql
- Connecting to Postgresql in a docker container from outside
    https://stackoverflow.com/questions/37694987/connecting-to-postgresql-in-a-docker-container-from-outside
- Connect From Your Local Machine to a PostgreSQL Database in Docker
    https://betterprogramming.pub/connect-from-local-machine-to-postgresql-docker-container-f785f00461a7
- Accessing a PostgreSQL Database in a Docker Container (incl. backup/restore)
    https://inedo.com/support/kb/1145/accessing-a-postgresql-database-in-a-docker-container
- psql Docker image
    https://hub.docker.com/r/jbergknoff/postgresql-client
- Access host from a docker container
    https://dev.to/bufferings/access-host-from-a-docker-container-4099
    https://docs.docker.com/network/host/
    https://docs.docker.com/network/network-tutorial-host/

### Launch Postgres container with a named volume
```bash
docker run -d --rm -p 9999:5432 --name some-postgres -v myPostgresVol:/var/lib/postgresql/data -e POSTGRES_PASSWORD=mysecretpassword postgres
```

### Execute psql as admin on Postgres container
```bash
docker exec -ti some-postgres psql -U postgres
```

### Create test database and test user at psql prompt
postgres=# 
```bash
create database testmooc;
```
postgres=# 
```bash
create role testuser with login password 'testpassword';
```

### Execute psql with test user on client container 
```bash
docker run -it --rm --network host --name psql postgres psql postgresql://testuser:testpassword@localhost:9999/testmooc
```
or
```bash
docker run -it --rm --network host --name psql jbergknoff/postgresql-client postgresql://testuser:testpassword@localhost:9999/testmooc
```

### Once in the testmooc psql prompt

#### Add the databases
By pasting the contents of the files in 
https://github.com/eliben/code-for-blog/tree/master/2021/go-postgresql/migrations

#### Back up database
```bash
docker exec -u postgres some-postgres pg_dump -Cc | xz > some-backup-$(date -u +%Y-%m-%d).sql.xz
```

#### Restore database from backup (replace date below as needed)
```bash
xz -dc some-backup-2022-02-17.sql.xz | docker exec -i -u postgres some-postgres psql –set ON_ERROR_STOP=on –single-transaction 
```

### Other ways to execute psql

#### Execute psql with another user on Postgres container
```bash
docker exec -ti some-postgres psql postgres://testuser:testpassword@localhost/testmooc
```

#### Execute psql with another user on host
```bash
psql postgres://testuser:testpassword@localhost:9999/testmooc
```
