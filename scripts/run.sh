#!/bin/bash

echo "Обновление зависимостей"
go mod tidy

echo "Сборка приложения"
go build -o app cmd/main.go

echo "Запуск приложения"
./app &

# Сохраняем PID запущенного процесса
APP_PID=$!
echo "Приложение запущено с PID: $APP_PID"

# Ожидание доступности сервера
echo "Ожидаем, пока сервер будет готов..."
for i in {1..10}; do
  if curl -s http://localhost:8080 &> /dev/null; then
    echo "Сервер успешно запущен и готов к работе."
    exit 0
  fi
  echo "Попытка $i: сервер пока не готов..."
  sleep 2
done

echo "Ошибка: сервер не запустился."
exit 1