# Chix

**Chix** - представляет собой изящный REST API веб-фреймворк, основанный на простом и функциональном HTTP маршрутизаторе Chi. Главной целью Chix является создание легковесного фреймворка, предоставляющего только необходимые компоненты для эффективной работы с проектами REST API. Встроенная Луковичная архитектура облегчает поддержку кода, его масштабирование и тестирование, делая приложение гибким к будущим изменениям.

## Команды проекта

Перечень доступных команд для развертывания проекта.

---

Автоматическая компиляция проекта при изменении файлов (режим локальной разработки).

```bash
make watch
```

---

**Компиляция проекта.**  
По умолчанию, если никакой аргумент после команды не введен, будет проставлен **local**.  
**Важно!** Введенный аргумент должен соответствовать имени конфигурационного файла, без расширения в конце.  
Поддерживаемые расширения конфигурационных файлов **.env** для godotenv и **.yml** для viper.

**В режиме разработки.**  
В этом режиме логи будут отображаться начиная от Debug уровня. 

```bash
go run ./cmd/rest/main.go local
```

**В режиме продакшн.**  
В этом режиме логи будут отображаться начиная от Info уровня.

```bash
go run ./cmd/rest/main.go prod
```

## Используемые пакеты

Перечень сторонних пакетов, используемых в проекте.

---

[modd - автоматическая компиляция при изменении файлов](https://github.com/cortesi/modd)

```bash
go install github.com/cortesi/modd/cmd/modd@latest
```

---

[tint - цветной лог](https://github.com/lmittmann/tint)

```bash
go get github.com/lmittmann/tint
```

---

[viper - yml конфигурация](https://github.com/spf13/viper)

---

```bash
go get github.com/spf13/viper
```

---

[godotenv - переменные окружения](https://github.com/joho/godotenv)

```bash
go get github.com/joho/godotenv
```

---

[sqlite3 - база данных sqlite](https://github.com/mattn/go-sqlite3)

```bash
go get github.com/mattn/go-sqlite3
```

---

[pgx - база данных postgres](https://github.com/jackc/pgx)  
[pgerrcode - коды ошибок postgres](https://github.com/jackc/pgerrcode)

```bash
go get github.com/jackc/pgx/v5
go get github.com/jackc/pgerrcode
```

---

[validator - валидация данных](https://github.com/go-playground/validator)

```bash
go get github.com/go-playground/validator/v10
```

---

[jwt - токен](https://github.com/golang-jwt/jwt)

```bash
go get -u github.com/golang-jwt/jwt/v5
```

---

[chi - роутер](https://github.com/go-chi/chi)  
[cors - защита межсайтового взаимодействия](https://github.com/go-chi/cors)

```bash
go get -u github.com/go-chi/chi/v5
go get github.com/go-chi/cors
```
