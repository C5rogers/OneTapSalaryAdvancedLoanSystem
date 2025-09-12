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

# How Validation Works

Validating and migrating customers to database is done on the following route

```
/api/validate_customers
```

It follows the following steps

1. read customers from `data/customers.json` file
2. load sample customer from `data/sample_customer.csv` file
3. validate customers read from json file against sample customer
4. if valid and the record doesn't exist on db, save to database
5. if not valid, take an error log for that specific record
6. finally map all customers with validity and also save the log under `logs/validation_log_<timestamp>.json` file
7. finally return the log data to the requesting user

# How Rating Calculation Works

Under the following route rating calculation is generated based on the customer transaction history

```
/api/process_transactions
```

It follows the following steps

1. read transactions from `data/transactions.json` file
2. load all customers from database
3. for each customer, filter transactions that belong to that specific customer
4. calculate rating based on the following rules

- `total` transaction : users total number of transactions
  - count all transactions
- `volumeScore`: total amount of transactions volume
  - minimum of `10` and `total/1000`
- `durationScore`: the average time gap between customer transactions
  - time range of the last and the recent transactions
- `stabilityScore`: the customer balance stability scoring calculation
  - checking the minimum balance with 0 and adding score of `2`
- `final`: final rating score
  - percentage multiplication of each scores and sum them up

# Security Measures

- applying `JWT` authentication system with role based token signing

  - the token is signed with `ES256` algorithm private key and decoded with `ES256` public key
  - using 2 key approach allow the jwt generated on this backend can also be used for another backend system by just sharing the public key only

- role level checking for every processing action
- basic authentication system using `email` and `password` system
- password is hashed using `bcrypt` algorithm which is a one way hashing algorithm

# Finally

The business logic routes are the following

```
/auth/login
/auth/register
/api/validate_customers
/api/process_transactions
```

and also the basic postman collection is under `one-tap.postman_collection.json` file

# Enjoy :)
