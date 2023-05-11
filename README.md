

# ZenMobile Service

Учебный проект для знакомства с *Redis*, *PostgreSQL*, *HMAC-SHA-512*

### Redis

Клиент Redis хранится в отдельном контейнере Docker.

Реализовано подключение к Redis указывая хост и порт в конфигурационном файле или передать в качестве параметров при запуске сервиса (-host и -port).

#### Функционал

1. Сохранение данных в Redis.

2. Получение данных из Redis по ключу.

3. Инкрементировать значение по ключу, только для целочисленных значений

   

#### Руководство пользователя

1. Чтобы сохранить данные в Redis, нужно отправить запрос к сервису:

   ```http
   POST /redis/ HTTP/1.1
   Content-Type: application/json; charset=utf-8
   Host: localhost:8080
   ```

   Тело запроса в формате JSON

   ```json
   {
       "key": "CountRegisteredUsers",
       "value": 444
   }
   ```

   Если операция сохранения значения выполнена успешно, получаете ответ такого формата:

   ```tex
   Key = CountRegisteredUsers Value = 444 saved in Redis
   ```

2. Чтобы получить данные из Redis по ключу, сделайте запрос:

   ```http
   GET /redis/key HTTP/1.1
   Host: localhost:8080
   ```

    После /redis/ указывается ключ значения который хотите получить.

   Пример ответа:

   ```json
   {
       "value": 444
   }
   ```

3. Инкрементировать значение по ключу можно с помощью запроса:

   ```http
   POST /redis/incr HTTP/1.1
   Content-Type: application/json; charset=utf-8
   Host: localhost:8080
   ```

   Тело запроса в формате JSON (value значение на которое инкрементируем)

   ```json
   {
       "key": "CountRegisteredUsers",
       "value": 4
   }
   ```

   Если операция выполнена успешно, вы получите ответ с новым значением. 

   ```json
   {
       "value": 448
   }
   ```

   Однако, если ключ не найден или значение не является числом, вы получите соответствующее сообщение об ошибке.

### HMAC-SHA-512

Подпись HMAC-SHA-512 реализовано в отдельном сервисе.

Реализовано подпись сообщения в соответствии с алгоритмом HMAC-SHA-512, используя переданный ключ.

#### Функционал

1. Подпись сообщения HMAC-SHA-512 и возрат hex строкой
2. Проверка подписи

#### Руководство пользователя

1. Чтобы подписать сообщение, нужно отправить запрос к сервису

   ```http
   POST /sign/hmacsha512 HTTP/1.1
   Content-Type: application/json; charset=utf-8
   Host: localhost:8080
   ```

   Тело запроса в формате JSON

   ```json
   {
       "text": "test",
       "key": "test123"
   }
   ```

   Если операция подписи выполнена успешно, получаете ответ такого формата:


```tex
9f7a0e6c7f4e026f5d8b3d2f6b652ba8a66a3c7e8b7c4dacc8a8e709f8ffbc63c4f499ae8d3c7b4aa3f73b4a5b3c273dc1d1adfb7c6c6e3e9fad9ecc9d347bc

```

### PostgreSQL

PostgreSQL хранится в отдельном контейнере Docker.

Реализовано подключение к PostgreSQL указывая хост, порт, имя базы данных, имя пользователя и пароль в конфигурационном файле.

#### Функционал

1. Создание таблицы users в базе данных PostgreSQL.

2. Добавление записи пользователя в базу данных.

3. Получение информации о пользователе.

#### Руководство пользователя

1. Чтобы сохранить данные пользователя в PostgreSQL, нужно отправить запрос к сервису:

   ```http
   POST /postgres/users/ HTTP/1.1
   Content-Type: application/json; charset=utf-8
   Host: localhost:8080
   ```

   Тело запроса в формате JSON

   ```json
   {
       "name": "Alex",
       "age": 21
   }
   ```

   В ответ сервис должен вернуть id пользователя в PostgreSQL:

   Пример ответа:

   ```json
   {
       "id": 1
   }
   ```

2. Чтобы получить информацию о пользователе по id сделайте запрос:

   ```http
   GET /postgres/userID HTTP/1.1
   Host: localhost:8080
   ```

    После /postgres/ указывается id пользователя информацию которого хотите получить.

   Пример ответа от сервиса:

   ```json
   {
       "id": 1,
       "name": "Alex",
       "age": 21
   }
   ```

### Запуск сервиса

Создание и запуск сервиса в Docker

```makefile
make build
make up
```

### Тестирование

Для сервиса написаны unit тесты:

Один раз протестировать:

```makefile
make test
```

Запсутить тесты сто раз:

```makefile
make test100
```

Посмотреть покрытие тестами:

```makefile
make cover
```

### Документация 

У сервиса есть описание API. Используется [Swagger](https://github.com/swaggo/http-swagger) 

Генерация документации:

```makefile
make swag
```

Просмотреть описание API:

[SwaggerURL](http://localhost:8080/swagger/index.html)

### Используемые фреймворкики, концепции

- [viper](https://github.com/spf13/viper)
- [go-chi/chi](https://github.com/go-chi/chi)
- [go-redis](https://github.com/redis/go-redis)
- [squirrel](https://github.com/Masterminds/squirrel)
- [pgx](https://github.com/jackc/pgx)
- [logrus](https://github.com/sirupsen/logrus)
- [gomock](https://github.com/golang/mock)

