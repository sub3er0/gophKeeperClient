// commands/commands.go
package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"gophKeeperClient/api"
	"gophKeeperClient/cli"
	"log"
	"strconv"
)

type CommandHandler struct {
	APIClient *api.APIClient
	CLIHelper *cli.CLIHelper
}

// UserData представляет структуру таблицы user_data.
type UserData struct {
	ID       uint
	UserID   uint
	UserData string
	DataType string
}

const (
	CommandSet = "Список команд:\n" +
		"Регистрация нового пользователя: register\n" +
		"Аутентификация пользователя: login\n" +
		"Добавить новые данные: add\n" +
		"Получеие даннных пользователя: get\n" +
		"Удалить данные по id: delete\n" +
		"Изменить данные по id: edit\n" +
		"Информация о версии и дате сборки бинарного файла клиента: info"
)

func (handler *CommandHandler) Run(ctx context.Context) {
	var command string
	for {
		fmt.Println("Введите команду:")

		select {
		case <-ctx.Done():
			fmt.Println("Программа завершена")
			return
		default:
			_, err := fmt.Scan(&command)
			if err != nil {
				log.Println("Ошибка ввода: ", err)
				continue
			}

			if command == "exit" {
				fmt.Println("Завершение работы...")
				return
			}

			switch command {
			case "add":
				handler.AddData()
			case "get":
				_, err = handler.GetData()
				if err != nil {
					return
				}
			case "edit":
				handler.EditData()
			case "delete":
				handler.DeleteData()
			case "login":
				handler.Authenticate()
			case "register":
				handler.Register()
			case "ping":
				handler.Ping()
			default:
				fmt.Println("Неизвестная команда")
			}
		}
	}
}

func (handler *CommandHandler) AddData() {
	metaInfo, err := handler.CLIHelper.GetMetaInfo()
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
			err = handler.addKeyPas(handler.APIClient, metaInfo)
			if err != nil {
				return
			}
			break
		case "text":
			err = handler.addText(handler.APIClient, metaInfo)
			if err != nil {
				return
			}
			break
		case "binary":
			err = handler.addBinary(handler.APIClient, metaInfo)
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

func (handler *CommandHandler) addKeyPas(apiClient *api.APIClient, metaInfo string) error {
	data, err := handler.CLIHelper.EnterKeyPas(metaInfo)
	if err != nil {
		log.Printf("Ошибка ввода")
		return err
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"Data-Type":    "log-pas",
		"Meta-Info":    metaInfo,
	}

	responseBody, err := apiClient.Post("add_data", data, headers)
	if err != nil {
		log.Printf("Ошибка при добавлении логин пароль: %v", err)
	}

	fmt.Printf("Ответ сервера: %s\n", responseBody)

	return nil
}

func (handler *CommandHandler) addText(apiClient *api.APIClient, metaInfo string) error {
	data, err := handler.CLIHelper.EnterText(metaInfo)

	if err != nil {
		return err
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"Data-Type":    "text",
	}

	responseBody, err := apiClient.Post("add_data", data, headers)
	if err != nil {
		log.Printf("Ошибка при добавлении текста: %v", err)
	}

	fmt.Printf("Ответ сервера: %s\n", responseBody)

	return nil
}

func (handler *CommandHandler) addBinary(apiClient *api.APIClient, metaInfo string) error {
	dataInfo, err := handler.CLIHelper.EnterBinary(metaInfo)

	if err != nil {
		log.Printf("Ошибка ввода")
		return err
	}

	headers := map[string]string{
		"Content-Type": "application/octet-stream",
		"Data-Type":    "binary",
	}

	responseBody, err := apiClient.Post("add_data", dataInfo, headers)
	if err != nil {
		log.Printf("Ошибка при добавлении бинарных данных: %v", err)
	}

	fmt.Printf("Ответ сервера: %s\n", responseBody)

	return nil
}

func (handler *CommandHandler) GetData() ([]UserData, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	responseBody, err := handler.APIClient.Get("get_data", headers)
	if err != nil {
		log.Printf("Ошибка при получении спика данных: %v", err)
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

func (handler *CommandHandler) EditData() {
	userDataArray, err := handler.GetData()
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	var dataID uint
	handler.CLIHelper.EnterDataID(&dataID)

	metaInfo, err := handler.CLIHelper.GetMetaInfo()

	if err != nil {
		log.Printf("Ошибка ввода")
		return
	}

	var data map[string]interface{}
	var dataType string
	for _, dataRow := range userDataArray {
		if dataRow.ID != dataID {
			continue
		}

		if dataRow.DataType == "log-pas" {
			data, err = handler.CLIHelper.EnterKeyPas(metaInfo)
		} else if dataRow.DataType == "text" {
			data, err = handler.CLIHelper.EnterText(metaInfo)
		} else if dataRow.DataType == "binary" {
			data, err = handler.CLIHelper.EnterBinary(metaInfo)
		}

		dataType = dataRow.DataType
		break
	}

	if err != nil {
		log.Printf("Ошибка ввода")
		return
	}

	headers := map[string]string{
		"Data-Id":      strconv.FormatUint(uint64(dataID), 10),
		"Content-Type": "application/json",
		"Data-Type":    dataType,
	}

	responseBody, err := handler.APIClient.Post("edit_data", data, headers)
	if err != nil {
		log.Printf("Ошибка при изменении данных: %v", err)
		return
	}

	fmt.Printf("Ответ сервера: %s\n", responseBody)
}

func (handler *CommandHandler) DeleteData() {
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

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	responseBody, err := handler.APIClient.Get("delete_data?id="+strconv.Itoa(dataID), headers)
	if err != nil {
		log.Printf("Ошибка при удалении записи: %v", err)
	}

	fmt.Printf("Ответ сервера: %s\n", responseBody)
}

func (handler *CommandHandler) Authenticate() {
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

	token, responseBody, err := handler.APIClient.Authenticate(login, password)
	handler.APIClient.SetToken(token)

	if err != nil {
		log.Printf("Ошибка при чтении ответа: %v", err)
	}

	fmt.Printf("Куки авторизации: %s\n", token)
	fmt.Printf("Ответ сервера: %s\n", responseBody)
}

func (handler *CommandHandler) Register() {
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

	token, responseBody, err := handler.APIClient.Registration(login, password)
	handler.APIClient.SetToken(token)

	if err != nil {
		log.Printf("Ошибка при чтении ответа: %v", err)
	}

	fmt.Printf("Куки авторизации: %s\n", token)
	fmt.Printf("Ответ сервера: %s\n", responseBody)
}

func (handler *CommandHandler) Ping() {
	handler.APIClient.Ping()
}
