## Deps
1. Postgres
2. Go

## Launch

```bash
    $ docker run --name postgres -p 5432:5432 -e POSTGRES_INITDB_ARGS="--data-checksums" -v $(PWD)/data/db:/var/lib/postgresql/data -d postgres:11

    $ docker exec -it postgres bash

    $ psql -U postgres

    $ create table tasks (id INT GENERATED ALWAYS AS IDENTITY primary key, created_at TIMESTAMP DEFAULT NOW() NOT NULL, rounds int not null default 1);

    $ create table tasks_hashes (id int generated always as identity primary key, task_id int references tasks(id), hash varchar not null, created_at timestamp default now() not null, round int not null);
```

```bash
    $ ./main.go #this command run app on 3001 port (as default)
```
