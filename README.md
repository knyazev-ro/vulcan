# **Vulcan ORM**

**Vulcan ORM** — это SQL vulcan builder + struct based ORM для Go с синтаксисом, вдохновлённым Laravel Query Builder. Проект является частью фреймворка **Gerard** и предназначен для безопасной, предсказуемой и детерминированной генерации SQL-запросов с поддержкой вложенных условий, join-операций, биндингов и работы с PostgreSQL.

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


## ORM-возможности (в разработке)

Vulcan включает экспериментальный ORM-уровень:

* Автоматическая генерация `SELECT` на основе структуры
* Поддержка relations:

  * `has-many`
  * `belongs-to`
  * `has-one`
  * `many-to-many` (через pivot-таблицу)
* Автоматическая генерация `JOIN` по тегам структуры
* Рекурсивная гидратация вложенных структур
* Группировка плоских SQL-строк в иерархические Go-объекты


## Базовый пример (новый API)

Вместо передачи `user.model` теперь используется **дженерик**:

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
    Build().
    Get()
```

 `Select` теперь опционален — по умолчанию список колонок генерируется из структуры.

Сгенерированный SQL:

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

## Вложенные условия

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

Результат:

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

## JOIN (ручный режим)

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

Результат:

```sql
SELECT ...
FROM users
INNER JOIN posts ON posts.user_id = users.id
LEFT JOIN categories ON categories.id = posts.category_id
WHERE "users.active" = $1;
```

---

## ORM Relations (автоматические JOIN)

Vulcan может автоматически строить JOIN на основе тегов структуры.

### Пример: One-to-Many + Many-to-Many

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
    PostTags []PostTag `type:"relation" table:"post_tags" reltype:"has-many" fk:"post_id"`
}

type UserTest struct {
    _        string     `type:"metadata" table:"users" pk:"id"`
    Id       int64      `type:"column" col:"id"`
    Name     string     `type:"column" col:"name"`
    LastName string     `type:"column" col:"last_name"`
    Posts    []PostTest `type:"relation" table:"posts" reltype:"has-many" fk:"user_id"`
}
```

Запрос:

```go
vulcan.NewQuery[UserTest]().
    Build().
    Get()
```

Vulcan автоматически сгенерирует:

```
users
LEFT JOIN posts
LEFT JOIN post_tags
LEFT JOIN tags
```

И соберёт результат в вложенные структуры:

```
User
 └── Posts
      └── PostTags
           └── Tag
```

---

## 🔹 UPDATE с JOIN

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

Сгенерированный SQL:

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

## Контракт данных

Ключевой архитектурный принцип Vulcan заключается в том, что **структуры сущностей являются контрактом выходных данных**.

Это означает:

* Структуры моделей определяют **ожидаемую форму результата запроса**.
* Изменение структуры автоматически меняет семантику данных, которые можно получить.
* Query Builder не «угадывает» типы и не выполняет неявный маппинг — ответственность за корректность лежит на описании структур.
* Такой подход делает поведение системы прозрачным, предсказуемым и удобным для сопровождения.

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
* Отсутствует AST и интерпретация SQL

Это осознанный компромисс между контролем, производительностью и сложностью реализации.

---

## Ограничения

* Проект не является полноценным ORM в классическом смысле
* Нет автоматического «магического» маппинга результатов в произвольные структуры
* Нет поддержки множества SQL-диалектов — приоритет отдан PostgreSQL

---

## Статус проекта

* Query Builder: реализован и стабилизирован
* ORM-уровень: **активно дополняется и улучшается по мере новых коммитов**
* API может изменяться
