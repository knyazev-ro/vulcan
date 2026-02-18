# **Vulcan ORM**
![Version](https://img.shields.io/badge/version-v1.0.0-pink)
![License](https://img.shields.io/github/license/knyazev-ro/vulcan)
[![Docs](https://img.shields.io/badge/docs-online-blue)](https://github.com/knyazev-ro/vulcan-orm/blob/main/docs/1.%20introduction.md)
[![Docs](https://img.shields.io/badge/sandbox-blue)](https://github.com/knyazev-ro/vulcan-orm/blob/main/sandbox/sandbox.md)
![Stars](https://img.shields.io/github/stars/knyazev-ro/vulcan?style=social)

**Vulcan ORM** — это высокопроизводительный SQL Query Builder и ORM на базе структур для Go. Синтаксис вдохновлен Laravel Query Builder, но архитектура построена на принципах **Data Mapper**.

> **Важно:** Это часть экосистемы **Gerard**. Подробная документация по всем методам и продвинутым техникам доступна в директории [`/docs`](https://github.com/knyazev-ro/vulcan-orm/blob/main/docs/1.%20introduction.md)

> Для того, чтобы протестировать возможности Vulcan самостоятельно [`/sandbox`](https://github.com/knyazev-ro/vulcan-orm/blob/main/sandbox/sandbox.md)
---

## Архитектурная философия

В отличие от традиционных Active Record ORM, Vulcan реализует принцип **Data Contract**:

* **Структура — это инструкция:** Модель описывает не только таблицу, но и то, как данные будут извлечены.
* **Детерминизм:** Вы знаете форму результата еще до выполнения запроса. Изменение структуры автоматически меняет семантику SQL.
* **Zero Over-fetching:** По умолчанию Vulcan выбирает только те колонки, которые описаны в вашей DTO.

---

## Ключевые возможности

* **Concurrent Hydration:** Рекурсивная загрузка графа отношений (`CLoad`) с использованием goroutines.
* **Smart Relations:** Поддержка `has-one`, `has-many`, `belongs-to` и `many-to-many` (pivot tables).
* **Aggregate Relations:** Встроенная поддержка агрегаций (`count`, `sum`, `avg`, `min`, `max`) прямо в отношениях.
* **Expressive Builder:** Вложенные условия (`WhereClause`), сложные JOIN и мутации.
* **Performance First:** Оптимизированные запросы (`WHERE ANY`) позволяют избегать лимитов параметров PostgreSQL (65k).

---

## Быстрый старт

### Определение модели

```go
type User struct {
    _        string    `type:"metadata" table:"users" pk:"id"`
    Id       int64     `type:"column" col:"id"`
    Name     string    `type:"column" col:"name"`
    // Отношение: Vulcan сам поймет, как собрать граф
    Posts    []Post    `type:"relation" table:"posts" reltype:"has-many" fk:"user_id" originalkey:"id"`
}

```

### Выполнение запроса

```go
ctx := context.Background()

// Конкурентная загрузка пользователей и всех их постов
users, err := vulcan.NewQuery[User]().
    Where("active", "=", 1).
    OrderBy("desc", "id").
    CLoad(ctx)

```

---

## Benchmarks

Vulcan оптимизирован для тяжелых нагрузок и больших наборов данных.

### Vulcan vs Laravel Eloquent

*Dataset: 23,000 records, 8 relations*

| Strategy | Time | Speedup |
| --- | --- | --- |
| Laravel **Eloquent ORM** | 2.9 s | 1.0× |
| **Vulcan** (Подгрузка отношений на базе JOIN) | 7.0 s | 0.41× |
| **Vulcan** (Имплементация отношений на базе WHERE ANY) | 600 ms | 4.8× faster |
| **Vulcan (Горутины + оптимизация)** | **300 ms** | **9.7× faster** |

### Vulcan vs GORM (Go)

*Heavy relations load (62,000 records)*

| ORM | Time | Speedup | Notes |
| --- | --- | --- | --- |
| **GORM** | 1.57 s | 1.0× | Slow SQL warnings |
| **Vulcan** | **843 ms** | **1.86× faster** | Concurrent hydration |

---

## Продвинутые возможности

### Агрегации в отношениях

Vulcan автоматически группирует данные, если видит тег `agg`.

```go
type UserStats struct {
    _          string `type:"metadata" table:"posts"`
    TotalPosts int64  `type:"column" col:"*" agg:"count"`
}

```

### Фильтрация отношений (With)

Метод `With` позволяет накладывать условия на вложенные данные, не затрагивая основной запрос.

```go
q, _ := vulcan.NewQuery[User]().
    With("Posts", func(q *vulcan.Query[User]) {
        q.Where("status", "=", "published")
    }).FindById(ctx, 1)

```

### UPDATE с JOIN

```go
vulcan.NewQuery[User]().
    From("posts").
    On("posts.user_id", "=", "users.id").
    Where("posts.views", ">", 1000).
    Update(ctx, map[string]any{"is_popular": true})

```

---

## Безопасность и ограничения

* **SQL Injection:** Все значения проходят через плейсхолдеры PostgreSQL (`$1`, `$2`). Конкатенация пользовательского ввода запрещена.
* **Диалекты:** Основная и единственная цель — **PostgreSQL**.
* **Trusted Names:** Имена таблиц и колонок считаются доверенными (задаются в коде через теги).

---

## Project Status

| Component | Status |
| --- | --- |
| **Query Builder Core** | ✅ Stable |
| **ORM / Relations** | ✅ Implemented |
| **Concurrent Hydration** | ✅ Implemented |
| **Composite PK** | ✅ Implemented |
| **Context & Semaphores** | ✅ Implemented |
| Chunk + Each | *in plan* |

---

## License

Vulcan ORM is open-sourced software licensed under the MIT license.

---

*Developed by @knyazev-ro as part of the Gerard ecosystem.*