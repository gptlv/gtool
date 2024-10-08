# gtool

[🇬🇧 English version](README.md)

Данное приложение предназначено для взаимодействия с системами JIRA и Active Directory и предлагает ряд инструментов для автоматизации задач.

## Зачем?

В ходе работы я начал замечать повторяющиеся, рутинные задачи, которые выполнялись в несколько простых действий. Родилась идея автоматизировать выполнение этих задач.

## Основной функционал

- **Работа с задачами JIRA**: обновление статуса сотрудников в CMDB Insight, назначение, обновление и закрытие задач на блокировку/прием сотрудников, предоставление доступов.
- **Работа с объектами в CMDB**: генерация информации для списания оборудования, вывод описания ноутбука для последующего оформления заказов в курьерской службе.
- **Работа с LDAP**: управление группами пользователей и их доступами.

## Установка

1. Склонировать репозиторий

```
git clone https://github.com/gptlv/gtool.git
cd gtool
```

2. Создать файл `.env` и указать необходимые данные (пример -- в файле `.env.example`)

```
touch .env
```

3. Заполнить конфигурационный файл `config.yml`. Необходимо указать актуальные фамилии сотрудников для генерации документов

4. Выполнить команду `go mod tidy` для установки необходимых библиотек

5. Собрать проект

```
go build -o gtool .
```

## Доступные команды

### Обработка заявок в JIRA

- `./gtool issue process-insight` Обработка заявок на деактивацию пользователя в Insight
- `./gtool issue process-ldap` Обработка заявок на деактивацию пользователей в Active Directory
- `./gtool issue process-staff --component=[all|hiring|dismissal]` Обработка заявок на прием и/или увольнение сотрудников
- `./gtool issue grant-access --key=<key>` Обработка заявки на выдачу доступа (добавление группы в Active Directory)
- `./gtool issue assign --component=[all|hiring|dismissal|insight|ldap]` Назначение заявок на текущуго пользователя
- `./gtool issue update-trainee --key=<key>` Обновление названия подзадач задачи на блокировку стажеров
- `./gtool issue show-empty` Вывод заявок, которые требуют указания компонента

### Обращения к CMDB

- `./gtool asset generate-records --start=<id>` Генерация `.csv` файла с необходимыми данными для использования в скрипте [wroffs](https://github.com/gptlv/wroffs)
- `./gtool asset print-description --isc=<isc>` Вывод описания ноутбука для добавления в заказ для курьерской службы

### Обращения к Active Directory

- `./gtool ldap add-group -emails <user1@ex.com user2@ex.com> -cns <cn1 cn2 cn3>` Добавление нескольких пользователей в несколько групп в Active Directory

### Основные библиотеки

- [go-arg](https://github.com/alexflint/go-arg) для парсинга аргументов командной строки.
- [go-jira](https://github.com/andygrunwald/go-jira) для взаимодействия с API JIRA. Использовал [собственный форк](https://github.com/gptlv/go-jira), который расширяет функционал для работы с JIRA Insight
- [go-ldap](https://github.com/go-ldap/ldap) для операций с LDAP.
