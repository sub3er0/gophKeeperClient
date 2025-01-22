package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type CLIHelper struct {
}

type CLIHelperInterface interface {
	GetLogin() string
	GetPassword() string
	GetMetaInfo() (string, error)
	EnterText(metaInfo string) (map[string]interface{}, error)
	EnterDataID() (uint, error)
	EnterBinary(metaInfo string) (map[string]interface{}, error)
	EnterKeyPas(metaInfo string) (map[string]interface{}, error)
	EnterInfoType() (string, error)
}

// GetLogin Запрашивает ввод логина у пользователя
func (cli *CLIHelper) GetLogin() string {
	var login string
	for {
		fmt.Println("Введите логин:")
		_, err := fmt.Scan(&login)

		if err == nil && login != "" {
			break
		}
		fmt.Println("Введите непустую строку!")
	}
	return login
}

// GetPassword Запрашивает ввод пароля у пользователя
func (cli *CLIHelper) GetPassword() string {
	var password string
	for {
		fmt.Println("Введите пароль:")
		_, err := fmt.Scan(&password)

		if err == nil && password != "" {
			break
		}
		fmt.Println("Введите непустую строку!")
	}
	return password
}

// GetMetaInfo Запрашивает метаинформацию у пользователя
func (cli *CLIHelper) GetMetaInfo() (string, error) {
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

func (cli *CLIHelper) EnterText(metaInfo string) (map[string]interface{}, error) {
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

func (cli *CLIHelper) EnterDataID() (uint, error) {
	var dataID uint
	for {
		_, err := fmt.Scan(&dataID)
		if err != nil {
			fmt.Printf("ошибка чтения id. Повторите попытку: %v", err)
		} else {
			break
		}
	}

	return dataID, nil
}

func (cli *CLIHelper) EnterBinary(metaInfo string) (map[string]interface{}, error) {
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

func (cli *CLIHelper) EnterKeyPas(metaInfo string) (map[string]interface{}, error) {
	var login string
	var password string
	login = cli.GetLogin()
	password = cli.GetPassword()

	data := map[string]interface{}{
		"login":     login,
		"password":  password,
		"meta-info": metaInfo,
	}

	return data, nil
}

func (cli *CLIHelper) EnterInfoType() (string, error) {
	var err error
	var infoType string
	fmt.Println("Выберите тип сохраняемой информации\n" +
		"Логин-пароль: key-pas\n" +
		"Текстовые данные: text\n" +
		"Бинарные данные: binary\n" +
		"Выход: exit")
	_, err = fmt.Scan(&infoType)

	if err != nil {
		return "", err
	}

	return infoType, nil
}
