package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const CommandSet = "Список команд:\n" +
	"Регистрация нового пользователя: register\n" +
	"Аутентификация пользователя: login\n" +
	"Добавить новые данные: add\n" +
	"Получеие даннных пользователя: get\n" +
	"Удалить данные по id: delete\n" +
	"Изменить данные по id: edit\n" +
	"Информация о версии и дате сборки бинарного файла клиента: info"

const cookieName string = "user_info"

var token string

// UserData представляет структуру таблицы user_data.
type UserData struct {
	ID       uint
	UserID   uint
	UserData string
	DataType string
}

func main() {
	var name string
	for {
		fmt.Println("Введите команду. Для получения списка команда введите help")
		_, err := fmt.Scan(&name)

		if err != nil {
			return
		}

		switch name {
		case "add":
			addData()
		case "get":
			_, err = getData()
			if err != nil {
				return
			}
		case "edit":
			editData()
		case "delete":
			deleteData()
		case "register":
			register()
		case "login":
			authenticate()
		case "info":
			fmt.Println(CommandSet)
		case "ping":
			ping()
		default:
			fmt.Println("Вы ввели несуществующую команду.\n" + CommandSet)
		}
	}
}

func deleteData() {
	var dataID int
	fmt.Println("Введите id записи на удаление")

	for {
		_, err := fmt.Scan(&dataID)
		if err != nil {
			fmt.Println("Ошибка чтения id. Повторите попытку")
		} else {
			break
		}
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/delete_data?id="+strconv.Itoa(dataID), nil)

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie := &http.Cookie{
		Name:  cookieName,
		Value: token,
	}

	req.AddCookie(cookie)
	client := &http.Client{}
	resp, err := client.Do(req)

	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		log.Printf("Необходимо авторизоваться для работы с хранилищем!")
		return
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка при чтении ответа: %v", err)
	}

	fmt.Printf("Ответ сервера: %s\n", responseBody)
}

func editData() {
	userDataArray, err := getData()
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	var dataID uint
	fmt.Println("Введите id записи на изменение")

	for {
		_, err := fmt.Scan(&dataID)
		if err != nil {
			fmt.Println("Ошибка чтения id. Повторите попытку")
		} else {
			break
		}
	}

	metaInfo, err := enterMetaInfo()

	if err != nil {
		log.Fatalf("Ошибка ввода")
	}

	var jsonData []byte
	var data map[string]interface{}
	var dataType string
	for _, dataRow := range userDataArray {
		if dataRow.ID != dataID {
			continue
		}

		if dataRow.DataType == "log-pas" {
			data, err = enterKeyPas(metaInfo)
		} else if dataRow.DataType == "text" {
			data, err = enterText(metaInfo)
		} else if dataRow.DataType == "binary" {
			data, err = enterBinary(metaInfo)
		}

		dataType = dataRow.DataType
		break
	}

	if err != nil {
		log.Fatalf("Ошибка ввода")
	}

	jsonData, err = json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/edit_data", bytes.NewBuffer(jsonData))

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Data-Id", strconv.FormatUint(uint64(dataID), 10))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Data-Type", dataType)

	cookie := &http.Cookie{
		Name:  cookieName,
		Value: token,
	}

	req.AddCookie(cookie)
	client := &http.Client{}
	resp, err := client.Do(req)

	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		return
	}

	if resp.StatusCode == http.StatusUnauthorized {
		log.Printf("Необходимо авторизоваться для работы с хранилищем!")
		return
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка при чтении ответа: %v", err)
		return
	}

	fmt.Printf("Ответ сервера: %s\n", responseBody)
}

func getData() ([]UserData, error) {
	req, err := http.NewRequest("GET", "http://localhost:8080/get_data", nil)

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie := &http.Cookie{
		Name:  cookieName,
		Value: token,
	}

	req.AddCookie(cookie)
	client := &http.Client{}
	resp, err := client.Do(req)

	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		log.Printf("Необходимо авторизоваться для работы с хранилищем!")
		return nil, err
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка при чтении ответа: %v", err)
		return nil, err
	}

	var userDataArray []UserData
	err = json.Unmarshal(responseBody, &userDataArray)

	if err != nil {
		log.Printf("Ошибка при чтении ответа: %v", err)
		return nil, err
	}

	for _, userData := range userDataArray {
		fmt.Printf("ID: %d, UserID: %d, Data: %s, DataType: %s\n", userData.ID, userData.UserID, userData.UserData, userData.DataType)
	}

	return userDataArray, nil
}

