```markdown
# GophKeeperClient

GophKeeperClient - это клиентское приложение для взаимодействия с GophKeeper, системой хранения и управления данными пользователей, обеспечивающее безопасную аутентификацию и управление данными.

## Установка

Для установки проекта вам потребуются следующие зависимости:

- Go (версии 1.16 или выше)
- PostgreSQL (для работы с базой данных)

Следуйте следующим шагам:

1. Клонируйте репозиторий:

```bash
git clone https://github.com/yourusername/gophKeeperClient.git
cd gophKeeperClient
```

2. Установите зависимости:

```bash
go mod tidy
```

3. Настройте конфигурацию, изменив файл конфигурации, например, `config.yaml`:

```yaml
ServerAddress: ":8080"
BaseURL: "http://localhost:8080"
DatabaseDsn: "host=localhost user=yourusername dbname=gophkeeper sslmode=disable" # обновите значения
```

## Использование

Запустите клиент, используя следующую команду:

```bash
go run main.go
```

После запуска программа предложит вам ввести команду. Доступные команды:

- `register` - Регистрация нового пользователя.
- `login` - Аутентификация пользователя.
- `add` - Добавить новые данные.
- `get` - Получение данных пользователя.
- `delete` - Удалить данные по идентификатору (id).
- `edit` - Изменить данные по идентификатору (id).
- `info` - Информация о версии и дате сборки клиентского бинарного файла.
- `ping` - Проверка связи с сервером.
- `exit` - Завершение работы приложения.

### Примеры команд

#### Регистрация нового пользователя
```bash
Введите команду: register
Введите логин: testuser
Введите пароль: yourpassword
```

#### Аутентификация
```bash
Введите команду: login
Введите логин: testuser
Введите пароль: yourpassword
```

#### Добавление данных
```bash
Введите команду: add
Введите метаинформацию: Example meta-info
Выберите тип данных: key-pas
Введите логин: testuser
Введите пароль: yourpassword
```

## API

### `/add_data`

**Метод**: POST
**Описание**: Добавляет новые данные пользователя.

**Тело запроса**:
```json
{
    "login": "example_user",
    "password": "example_password"
}
```

### `/get_data`

**Метод**: GET
**Описание**: Получает данные пользователя.

### `/delete_data`

**Метод**: GET
**Описание**: Удаляет данные пользователя по идентификатору (id).

**Параметры**:
- `id`: идентификатор записи данных.

### `/edit_data`

**Метод**: POST
**Описание**: Изменяет данные пользователя по идентификатору (id).

### `/registration`

**Метод**: POST
**Описание**: Регистрирует нового пользователя.

**Тело запроса**:
```json
{
    "login": "example_user",
    "password": "example_password"
}
```

### `/authentication`

**Метод**: POST
**Описание**: Аутентифицирует существующего пользователя.

**Тело запроса**:
```json
{
    "login": "example_user",
    "password": "example_password"
}
```

### `/ping`

**Метод**: GET
**Описание**: Проверяет соединение с сервером.

```
