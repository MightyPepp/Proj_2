#!/bin/bash

echo "=== Тестирование Task API ==="
echo

# Проверка сервера
echo "1. Проверяем сервер..."
curl -s http://localhost:8080/
echo

# Создание задач
echo "2. Создаем задачи..."
echo "Задача 1:"
curl -s -X POST http://localhost:8080/task/ \
  -H "Content-Type: application/json" \
  -d '{"text": "Купить молоко", "tags": ["shopping", "home"], "due": "2024-12-25T10:00:00Z"}'
echo

echo "Задача 2:"
curl -s -X POST http://localhost:8080/task/ \
  -H "Content-Type: application/json" \
  -d '{"text": "Изучить Go", "tags": ["programming", "learning"], "due": "2024-12-26T14:00:00Z"}'
echo

# Все задачи
echo "3. Все задачи:"
curl -s http://localhost:8080/task/ 
echo

# Поиск по тегу
echo "4. Задачи с тегом 'shopping':"
curl -s http://localhost:8080/tag/shopping/ |
echo

# Удаление
echo "5. Удаляем все задачи..."
curl -s -X DELETE http://localhost:8080/task/
echo

echo "6. Проверяем удаление:"
curl -s http://localhost:8080/task/ |

echo "=== Тест завершен ==="