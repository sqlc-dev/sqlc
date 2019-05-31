This is an attempt to "hand-generate" the interface we'd like to tables, using
a sample table from Equinox.

```sql
CREATE TABLE credentials (
        id         SERIAL       UNIQUE NOT NULL,
        sid        varchar(64)  UNIQUE NOT NULL,
        created    timestamp    DEFAULT NOW(),
        accountid  bigint       NOT NULL,
        tokenhash  varchar(255) NOT NULL
)
```