func addData() {
	metaInfo, err := enterMetaInfo()
	if err != nil {
		return
	}

	var infoType string
	for {
		fmt.Println("Выберите тип сохраняемой информации\n" +
			"Логин-пароль: key-pas\n" +
			"Текстовые данные: text\n" +
			"Бинарные данные: binary\n" +
			"Выход: exit")
		_, err = fmt.Scan(&infoType)

		if err != nil {
			return
		}

		switch infoType {
		case "key-pas":
			err = addKeyPas(metaInfo)
			if err != nil {
				return
			}
			break
		case "text":
			err = addText(metaInfo)
			if err != nil {
				return
			}
			break
		case "binary":
			err = addBinary(metaInfo)
			if err != nil {
				return
			}
			break
		case "exit":
			return
		default:
			fmt.Println("Вы ввели несуществующую команду.\n" + CommandSet)
		}
	}
}

func enterMetaInfo() (string, error) {
	var answer string
	for {
		fmt.Println("Хотите ввести метаинформацию к данным? yes/no")
		_, err := fmt.Scanln(&answer)

		if err != nil {
			fmt.Println("Ошибка ввода")
			return "", err
		}

		if answer == "yes" || answer == "no" {
			break
		} else {
			fmt.Println("Введите 'да' или 'нет'")
		}
	}

	var metaInfo string

	if answer == "yes" {
		fmt.Println("Введите метаинформацию")
		reader := bufio.NewReader(os.Stdin)

		var err error
		metaInfo, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка ввода")
			return "", err
		}
	}
	metaInfo = strings.TrimSuffix(metaInfo, "\n")

	return metaInfo, nil
}

func addKeyPas(metaInfo string) error {
	data, err := enterKeyPas(metaInfo)

	if err != nil {
		log.Fatalf("Ошибка ввода")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/add_data", bytes.NewBuffer(jsonData))

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Data-Type", "log-pas")
	req.Header.Set("Meta-Info", metaInfo)
	cookie := &http.Cookie{
		Name:  cookieName,
		Value: token,
	}

	req.AddCookie(cookie)
	client := &http.Client{}
	resp, err := client.Do(req)

	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		log.Printf("Необходимо авторизоваться для работы с хранилищем")
		return errors.New("Необходимо авторизоваться для работы с хранилищем")
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ошибка при чтении ответа: %v", err)
	}

	fmt.Printf("Ответ сервера: %s\n", responseBody)

	return nil
}

func enterKeyPas(metaInfo string) (map[string]interface{}, error) {
	var login string
	var password string
	for {
		fmt.Println("Введите логин:")
		_, err := fmt.Scan(&login)

		if err != nil {
			return nil, err
		}

		if login == "" {
			fmt.Println("Введите непустую строку!")
		} else {
			break
		}
	}

	for {
		fmt.Println("Введите пароль:")
		_, err := fmt.Scan(&password)

		if err != nil {
			return nil, err
		}

		if login == "" {
			fmt.Println("Введите непустую строку!")
		} else {
			break
		}
	}

	data := map[string]interface{}{
		"login":     login,
		"password":  password,
		"meta-info": metaInfo,
	}

	return data, nil
}

func addText(metaInfo string) error {
	data, err := enterText(metaInfo)

	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/add_data", bytes.NewBuffer(jsonData))

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Data-Type", "text")

	cookie := &http.Cookie{
		Name:  cookieName,
		Value: token,
	}

	req.AddCookie(cookie)
	client := &http.Client{}
	resp, err := client.Do(req)

	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		log.Printf("Необходимо авторизоваться для работы с хранилищем!")
		return errors.New("необходимо авторизоваться для работы с хранилищем")
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка при чтении ответа: %v", err)
	}

	fmt.Printf("Ответ сервера: %s\n", responseBody)

	return nil
}

