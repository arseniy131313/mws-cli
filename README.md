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
./mws.exe profile get --name dev
```

## Команды

- `mws profile create --name=<name> --user=<user> --project=<project>`
- `mws profile get <name> [--output json]`
- `mws profile list [--output json]`
- `mws profile delete--name=<name>`
- `mws version`

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
```

Получить профиль:

```bash
./mws profile get --name=prod
```

Посмотреть список:

```bash
./mws profile list
```

Удалить профиль:

```bash
./mws profile delete --name=prod
```

Справка:

```bash
./mws profile --help
```