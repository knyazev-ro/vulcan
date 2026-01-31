# **Vulcan ORM**

**Vulcan ORM** is an SQL Query Builder + struct-based ORM for Go, with syntax inspired by the Laravel Query Builder.
The project is part of the **Gerard** framework and is designed for **safe, predictable, and deterministic SQL generation**, with support for nested conditions, join operations, bindings, and PostgreSQL.

Vulcan implements the **Data Mapper** principle in the ORM world.
Structs act as instructions describing **how data will be fetched**. This creates a convenient and explicit contract: you don’t need to execute a query to understand the shape of the result.
The way a struct is defined *before* executing a query guarantees the structure of the output.

---

## Supported Features

* `SELECT`
* `WHERE` / `OR WHERE`
* Nested conditions via `WhereClause` / `OrWhereClause`
* `INNER JOIN`, `LEFT JOIN`, `RIGHT JOIN`
* `ORDER BY`
* `LIMIT`, `OFFSET`
* `INSERT` / `CREATE`
* `UPDATE`
* PostgreSQL placeholders (`$1`, `$2`, …)
* Integration with `database/sql`

---

## ORM Features

Vulcan includes an ORM layer:

* Automatic `SELECT` generation based on structs
* Relation support:

  * `has-many`
  * `belongs-to`
  * `has-one`
  * `many-to-many` (via pivot tables)
* Automatic preload generation
* Recursive hydration of nested structures
* Grouping via nested queries

---

## Basic Example

Uses generics:

```go
type UserTest struct {
    _        string `type:"metadata" table:"users" pk:"id"`
    Id       int64  `type:"column" col:"id"`
    Name     string `type:"column" col:"name"`
    LastName string `type:"column" col:"last_name"`
}

vulcan.NewQuery[UserTest]().
    Where("id", ">", 1).
    Where("id", "!=", 3).
    OrderBy([]string{"id"}, "desc").
    Load()
```

`Select` is now optional — by default, the column list is generated from the struct, and this is the **preferred approach**.
The idea is that structs define the output, so using `Select` in Vulcan is now considered legacy and will be removed in the future.

Generated SQL:

```sql
SELECT "users"."id" AS users_id, "users"."name" AS users_name, "users"."last_name" AS users_last_name
FROM users
WHERE "users"."id" > $1 AND "users"."id" != $2
ORDER BY "users"."id" DESC;
```

Bindings:

```
[1, 3]
```

---

## Nested Conditions

```go
vulcan.NewQuery[UserTest]().
    Where("status", "=", 1).
    WhereClause(func(q *vulcan.Query[UserTest]) {
        q.
            Where("age", ">", 18).
            OrWhereClause(func(q *vulcan.Query[UserTest]) {
                q.
                    Where("role", "=", "admin").
                    Where("last_login", ">", "2026-01-01")
            })
    }).
    Where("active", "=", 1).
    Build().
    SQL()
```

Result:

```sql
SELECT "users"."id"
FROM users
WHERE "status" = $1
AND ("age" > $2 OR ("role" = $3 AND "last_login" > $4))
AND "active" = $5;
```

Bindings:

```
[1, 18, "admin", "2026-01-01", 1]
```

---

## JOIN (Manual Mode)

In the current version, if you need a flat result structure, it is enough to declare a single struct and specify the `table:<table_name>` tag for fields mapped from related tables.

```go
q := vulcan.NewQuery[UserTest]().
    InnerJoin("posts", func(jc *vulcan.Join) {
        jc.On("posts.user_id", "=", "users.id")
    }).
    LeftJoin("categories", func(jc *vulcan.Join) {
        jc.On("categories.id", "=", "posts.category_id")
    }).
    Where("users.active", "=", 1).
    Build().
    SQL()
```

Result:

```sql
SELECT <all fields defined in the struct>
FROM users
INNER JOIN posts ON posts.user_id = users.id
LEFT JOIN categories ON categories.id = posts.category_id
WHERE "users.active" = $1;
```

---

## ORM Relations

Vulcan can automatically build and load related models based on **nested structs**.

### Example: One-to-Many + Many-to-Many

```go
type TagTest struct {
    _    string `type:"metadata" table:"tags" pk:"id"`
    Id   int64  `type:"column" col:"id"`
    Name string `type:"column" col:"name"`
}

type PostTag struct {
    _      string  `type:"metadata" table:"post_tags" pk:"post_id,tag_id" tabletype:"pivot"`
    PostId int64   `type:"column" col:"post_id"`
    TagId  int64   `type:"column" col:"tag_id"`
    Tag    TagTest `type:"relation" table:"tags" reltype:"belongs-to" fk:"tag_id" originalkey:"id"`
}

type PostTest struct {
    _        string    `type:"metadata" table:"posts" pk:"id"`
    Id       int64     `type:"column" col:"id"`
    Name     string    `type:"column" col:"name"`
    UserId   int64     `type:"column" col:"user_id"`
    PostTags []PostTag `type:"relation" table:"post_tags" reltype:"has-many" fk:"post_id" originalkey:"id"`
}

type UserTest struct {
    _        string     `type:"metadata" table:"users" pk:"id"`
    Id       int64      `type:"column" col:"id"`
    Name     string     `type:"column" col:"name"`
    LastName string     `type:"column" col:"last_name"`
    Posts    []PostTest `type:"relation" table:"posts" reltype:"has-many" fk:"user_id" originalkey:"id"`
}
```

Query:

```go
vulcan.NewQuery[UserTest]().Load()
```

Vulcan automatically generates **4 queries**
and assembles the result into nested structures:

```
User
 └── Posts
      └── PostTags
           └── Tag
```

---

## UPDATE with JOIN

