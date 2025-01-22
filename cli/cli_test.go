// cli/cli_test.go
package cli

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCLIHelper_GetLogin(t *testing.T) {
	// Создаем временный файл для имитации стандартного ввода
	tempFile, err := ioutil.TempFile("", "input.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name()) // Удаляем файл после завершения теста

	// Записываем тестовые данные в временный файл
	if _, err := tempFile.WriteString("testUser\n"); err != nil {
		t.Fatal(err)
	}

	// Перенаправляем стандартный ввод на временный файл
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin, err = os.Open(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	cli := &CLIHelper{}
	login := cli.GetLogin()
	assert.Equal(t, "testUser", login)
}

func TestCLIHelper_GetPassword(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "input.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString("testPassword\n"); err != nil {
		t.Fatal(err)
	}

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin, err = os.Open(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	cli := &CLIHelper{}
	password := cli.GetPassword()
	assert.Equal(t, "testPassword", password)
}

func TestCLIHelper_GetMetaInfo_Yes(t *testing.T) {
	// Creación de un archivo temporal para simular la entrada estándar
	tempFile, err := ioutil.TempFile("", "input.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("yes\nmetaInfoValue\n")
	assert.NoError(t, err)

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin, err = os.Open(tempFile.Name())
	assert.NoError(t, err)

	cli := CLIHelper{}
	metaInfo, err := cli.GetMetaInfo()
	assert.NoError(t, err)
	assert.Equal(t, "metaInfoValue", metaInfo) // Обратите внимание на '\n' в конце
}

func TestCLIHelper_GetMetaInfo_No(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "input.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString("no\n"); err != nil {
		t.Fatal(err)
	}

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin, err = os.Open(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	cli := &CLIHelper{}
	metaInfo, err := cli.GetMetaInfo()
	assert.NoError(t, err)
	assert.Empty(t, metaInfo) // При выборе "no" метаинформация должна быть пустой
}

func TestCLIHelper_EnterText(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "input.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("some text\n")
	assert.NoError(t, err)

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin, err = os.Open(tempFile.Name())
	assert.NoError(t, err)

	cli := &CLIHelper{}
	data, err := cli.EnterText("meta")
	assert.NoError(t, err)
	assert.Equal(t, "some text", data["text"])
	assert.Equal(t, "meta", data["meta-info"])
}

func TestCLIHelper_EnterDataID(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "input.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("42\n")
	assert.NoError(t, err)

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin, err = os.Open(tempFile.Name())
	assert.NoError(t, err)

	cli := &CLIHelper{}
	dataID, err := cli.EnterDataID()
	assert.NoError(t, err)
	assert.Equal(t, uint(42), dataID)
}

func TestCLIHelper_EnterBinary(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "input.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte("binary data\n"))
	assert.NoError(t, err)

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin, err = os.Open(tempFile.Name())
	assert.NoError(t, err)

	cli := &CLIHelper{}
	dataInfo, err := cli.EnterBinary("meta")
	assert.NoError(t, err)
	assert.Equal(t, "meta", dataInfo["meta-info"])
	assert.Equal(t, []byte("binary data\n"), dataInfo["binary-data"])
}

func TestCLIHelper_EnterKeyPas(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "input.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("testUser\ntestPassword\n")
	assert.NoError(t, err)

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin, err = os.Open(tempFile.Name())
	assert.NoError(t, err)

	cli := &CLIHelper{}
	data, err := cli.EnterKeyPas("meta")
	assert.NoError(t, err)
	assert.Equal(t, "testUser", data["login"])
	assert.Equal(t, "testPassword", data["password"])
	assert.Equal(t, "meta", data["meta-info"])
}

func TestCLIHelper_EnterInfoType(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "input.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("key-pas\n")
	assert.NoError(t, err)

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin, err = os.Open(tempFile.Name())
	assert.NoError(t, err)

	cli := &CLIHelper{}
	infoType, err := cli.EnterInfoType()
	assert.NoError(t, err)
	assert.Equal(t, "key-pas", infoType)
}
