<div align="center">
  <h1>🔒 Confessly</h1>
  <h3>Анонимная платформа для откровений</h3>
  
  [![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)](https://golang.org/)
  [![Gin Framework](https://img.shields.io/badge/Gin-1.8.1-00ADD8?style=flat&logo=go)](https://github.com/gin-gonic/gin)
  [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-4169E1?style=flat&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
  
  [![Swagger](https://img.shields.io/badge/Swagger-85EA2D?style=flat&logo=swagger&logoColor=white)](/swagger/index.html)
</div>

## 📝 О проекте

**Confessly** — это высоконадежная платформа для анонимных признаний, построенная на Go с использованием Gin и PostgreSQL. Позволяет пользователям делиться своими мыслями анонимно, обеспечивая при этом безопасность и конфиденциальность.

### 🌟 Ключевые особенности

- 🔐 **Аутентификация** через JWT токены
- 📱 **RESTful API** с полной документацией Swagger
- 🛡️ **Модерация контента** с системой жалоб
- 👥 **Гостевой доступ** к просмотру признаний
- ⚡ **Высокая производительность** благодаря Golang
- 📊 **Подробное логирование** всех операций

## 🚀 Быстрый старт

### 📋 Требования

- Go 1.19+
- PostgreSQL 13+
- Git

### ⚙️ Установка

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/hadisjane/confessly.git
   cd confessly
   ```

2. Создайте файл `.env` в корне проекта:
   ```env
   DB_PASSWORD=your-password
   JWT_SECRET_KEY=your-secret-key
   ```

3. Установите зависимости:
   ```bash
   go mod download
   ```

4. Запустите приложение:
   ```bash
   go run main.go
   ```

5. Откройте документацию API:
   ```
   http://localhost:8081/swagger/index.html
   ```

## 📚 API Endpoints

### 🔐 Аутентификация

| Метод | Эндпоинт | Описание |
|-------|----------|-----------|
| `POST` | `/auth/register` | Регистрация нового пользователя |
| `POST` | `/auth/login` | Вход в систему |

### 📝 Признания

| Метод | Эндпоинт | Описание |
|-------|----------|-----------|
| `GET` | `/public/confessions` | Получить список публичных признаний |
| `GET` | `/public/confessions/:id` | Получить признание по ID |
| `POST` | `/public/confessions` | Создать новое анонимное признание |
| `PUT` | `/api/confessions/:id` | Обновить признание (только автор) |
| `DELETE` | `/api/confessions/:id` | Удалить признание (только автор) |

### 🚨 Модерация

| Метод | Эндпоинт | Описание |
|-------|----------|-----------|
| `GET` | `/api/admin/reports` | Список жалоб (админ) |
| `POST` | `/api/reports` | Пожаловаться на признание |
| `PUT` | `/api/admin/reports/:id` | Обновить статус жалобы (админ) |
| `DELETE` | `/api/admin/confessions/:id` | Удалить признание (админ) |

## 🏗️ Архитектура

```
confessly/
├── internal/            # Внутренние пакеты
│   ├── configs/         # Конфигурация приложения
│   ├── controller/      # HTTP обработчики
│   ├── db/              # Работа с базой данных
│   ├── errs/            # Кастомные ошибки
│   ├── middleware/      # Промежуточное ПО
│   ├── models/          # Модели данных
│   ├── repository/      # Слой доступа к данным
│   └── service/         # Бизнес-логика
├── logger/              # Логирование
├── utils/               # Вспомогательные утилиты
├── docs/                # Документация Swagger
├── .env                 # Переменные окружения
├── go.mod               # Зависимости
└── main.go              # Точка входа
```

## 🔧 Технологии

- **Backend**: Go 1.19+
- **Фреймворк**: Gin
- **База данных**: PostgreSQL
- **Аутентификация**: JWT
- **Документация**: Swagger
- **Логирование**: Кастомный логгер
- **Конфигурация**: Переменные окружения и конфиги

## 📬 Контакты

- **Автор**: [hadisjane](https://github.com/hadisjane)

---

> 🧠 Данный проект создан для практики и изучения чистого Go, REST API и архитектурных подходов. А так же для понимания работы с базами данных и веб-серверами.