```go
vulcan.NewQuery[UserTest]().
    From("posts").
    On("posts.id", "=", "users.post_id").
    Where("users.id", "=", 10).
    LeftJoin("categories", func(jc *vulcan.Join) {
        jc.On("categories.id", "=", "posts.category_id")
    }).
    Where("categories.name", "like", "%Tech%").
    Update(map[string]any{
        "users.role_id":  1,
        "users.owner_id": 2,
    })
```

Generated SQL:

```sql
UPDATE users
SET users.role_id = $1, users.owner_id = $2
FROM posts
LEFT JOIN categories ON categories.id = posts.category_id
WHERE posts.id = users.post_id
AND users.id = $3
AND categories.name LIKE $4;
```

Bindings:

```
[1, 2, 10, "%Tech%"]
```

---

## Data Contract

A core architectural principle of Vulcan is that **entity structs are the data output contract**.

This means:

* Model structs define the **expected shape of query results**
* Changing a struct automatically changes the data semantics
* The Query Builder does not “guess” types or perform implicit mapping
* Responsibility for correctness lies in struct definitions

This approach makes the system transparent, predictable, and easy to maintain.

---

## Security

* All values are passed via placeholders
* Values are stored separately in `Bindings`
* Concatenation of user input into SQL is forbidden

> Column and table names are considered trusted, as in Laravel / GORM.

---

## Architectural Decisions

* SQL is built deterministically
* Nested conditions are implemented via builder state control
* Parentheses are formed logically via `WhereClause` and `OrWhereClause`
* No AST or SQL interpretation layer

This is a deliberate tradeoff between control, performance, and implementation complexity.

---

## Limitations

* No multi-dialect SQL support — PostgreSQL is the primary target

---

## New Features

New methods such as `DeleteById`, `FindById`, and most importantly in the current version — **subquery mutation support via `With`** — have been added.

Currently supported:

* Relation-level filtering via closures on the parent struct
* Nested `With` calls

Important note:
The `With` method **only affects how relations are loaded**.
Whether a relation exists in the final result is determined **by the struct**, not the query.

```go
q2, _ := vulcan.NewQuery[ReportData]().With("City", func(q *vulcan.Query[ReportData]) {
    q.Where("city", "like", "Москва")
}).FindById(2)
```

Model relations, if present, are now loaded **concurrently using goroutines**.

---

## Benchmarks (Laravel Eloquent ORM vs Vulcan)

**Dataset:** 23,000 records
**Relations:** 7 × `belongsTo`, 1 × `hasMany`
**Baseline:** Laravel Eloquent ORM (2.9 s)

| ORM / Strategy                                                                   | Time   | Speedup         |
| -------------------------------------------------------------------------------- | ------ | --------------- |
| Laravel **Eloquent ORM**                                                         | 2.9 s  | 1.0×            |
| **Vulcan** (v1, relations via `JOIN`)                                            | 7.0 s  | 0.41× (slower)  |
| **Vulcan** (`WHERE ANY` eager loading)                                           | 600 ms | **4.8× faster** |
| **Vulcan** (`WHERE ANY` + concurrent relation loading)                           | 500 ms | **5.8× faster** |
| **CURRENT Vulcan** (`WHERE ANY` + concurrent loading, post-optimized, NULL-safe) | 300 ms | **9.7× faster** |

---

## Benchmarks (Go ORM: GORM vs Vulcan)

### Case 1 — Full load without relations

**Dataset:** 107,536 records
**Query:** `SELECT * FROM report_data`
**Baseline:** GORM (2.27 s)

| ORM        | Time       | Speedup         | Notes                                           |
| ---------- | ---------- | --------------- | ----------------------------------------------- |
| **GORM**   | 2.27 s     | 1.0×            | `extended protocol limited to 65535 parameters` |
| **Vulcan** | **1.19 s** | **1.9× faster** | Full table load                                 |

---

### Case 2 — Heavy relations load

**Dataset:** 62,000 records
**Relations:** 8 × `belongsTo`, 1 × `hasMany`
**Baseline:** GORM (1.57 s)

| ORM        | Time       | Speedup          | Notes                        |
| ---------- | ---------- | ---------------- | ---------------------------- |
| **GORM**   | 1.57 s     | 1.0×             | `SLOW SQL >= 200ms`          |
| **Vulcan** | **843 ms** | **1.86× faster** | Concurrent relations loading |

---

### So…

* Vulcan outperforms GORM in **both flat and relational workloads**
* GORM hits **PostgreSQL extended protocol limits (65k params)**
* Vulcan avoids parameter explosion and scales linearly
* Best result: **~2× faster** on heavy relations, **~2.5× faster** on full table scans

---

## Endurance Test (500 Iterations)

**Test:** 23,000 records with 8 relations (7 × `belongsTo`, 1 × `hasMany`)  
**Source:** [github.com/knyazev-ro/vulcan/tests](https://github.com/knyazev-ro/vulcan/tests)  
**Configuration:** Semaphore 100, PostgreSQL connections 100  

| Iterations | Records per Iteration | Relations | Total Records Processed | Total Time | Concurrency | Notes |
|------------|---------------------|-----------|------------------------|------------|-------------|-------|
| 500        | 23,000              | 8 (7 `belongsTo`, 1 `hasMany`) | 138,000,000 (main rows + belongs tables + avg 4 has many per main row) | 221.532 s | 100 goroutines | Endurance test passed without errors |


## Project Status

| Component        | Status        | Notes |
|------------------|---------------|-------|
| Query Builder    | Implemented   | Stable |
| ORM Layer        | Implemented   | Core functionality complete |
| Relations        | Implemented   | Fully supported |
| NULL Support     | Implemented   | Ptr Go idiomatic |
| Context Support  | In progress   | API may change |
| `GROUP BY`       | testing       | now is functioning, but needs more tests |
