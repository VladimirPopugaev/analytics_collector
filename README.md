# Сервис сбора аналитики

## Описание

Сервис предназначен для сбора и хранение информации о действиях совершенных
пользователями. С какой-то периодичностью другие сервисы отправляют в сервис
аналитики сообщения о действиях пользователя, например: пользователь
авторизовался, пользователь изменил свои настройки, добавил аватарку.

### Пример отправляемых сообщений

```curl
curl -location -request POST 'http://localhost:8080/analitycs' \
--header 'X-Tantum-UserAgent: DeviceID=G1752G75-7C56-4G49-BGFA-5ACBGC963471;DeviceType=iOS' \
--header 'X-Tantum-Authorization: 2daba111-1e48-4ba1-8753-2daba1119a09' \
--header 'Content-Type: application/json' \
--data-raw '{
    "module" : "settings",
    "type" : "alert",
    "event" : "click",
    "name" : "подтверждение выхода",
    "data" : {"action" : "cancel"}
}'
```

Сервис принимает данные через POST запрос по адресу ```/analytics```, валидирует 
принимаемые значения и в случае корректных данных возвращает положительный ответ
со статусом ```202``` и телом:

```json
{
  "status": "OK"
}
```

Далее сервис передаёт работу worker's pool, где уже происходит сохранение в базу данных.

## Конфигурирование сервиса

Сервис поддерживает конфигурирование с помощью yaml файла. Для этого необходимо в
директории ```configs/``` создать файл со следующим содержимым:

```yaml
env: "local"
server:
  host: "localhost"
  port: "8080"
  workers_count: 5
database:
  username: "admin"
  password: "admin"
  db_name: "analysis_db"
  ssl_mode: "disable"
  address: "localhost:5432"
```

Затем необходимо в переменной окружения под названием ```CONFIG_PATH``` указать
полный путь к файлу с вашей конфигурацией.

## Запуск сервера

#### 1. Необходимо склонировать репозиторий на локальную машину с помощью команды
```
git clone https://github.com/VladimirPopugaev/analytics_collector.git
```

#### 2. Запустить docker-compose с помощью команды

```
make services-up
```
Эта команда запустит сервисы с Postgres и сам сервер аналитики по адресу, 
прописанному в вашем конфиге

#### 3. Запустить миграции базы данных с помощью команды

```
make migrate-up
```

#### Готово, теперь ваш сервис доступен на локальной машине!

## Стек технологий
- Golang v.1.21.3
- Logging with __slog__
- __Postgres__ database
- Docker compose
