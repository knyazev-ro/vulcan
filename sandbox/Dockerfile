FROM postgres:latest

# Устанавливаем переменные окружения для БД
ENV POSTGRES_DB=vulcan_test
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=password

# Копируем скрипт инициализации. 
# Скрипты в этой папке выполняются только при первом запуске контейнера.
COPY init.sql /docker-entrypoint-initdb.d/