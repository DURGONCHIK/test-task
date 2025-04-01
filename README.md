## 🛠 Установка и запуск  

### 1️⃣ **Клонирование репозитория**  
```sh
git clone https://github.com/DURGONCHIK/test-task.git
cd test-task


настройка файла окружени
2. cp .env.example .env

внутри указать содержимое

PORT=8080
POSTGRES_CONN=postgres://user:password@db:5432/user?sslmode=disable
OLLAMA_URL=http://ollama:11434
OLLAMA_MODEL=tinydolphin:latest


3. Производить запуск через докер

docker-compose up --build -d


4. Далее установка NLP модели

docker exec -it ollama bash

 ввести внутри контейнера для установки нлп
ollama pull tinydolphin:latest

выходим из контейнера

exit


тестировать рекомендую через постман

указываем метод POST и вводим в строку http://localhost:8080/query

во вкладке Headers в поле key - Content-Type ,  в поле value - application/json

во вкладке body пишем тело запроса, например

{
  "text": "когда мне ждать мою доставку?"
}

в бд записаны примерные ключевые слова клиента и примерные ответы на них
рекомендую использовать запросы с ключевыми словами "доставка" "гарантия" "жалоба" "скидка" как пример работы сервиса
