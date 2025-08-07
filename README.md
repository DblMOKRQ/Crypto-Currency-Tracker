# Crypto-Currency-Tracker

#### Описание проекта

Этот микросервис предназначен для сбора, хранения и отображения стоимости криптовалют. Сервис позволяет:

* Добавлять криптовалюты в список наблюдения
* Удалять криптовалюты из списка наблюдения
* Получать цену криптовалюты в конкретный момент времени

Сервис автоматически обновляет цены отслеживаемых криптовалют через заданные интервалы времени, используя CoinGecko API.

##### Технологии:

* Язык: Go 1.24
* База данных: PostgreSQL
* Логирование: Zap
* Конфигурация: YAML
* Контейнеризация: Docker Compose
* API: HTTP REST

#### Запуск проекта
###### Требования

* Установленный Docker и Docker Compose
* API ключ от CoinGecko (можно получить на [CoinGecko](https://www.coingecko.com/))

###### Инструкция по запуску
1. Клонируйте репозиторий:
```bash
git clone github.com/DblMOKRQ/Crypto-Currency-Tracker
cd Crypto-Currency-Tracker
```
2. Создайте файл конфигурации config/config.yaml (пример содержимого в разделе ниже)
3. Запустите сервис с помощью Docker Compose:

```bash
docker-compose up --build
```
4. Сервис будет доступен по адресу: http://localhost:8080

##### Конфигурация

Пример содержимого config/config.yaml:
```yaml
# Конфигурация хранилища (PostgreSQL)
storage:
  user: "postgres"          # Имя пользователя PostgreSQL
  password: "123"           # Пароль пользователя PostgreSQL
  host: "postgres"          # Хост базы данных (в Docker Compose - имя сервиса)
  port: "5432"              # Порт базы данных
  db_name: "coins"          # Название базы данных
  ssl_mode: "disable"       # Режим SSL (disable/prefer/require)

# Настройки обновления цен
price_updates: 10s          # Интервал обновления цен (например: 10s, 1m)

# Настройки CoinGecko API
api_key: "ваш_api_ключ"     # API ключ для доступа к CoinGecko API
vs_currency: "usd"          # Валюта, в которой отображаются цены (usd, eur, rub и т.д.)

# Настройки REST сервера
rest:
  address: ":8080"          # Адрес и порт, на котором будет работать сервер
```
1. **price_updates** - поддерживает значения в формате:
    - `10s` - 10 секунд
    - `1m` - 1 минута
    - `1h` - 1 час
2. **ssl_mode** - рекомендуется использовать:
    - `disable` - для локальной разработки
    - `require` - для production
3. **vs_currency** - поддерживает все валюты, доступные в CoinGecko API:
    - Основные: `usd`, `eur`, `gbp`, `jpy`
    - Криптовалюты: `btc`, `eth`
    - Другие: `rub`, `cny`, и т.д.
4. **api_key** - можно получить бесплатный ключ на [CoinGecko API](https://www.coingecko.com/en/api)
### Переменные окружения

Для корректной работы необходимо установить следующие переменные окружения:

- `CONFIG_PATH` - путь к файлу конфигурации (по умолчанию `/app/config/config.yaml`)

- `MIGRATIONS_DIR` - путь к папке с миграциями (по умолчанию `/app/migrations`)

#### ## API Endpoints

`POST /currency/add` - Добавление криптовалюты в список наблюдения
```json
{
  "coin": "BTC"
}
```
`POST /currency/remove` - Удаление криптовалюты из списка наблюдения
```json
{
  "coin": "BTC"
}
```
`POST /currency/get` - Получение цены криптовалюты
```json
{
  "coin": "BTC",
  "timestamp": 1736500490
}
```
