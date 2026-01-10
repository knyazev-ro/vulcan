# Vulcan Query Builder

**Vulcan Query Builder** — это SQL query builder для Go с синтаксисом, вдохновлённым Laravel Query Builder. Проект является частью фреймворка Gerard и предназначен для безопасной и детерминированной генерации SQL-запросов с поддержкой вложенных условий, join, биндингов и PostgreSQL.

## Поддерживаемые возможности

* `SELECT`
* `WHERE` / `OR WHERE`
* Вложенные условия через `WhereClause` / `OrWhereClause`
* `INNER JOIN`, `LEFT JOIN`
* `ORDER BY`
* `LIMIT`, `OFFSET`
* `INSERT` / `CREATE`
* `UPDATE`
* PostgreSQL placeholders (`$1`, `$2`, …)
* Интеграция с `database/sql`

## Пример использования

```go
user := NewUser()

query.NewQuery(user.model).
    Select([]string{"id", "name"}).
    Where("id", ">", 1).
    Where("id", "!=", 3).
    OrderBy([]string{"id"}, "desc").
    Build().
    Get()
```

Сгенерированный SQL:

```sql
SELECT "id", "name" FROM users WHERE "id" > $1 AND "id" != $2 ORDER BY "id" DESC;
```

Bindings:

```text
[1, 3]
```

---

## Вложенные условия

Для работы со скобками используется функциональный подход через `WhereClause` и `OrWhereClause`:

```go
query.NewQuery(user.model).
    Select([]string{"id"}).
    Where("status", "=", 1).
    WhereClause(func(q *query.Query) {
        q.Where("age", ">", 18).
          OrWhereClause(func(q *query.Query) {
              q.Where("role", "=", "admin").
                Where("last_login", ">", "2026-01-01")
          })
    }).
    Where("active", "=", 1).
    Build().
    SQL()
```

Результат:

```sql
SELECT "id" FROM users
WHERE "status" = $1 AND ("age" > $2 OR ("role" = $3 AND "last_login" > $4)) AND "active" = $5;
```

Bindings:

```text
[1, 18, "admin", "2026-01-01", 1]
```

---

## JOIN

Поддерживаются `INNER JOIN` и `LEFT JOIN` с кастомными `On`-условиями:

```go
q := query.NewQuery(user.model).
    Select([]string{"users.id", "posts.title", "categories.name"}).
    InnerJoin("posts", func(jc *query.Join) {
        jc.On("posts.user_id", "=", "users.id")
    }).
    LeftJoin("categories", func(jc *query.Join) {
        jc.On("categories.id", "=", "posts.category_id")
    }).
    Where("users.active", "=", 1).
    Build().
    SQL()
```

Результат:

```sql
SELECT "users.id", "posts.title", "categories.name"
FROM users
INNER JOIN posts ON posts.user_id = users.id
LEFT JOIN categories ON categories.id = posts.category_id
WHERE "users.active" = $1;
```

---

## UPDATE

Пример обновления данных с join и вложенными условиями:

```go
query.NewQuery(user.model).
    From("posts").
    On("posts.id", "=", "users.post_id").
    Where("users.id", "=", 10).
    LeftJoin("categories", func(jc *query.Join) {
        jc.On("categories.id", "=", "posts.category_id")
    }).
    Where("categories.name", "like", "%Tech%").
    Update(map[string]any{
        "users.role_id":  1,
        "users.owner_id": 2,
    })
```

Сгенерированный SQL:

```sql
UPDATE users
SET users.role_id = $1, users.owner_id = $2
FROM posts
LEFT JOIN categories ON categories.id = posts.category_id
WHERE posts.id = users.post_id AND users.id = $3 AND categories.name LIKE $4;
```

Bindings:

```text
[1, 2, 10, "%Tech%"]
```

---

## Безопасность

* Все значения передаются через placeholders
* Значения хранятся отдельно в `Bindings`
* Конкатенация пользовательских значений в SQL запрещена

> Имена колонок и таблиц считаются доверенными, как в Laravel / GORM.

---

## Архитектурные решения

* SQL собирается детерминированно
* Вложенность условий реализована через управление состоянием билдера
* Скобки формируются логически через `WhereClause` и `OrWhereClause`
* Нет AST и нет попытки интерпретации SQL

Это осознанный компромисс между контролем и сложностью реализации.

---

## Ограничения

* Проект не является полноценным ORM
* Нет автоматического маппинга результатов в структуры
* Нет поддержки разных диалектов SQL

---

## Статус проекта

* Query Builder: реализован и протестирован
* ORM-уровень: в разработке
* API может изменяться
