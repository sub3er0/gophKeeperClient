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
	APIClient api.APIClientInterface
	CLIHelper cli.CLIHelperInterface
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
	var err error
	for {
		fmt.Println("Введите команду:")

		select {
		case <-ctx.Done():
			fmt.Println("Программа завершена")
			return
		default:
			_, err = fmt.Scan(&command)
			if err != nil {
				log.Println("Ошибка ввода: ", err)
				continue
			}

			if command == "exit" {
				fmt.Println("Завершение работы...")
				return
			}

			switch command {
			case "help":
				fmt.Println(CommandSet)
			case "add":
				err = handler.AddData()
			case "get":
				_, err = handler.GetData()
			case "edit":
				err = handler.EditData()
			case "delete":
				err = handler.DeleteData()
			case "login":
				err = handler.Authenticate()
				if err != nil {
					return
				}
			case "register":
				err = handler.Register()
			case "ping":
				err = handler.Ping()
			default:
				fmt.Println("Неизвестная команда")
			}
		}

		if err != nil {
			log.Println(err)
		}
	}
}

func (handler *CommandHandler) AddData() error {
	metaInfo, err := handler.CLIHelper.GetMetaInfo()
	if err != nil {
		return err
	}

	var infoType string
	for {
		infoType, err = handler.CLIHelper.EnterInfoType()

		if err != nil {
			return err
		}

		switch infoType {
		case "key-pas":
			err = handler.addKeyPas(handler.APIClient, metaInfo)
			if err != nil {
				return err
			}
			break
		case "text":
			err = handler.addText(handler.APIClient, metaInfo)
			if err != nil {
				return err
			}
			break
		case "binary":
			err = handler.addBinary(handler.APIClient, metaInfo)
			if err != nil {
				return err
			}
			break
		case "exit":
			return err
		default:
			fmt.Println("Вы ввели несуществующую команду.\n" + CommandSet)
		}

		return nil
	}
}

func (handler *CommandHandler) addKeyPas(apiClient api.APIClientInterface, metaInfo string) error {
	data, err := handler.CLIHelper.EnterKeyPas(metaInfo)
	if err != nil {
		return fmt.Errorf("ошибка ввода: %v", err)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"Data-Type":    "log-pas",
		"Meta-Info":    metaInfo,
	}

	responseBody, err := apiClient.Post("add_data", data, headers)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении логин пароль: %v", err)
	}

	fmt.Printf("Ответ сервера: %s\n", responseBody)

	return nil
}

func (handler *CommandHandler) addText(apiClient api.APIClientInterface, metaInfo string) error {
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
		return fmt.Errorf("ошибка при добавлении текста: %v", err)

	}

	fmt.Printf("Ответ сервера: %s\n", responseBody)

	return nil
}

func (handler *CommandHandler) addBinary(apiClient api.APIClientInterface, metaInfo string) error {
	dataInfo, err := handler.CLIHelper.EnterBinary(metaInfo)

	if err != nil {
		return fmt.Errorf("ошибка ввода: %v", err)
	}

	headers := map[string]string{
		"Content-Type": "application/octet-stream",
		"Data-Type":    "binary",
	}

	responseBody, err := apiClient.Post("add_data", dataInfo, headers)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении бинарных данных: %v", err)
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
		return nil, fmt.Errorf("ошибка при получении спика данных: %v", err)
	}

	var userDataArray []UserData
	err = json.Unmarshal(responseBody, &userDataArray)

	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	for _, userData := range userDataArray {
		fmt.Printf("ID: %d, UserID: %d, Data: %s, DataType: %s\n", userData.ID, userData.UserID, userData.UserData, userData.DataType)
	}

	return userDataArray, nil
}

func (handler *CommandHandler) EditData() error {
	userDataArray, err := handler.GetData()
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return err
	}

	fmt.Println("Введите id записи на изменение")
	dataID, err := handler.CLIHelper.EnterDataID()
	if err != nil {
		log.Print(err)
		return err
	}

	metaInfo, err := handler.CLIHelper.GetMetaInfo()

	if err != nil {
		return fmt.Errorf("ошибка ввода: %v", err)
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
		return fmt.Errorf("ошибка ввода: %v", err)
	}

	headers := map[string]string{
		"Data-Id":      strconv.FormatUint(uint64(dataID), 10),
		"Content-Type": "application/json",
		"Data-Type":    dataType,
	}

	responseBody, err := handler.APIClient.Post("edit_data", data, headers)
	if err != nil {
		return fmt.Errorf("ошибка при изменении данных: %v", err)
	}

	fmt.Printf("Ответ сервера: %s\n", responseBody)

	return nil
}

func (handler *CommandHandler) DeleteData() error {
	fmt.Println("Введите id записи на удаление")
	dataID, err := handler.CLIHelper.EnterDataID()
	if err != nil {
		return err
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	responseBody, err := handler.APIClient.Get("delete_data?id="+strconv.FormatUint(uint64(dataID), 10), headers)
	if err != nil {
		return fmt.Errorf("ошибка при удалении записи: %v", err)
	}

	fmt.Printf("Ответ сервера: %s\n", responseBody)

	return nil
}

func (handler *CommandHandler) Authenticate() error {
	login := handler.CLIHelper.GetLogin()
	password := handler.CLIHelper.GetPassword()

	token, responseBody, err := handler.APIClient.Authenticate(login, password)
	handler.APIClient.SetToken(token)

	if err != nil {
		return fmt.Errorf("ошибка при чтении ответа: %v", err)

	}

	fmt.Printf("Куки авторизации: %s\n", token)
	fmt.Printf("Ответ сервера: %s\n", responseBody)

	return nil
}

func (handler *CommandHandler) Register() error {
	login := handler.CLIHelper.GetLogin()
	password := handler.CLIHelper.GetPassword()

	token, responseBody, err := handler.APIClient.Registration(login, password)
	handler.APIClient.SetToken(token)

	if err != nil {
		return fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	fmt.Printf("Куки авторизации: %s\n", token)
	fmt.Printf("Ответ сервера: %s\n", responseBody)

	return nil
}

func (handler *CommandHandler) Ping() error {
	return handler.APIClient.Ping()
}
