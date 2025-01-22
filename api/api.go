package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const CookieName string = "user_info"

type APIClient struct {
	BaseURL string
	Token   string
}

// APIClientInterface интерфейс для APIClient
type APIClientInterface interface {
	Authenticate(login, password string) (string, []byte, error)
	Registration(login, password string) (string, []byte, error)
	Get(endpoint string, headers map[string]string) ([]byte, error)
	Post(endpoint string, data interface{}, headers map[string]string) ([]byte, error)
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{BaseURL: baseURL}
}

func (client *APIClient) SetToken(token string) {
	client.Token = token
}

func (client *APIClient) Get(endpoint string, headers map[string]string) ([]byte, error) {
	var err error

	url := fmt.Sprintf("%s/%s", client.BaseURL, endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}

	cookie := &http.Cookie{
		Name:  CookieName,
		Value: client.Token,
	}
	req.AddCookie(cookie)

	// Установка пользовательских заголовков
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	clientHTTP := &http.Client{}
	resp, err := clientHTTP.Do(req)
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		return nil, err
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка при чтении ответа: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Ошибка: статус код %d, ответ: %s", resp.StatusCode, responseBody)
		return nil, fmt.Errorf("Ошибка: статус код %d", resp.StatusCode)
	}

	return responseBody, nil
}

func (client *APIClient) Post(endpoint string, data interface{}, headers map[string]string) ([]byte, error) {
	var jsonData []byte
	var err error

	if data != nil {
		jsonData, err = json.Marshal(data)
		if err != nil {
			log.Printf("Error marshaling JSON: %v", err)
			return nil, err
		}
	}

	url := fmt.Sprintf("%s/%s", client.BaseURL, endpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}

	cookie := &http.Cookie{
		Name:  CookieName,
		Value: client.Token,
	}
	req.AddCookie(cookie)

	// Установка пользовательских заголовков
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	clientHTTP := &http.Client{}
	resp, err := clientHTTP.Do(req)
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		return nil, err
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка при чтении ответа: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Ошибка: статус код %d, ответ: %s", resp.StatusCode, responseBody)
		return nil, fmt.Errorf("Ошибка: статус код %d", resp.StatusCode)
	}

	return responseBody, nil
}

// Authenticate Метод для аутентификации пользователя
func (client *APIClient) Authenticate(login, password string) (string, []byte, error) {
	data := map[string]string{
		"login":    login,
		"password": password,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", nil, fmt.Errorf("ошибка сериализации JSON: %v", err)
	}

	resp, err := http.Post(client.BaseURL+"/authentication", "application/json", bytes.NewBuffer(jsonData))
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	if err != nil {
		return "", nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	// Извлекаем куку
	token, err := getCookieValue(resp.Header.Get("Set-Cookie"), CookieName)
	if err != nil {
		return "", nil, fmt.Errorf("ошибка при извлечении куки: %v", err)
	}

	return token, responseBody, nil
}

// Registration Метод для регистрации пользователя
func (client *APIClient) Registration(login, password string) (string, []byte, error) {
	data := map[string]string{
		"login":    login,
		"password": password,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", nil, fmt.Errorf("ошибка сериализации JSON: %v", err)
	}

	resp, err := http.Post(client.BaseURL+"/registration", "application/json", bytes.NewBuffer(jsonData))
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	if err != nil {
		return "", nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	// Извлекаем куку
	token, err := getCookieValue(resp.Header.Get("Set-Cookie"), CookieName)
	if err != nil {
		return "", nil, fmt.Errorf("ошибка при извлечении куки: %v", err)
	}

	return token, responseBody, nil
}

func (client *APIClient) Ping() {
	resp, err := http.Get(client.BaseURL + "/ping")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.Status)
}

// Функция для извлечения значения куки
func getCookieValue(cookieStr string, cookieName string) (string, error) {
	parts := strings.Split(cookieStr, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, cookieName+"=") {
			value := strings.SplitN(part, "=", 2)
			if len(value) == 2 {
				return value[1], nil
			}
		}
	}
	return "", fmt.Errorf("кука с именем %s не найдена", cookieName)
}
