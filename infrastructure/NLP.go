package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"service/usecases"
	"strings"
	"time"

	"github.com/kljensen/snowball"
)

var ollamaURL = getOllamaURL()

func getOllamaURL() string {
	url := os.Getenv("OLLAMA_URL")
	if url == "" {
		url = "http://localhost:11434" // Значение по умолчанию для локального запуска
	}
	return url
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

type LocalNLPService struct {
	model string
}

func NewLocalNLPService() (*LocalNLPService, error) {
	modelName := os.Getenv("OLLAMA_MODEL")
	if modelName == "" {
		modelName = "mistral" // Дефолтная модель
	}

	// Проверяем, запущен ли Ollama
	if !isOllamaRunning() {
		return nil, errors.New("Ollama не запущен. Запустите его: 'ollama serve'")
	}

	// Проверяем, есть ли нужная модель
	if !isModelAvailable(modelName) {
		if err := pullModel(modelName); err != nil {
			return nil, err
		}
	}

	return &LocalNLPService{model: modelName}, nil
}

func (l *LocalNLPService) AnalyzeIntent(query string, db usecases.Database) (string, string, error) {
	log.Printf("Original query: %s", query)
	query = stemText(query)
	log.Printf("Stemmed query: %s", query)

	// Нормализуем запрос через стемминг
	stemmedQuery := stemText(query)

	// Сначала пытаемся найти по ключевым словам
	intent, response, err := db.FindIntentByKeywords(stemmedQuery)
	if err == nil {
		return intent, response, nil
	}

	// Получаем список всех интентов из базы
	intents, err := db.GetAllIntents()
	if err != nil {
		return "", "", err
	}

	// Формируем запрос в NLP
	prompt := fmt.Sprintf(
		"Из списка интентов: [%s] выбери наиболее близкий по смыслу к запросу: \"%s\". Ответь только названием интента.",
		strings.Join(intents, ", "), stemmedQuery,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	reqBody, _ := json.Marshal(OllamaRequest{
		Model:  l.model,
		Prompt: prompt,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", ollamaURL+"/api/generate", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var res OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", "", err
	}

	if res.Response == "" {
		return "", "", errors.New("Ollama вернула пустой ответ")
	}

	// Проверяем, есть ли такой интент в базе
	response, err = db.GetResponse(res.Response)
	if err != nil {
		return res.Response, "", nil // Интент есть, но ответа в БД нет
	}

	return res.Response, response, nil
}

func stemText(text string) string {
	words := strings.Fields(text) // Разбиваем текст на слова
	for i, word := range words {
		stemmedWord, _ := snowball.Stem(word, "russian", true) // Второе значение игнорируем
		words[i] = stemmedWord
	}
	return strings.Join(words, " ")
}

// --- Автоматизация ---

func isOllamaRunning() bool {
	resp, err := http.Get(ollamaURL + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func isModelAvailable(modelName string) bool {
	resp, err := http.Get(ollamaURL + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var data struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return false
	}

	for _, model := range data.Models {
		if model.Name == modelName {
			return true
		}
	}
	return false
}

func pullModel(modelName string) error {
	reqBody, _ := json.Marshal(map[string]string{"name": modelName})
	resp, err := http.Post(ollamaURL+"/api/pull", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("не удалось загрузить модель %s: %s", modelName, string(body))
	}
	return nil
}
