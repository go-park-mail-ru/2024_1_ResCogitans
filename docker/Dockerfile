# Используем базовый образ PostgreSQL
FROM postgres:13.2

# Переменные окружения для настройки базы данных
ENV POSTGRES_DB=jantugan
ENV POSTGRES_USER=mrdzhofik
ENV POSTGRES_PASSWORD=246858

# Копируем init.sql в контейнер
COPY init.sql /docker-entrypoint-initdb.d/

# По умолчанию PostgreSQL слушает порт 5432
EXPOSE 5432

# Команда для запуска PostgreSQL сервера
CMD ["postgres"]