func enterText(metaInfo string) (map[string]interface{}, error) {
	var text string
	for {
		fmt.Println("Введите текст:")
		reader := bufio.NewReader(os.Stdin)

		var err error
		text, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка ввода")
			return nil, err
		}

		if text == "" {
			fmt.Println("Введите непустую строку!")
		} else {
			break
		}
	}
	text = strings.TrimSuffix(text, "\n")

	data := map[string]interface{}{
		"text":      text,
		"meta-info": metaInfo,
	}

	return data, nil
}

func addBinary(metaInfo string) error {
	dataInfo, err := enterBinary(metaInfo)

	if err != nil {
		log.Fatalf("Ошибка ввода")
	}

	jsonData, err := json.Marshal(dataInfo)
	if err != nil {
		log.Fatalf("Ошибка при сериализации JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/add_data", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Ошибка при создании запроса: %v", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Data-Type", "binary")

	cookie := &http.Cookie{
		Name:  cookieName,
		Value: token,
	}

	req.AddCookie(cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		return err
	}

	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	if resp.StatusCode == http.StatusUnauthorized {
		log.Printf("Необходимо авторизоваться для работы с хранилищем!")
		return errors.New("необходимо авторизоваться для работы с хранилищем")
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка при чтении ответа: %v", err)
		return err
	}

	fmt.Printf("Ответ сервера: %s\n", responseBody)

	return nil
}

func enterBinary(metaInfo string) (map[string]interface{}, error) {
	var data []byte

	fmt.Println("Введите бинарные данные (Ctrl+D для завершения):")
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения из стандартного ввода: %v", err)
	}

	dataInfo := map[string]interface{}{
		"meta-info":   metaInfo,
		"binary-data": data,
	}

	return dataInfo, nil
}

func ping() {
	resp, err := http.Get("http://localhost:8080/ping")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp.Status)
}

func register() {
	var login string
	var password string
	for {
		fmt.Println("Введите логин:")
		_, err := fmt.Scan(&login)

		if err != nil {
			return
		}

		if login == "" {
			fmt.Println("Введите непустую строку!")
		} else {
			break
		}
	}

	for {
		fmt.Println("Введите пароль:")
		_, err := fmt.Scan(&password)

		if err != nil {
			return
		}

		if login == "" {
			fmt.Println("Введите непустую строку!")
		} else {
			break
		}
	}

	data := map[string]string{
		"login":    login,
		"password": password,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	resp, err := http.Post(
		"http://localhost:8080/registration",
		"application/json",
		bytes.NewBuffer(jsonData))

	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка при чтении ответа: %v", err)
	}

	token, err = getCookieValue(resp.Header.Get("Set-Cookie"), cookieName)

	if err != nil {
		log.Printf("Ошибка при чтении ответа: %v", err)
	}

	fmt.Printf("Куки авторизации: %s\n", token)
	fmt.Printf("Ответ сервера: %s\n", responseBody)
}

// getCookieValue извлекает значение куки по её имени.
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

func authenticate() {
	var login string
	var password string
	for {
		fmt.Println("Введите логин:")
		_, err := fmt.Scan(&login)

		if err != nil {
			return
		}

		if login == "" {
			fmt.Println("Введите непустую строку!")
		} else {
			break
		}
	}

	for {
		fmt.Println("Введите пароль:")
		_, err := fmt.Scan(&password)

		if err != nil {
			return
		}

		if login == "" {
			fmt.Println("Введите непустую строку!")
		} else {
			break
		}
	}

	data := map[string]string{
		"login":    login,
		"password": password,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	resp, err := http.Post(
		"http://localhost:8080/authentication",
		"application/json",
		bytes.NewBuffer(jsonData))

	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка при чтении ответа: %v", err)
	}

	token, err = getCookieValue(resp.Header.Get("Set-Cookie"), cookieName)

	if err != nil {
		log.Printf("Ошибка при чтении ответа: %v", err)
	}

	fmt.Printf("Куки авторизации: %s\n", token)
	fmt.Printf("Ответ сервера: %s\n", responseBody)
}
