**Usage**

**First step:**

install and configure [go-task](https://github.com/go-task/task) on your machine. [installation instructions](https://taskfile.dev/installation/)

Note: you always use taskfile on local development but not recommended to use on staging and production environments because it adds extra dependencies on servers

**Setup ENVs**

take a look at `.env.sample`

rename `dotenv: ['.env.sample']` to `dotenv: ['.env']` to start using it locally.

**Run Docker Compose:**
to run docker compose commands do the following

```
task compose -- up
task compose -- logs
task compose -- restart
task compose -- down
```

**To run server:**

the following will do cd into server and run main.go

```
task server
```
