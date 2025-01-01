-- Создаем пользователя
CREATE USER validator WITH PASSWORD 'val1dat0r';

-- Создаем базу данных
CREATE DATABASE "project-sem-1"
    WITH 
    OWNER = validator
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TEMPLATE = template0;

-- Предоставляем пользователю полные права на базу данных
GRANT ALL PRIVILEGES ON DATABASE "project-sem-1" TO validator;