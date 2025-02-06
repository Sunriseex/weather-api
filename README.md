# Weather API на Go

Это проект Weather API, написанный на Go, который получает данные о погоде из стороннего API (Visual Crossing), кэширует результаты в Redis и логирует ключевые события. Проект демонстрирует интеграцию с внешними API, использование in-memory кэша, работу с конфигурационными файлами и расширенное логирование.

## Функциональность

Получение данных о погоде: API обращается к Visual Crossing для получения актуальной информации о погоде по запрашиваемому городу.

Кэширование: Результаты запросов к погодному API сохраняются в Redis с установленным временем истечения (например, 12 часов), что снижает нагрузку на внешний API.

Конфигурация: Важные настройки (API-ключ, URL Redis, время кэширования) вынесены в файл config.json.

Логирование: Все ключевые действия и ошибки логируются с использованием настроенного логгера. Логи записываются в файл app.log.



## Структура проекта

weather-api/
├── config.json    # Файл конфигурации с API-ключом, настройками Redis и временем кэширования
├── main.go        # Основной код приложения
├── app.log        # Файл для записи логов (будет создан автоматически)
└── README.md      # Описание проекта

## Требования

Go (рекомендуется версия 1.16 и выше)

Redis

Локально установленный сервер Redis или контейнер Docker.

Internet-доступ для обращения к Visual Crossing API.

## Установка и запуск

1. Клонируйте репозиторий
```sh
git clone https://github.com/Sunriseex/weather-api
cd weather-api
```
2. Установите зависимости
```go
go mod tidy
```
3. Настройте config.json

Создайте файл config.json (если его нет) и укажите API-ключ, данные для подключения к Redis и время кэширования:
```json
{
  "apiKey": "YOUR_VISUAL_CROSSING_API_KEY",
  "redisAddr": "localhost:6379",
  "cacheExpiration": 43200
}
```
4. Запустите сервер
```go
go run main.go
```
5. Использование API

Для получения данных о погоде выполните запрос:
```sh
    curl http://localhost:8080/weather/London
```

Пример ответа:
```json
    {
      "location": "London",
     "temperature": "15.2",
     "condition": "Partly Cloudy"
    }
```
6. Логирование

Логи записываются в app.log. В случае ошибок можно проверить содержимое файла:
```sh
tail -f app.log
```

Этот проект является учебным и демонстрирует основы работы с внешними API, кэшированием и логированием в Go.
