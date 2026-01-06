# Vulcan Query Builder

Vulcan Query Builder — это SQL query builder для Go с синтаксисом, вдохновлённым Laravel Query Builder. Проект является частью фреймворка Gerard и предназначен для безопасной и детерминированной генерации SQL-запросов с поддержкой вложенных условий, биндингов и PostgreSQL.

## Поддерживаемые возможности

* SELECT
* WHERE / OR WHERE
* Вложенные условия через скобки
* LEFT JOIN
* ORDER BY
* INSERT
* PostgreSQL placeholders ($1, $2, ...)
* Интеграция с database/sql

## Пример использования

```go
user := NewUser()

query.NewQuery(user.model).
    Select([]string{"id", "name"}).
    Where("id", ">", "1").
    Where("id", "!=", "3").
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

## Вложенные условия

Для работы со скобками используется функциональный подход:

```go
query.NewQuery(user.model).
    Select([]string{"id"}).
    Where("status", "=", "1").
    WhereClause(func(q *query.Query) {
        q.Where("age", ">", "18").
          OrWhere("role", "=", "admin")
    }).
    Build().
    Get()
```

Результат:

```sql
SELECT "id" FROM users WHERE "status" = $1 AND ("age" > $2 OR "role" = $3);
```

## Безопасность

* Все значения передаются через placeholders
* Значения хранятся отдельно в `Bindings`
* Конкатенация пользовательских значений в SQL запрещена

Имена колонок и таблиц считаются доверенными (как и в Laravel / GORM).

## Архитектурные решения

* SQL собирается детерминированно
* Вложенность условий реализована через управление состоянием билдера
* Скобки не парсятся, а формируются логически
* Нет AST и нет попытки интерпретации SQL

Это осознанный компромисс между сложностью и контролем.

## Ограничения

* Проект не является ORM (пока что)
* Нет автоматического маппинга результатов в структуры
* Нет поддержки разных диалектов SQL

## Статус проекта

* Query Builder: реализован и протестирован
* ORM-уровень: в разработке
* API может изменяться