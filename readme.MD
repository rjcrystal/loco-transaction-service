# Transactions service

## APIs 
- Put Transaction  `/transactionservice/transaction/:transactionId` 
- Get Transactions by type `/transactionservice/types/:type`
- Get Transaction sum `/transactionservice/sum/:transactionId`
- Get Transaction details `/transactionservice/transaction/:transactionId`
- Check service health  `/health`


## Caveats/Improvements 
1. Transaction ID should be Generated from the backend It can be a serial/bigserial type in postgres or randomly generated UUID, etc.
2. SUM API considerations 
- Circular referencing of transactions will not be an issue as cyclic checks are in place, but it should be avoided e.g. A->B->C->A
- As this API requires recursive joins nesting should not be done more than 5-6 levels deep for best performance, If the default Stack depth of 2 MB is exceeded it will cause the query to be killed and API will fail


### Requirements 
- Postgresql 

### How to run 
- Create a .env file with postgres connection string 
- Crate a table in postgres with the below definition 
- Run `go run handler.go logger.go main.go models.go store.go`

### DDL for transactions table
```
CREATE TABLE public.transactions (
	id int4 NOT NULL,
	amount float8 NULL,
	parent_id int4 NULL,
	"type" text NOT NULL,
	CONSTRAINT transactions_pk_id PRIMARY KEY (id),
	CONSTRAINT transactions_fk FOREIGN KEY (parent_id) REFERENCES public.transactions(id) ON DELETE RESTRICT ON UPDATE RESTRICT
);