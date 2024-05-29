# 2024_1_ResCogitans
Бэкенд проекта "TripAdvisor" команды "res cogitans"

## Авторы

[Спасскова Ксения](https://github.com/sp20ks)

[Джафаров Эмиль](https://github.com/MrDzhofik)

[Саяпин Егор](https://github.com/exg0rd)

[Горевой Александр](https://github.com/AlGrItm/)

## Менторы

[Иван Червоный](https://github.com/mzingwrld) - _Frontend_

[Седойкин Георгий](https://github.com/GeorgiyX) - _Backend_

TODO - _UX_

## Ссылки
[Бэкенд проекта](https://github.com/go-park-mail-ru/2024_1_ResCogitans)

# Сборка Фронта

## 1. Локальная разработка
* npm i
* npm run start

## 2. Сборка
* npm i
* npm run build

## 3. Форматирование 
* npm run lint:fix - прогнать eslint с правкой и форматированием кода под кодстайл AirBnB

# Сборка Бэка

## 1. Создание .env
* файл должен содержать абсолютный путь до файла конфигурации yaml
* файл должен лежать рядом с main.go

## 2. Поднятие баз данных через docker

### Postgres
`docker build -t cogitans .`<br>
`docker run -d --name <name> -p 5433:5432 <name>`

### Redis
`docker run -d --name <name> -p 6379:6379 redis redis-server --requirepass <password>`

## 3. Запуск
* go run .
## 4. Деплой
* [Сайт](http://jantugan.ru)
Тестовый юзер
{
"username": "san@boy.ru",
"password": "ABC123abc123!"
}
