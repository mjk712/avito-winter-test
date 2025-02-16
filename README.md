# **Магазин мерча для сотрудников Avito**

[![CI](https://github.com/mjk712/avito-winter-test/actions/workflows/ci.yaml/badge.svg)](https://github.com/mjk712/avito-winter-test/actions/workflows/ci.yaml)
[![Coverage Status](https://coveralls.io/repos/github/mjk712/avito-winter-test/badge.svg?branch=main)](https://coveralls.io/github/mjk712/avito-winter-test?branch=main)

Этот проект представляет собой сервис для внутреннего магазина мерча компании Avito. Сотрудники могут использовать монеты для покупки товаров и передачи монет друг другу. Сервис предоставляет API для управления балансом, покупками и историей транзакций.

---

## **Функциональность**

- **Покупка мерча**: Сотрудники могут покупать товары из магазина за монеты.
- **Передача монет**: Сотрудники могут передавать монеты друг другу.
- **История транзакций**: Каждый сотрудник может просматривать историю полученных и отправленных монет.
- **Инвентарь**: Сотрудники могут просматривать список купленных товаров.
- **Авторизация**: Используется JWT для аутентификации и авторизации пользователей.

---

## **Технические детали**

### **Стек технологий**
- **Язык программирования**: Go
- **База данных**: PostgreSQL
- **Контейнеризация**: Docker и Docker Compose
- **Тестирование**: Юнит-тесты, интеграционные тесты
- **CI/CD**: GitHub Actions
- **Линтинг**: Настроен через `.golangci.yaml`

### **API**
Подробное описание API находится в папке [docs](/docs).

---

## **Инструкция по запуску**

### **1. Настройка окружения**
Сервис настраивается с помощью следующих переменных окружения:

- `SERVER_ADDRESS` — адрес сервера (по умолчанию `:8080`)
- `POSTGRES_USERNAME` — имя пользователя PostgreSQL
- `POSTGRES_PASSWORD` — пароль PostgreSQL
- `POSTGRES_HOST` — хост PostgreSQL
- `POSTGRES_PORT` — порт PostgreSQL
- `POSTGRES_DATABASE` — имя базы данных
- `POSTGRES_CONN` — строка подключения к PostgreSQL
- `ENV` — окружение (например, `dev` или `prod`)
- `JWT_SECRET` — секретный ключ для JWT

### **2. Запуск сервиса**
1. Убедитесь, что у вас установлены Docker и Docker Compose.
2. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/mjk712/avito-winter-test.git
   cd avito-winter-test
3. Запустите сервис с помощью Docker Compose:
    ```bash
   docker-compose up
4. Сервис будет доступен по адресу http://localhost:8080.

### **Тестирование**
**Юнит-тесты**

- Проект покрыт юнит-тестами, которые проверяют основные сценарии работы сервиса. 
- Общее покрытие тестами превышает 40%.

**Интеграционные тесты**

- Реализованы интеграционные тесты для основных сценариев.

**Запуск тестов**
   ```bash
   go test ./...

