# MWS CLI

MWS — утилита для управления профилями.


## Быстрый старт

### Linux / MacOS
```bash
go build -o mws .
```
```bash
./mws profile create --name=dev --user=dev-account --project=delivery-dev
./mws profile list
./mws profile get --name dev
```

### Windows
```bash
go build -o mws.exe .
```

```bash
./mws.exe profile create --name=dev --user=dev-account --project=delivery-dev
./mws.exe profile list
./mws.exe profile get --name=dev
```

## Команды

- `./mws profile create --name=<name> --user=<user> --project=<project>`
- `./mws profile get --name=<name> [--output json]`
- `./mws profile list [--output json]`
- `./mws profile delete --name=<name>`
- `./mws version`
- `./mws help`

## Хранение профилей

Профили хранятся в папке `profiles` в текущей директории:

```
profiles/
  dev.yaml
  prod.yaml
```

Каждый профиль — это YAML-файл:

```yaml
user: dev-account
project: delivery-dev
```

## Примеры

Создать профиль:

```bash
./mws profile create --name=prod --user=prod-account --project=delivery-prod

Profile "prod" created.
```

Получить профиль:

```bash
./mws profile get --name=prod

+------+--------------+---------------+
| NAME | USER         | PROJECT       |
+------+--------------+---------------+
| prod | prod-account | delivery-prod |
+------+--------------+---------------+
```

Посмотреть список:

```bash
./mws profile list

+------+--------------+---------------+
| NAME | USER         | PROJECT       |
+------+--------------+---------------+
| dev  | dev-account  | delivery-dev  |
+------+--------------+---------------+
| prod | prod-account | delivery-prod |
+------+--------------+---------------+
```

Удалить профиль:

```bash
./mws profile delete --name=prod

Profile "prod" deleted.
```

Справка:

```bash
./mws profile --help

Manage profiles.

Usage:
  mws profile <command> [flags]

Available Commands:
  create      Create a profile
  get         Show profile details
  list        List profiles
  delete      Delete a profile

Examples:
  mws profile create --name dev --user dev-account --project delivery-dev
  mws profile get dev
  mws profile list
  mws profile delete dev

Use "mws profile <command> --help" for command-specific help.
```
