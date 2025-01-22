// commands/commands_test.go
package commands

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAPIClient представляет мок для APIClientInterface
type MockAPIClient struct {
	mock.Mock
}

// Authenticate мокированного метода
func (m *MockAPIClient) Authenticate(login, password string) (string, []byte, error) {
	args := m.Called(login, password)
	return args.String(0), args.Get(1).([]byte), args.Error(2)
}

// Registration мокированного метода
func (m *MockAPIClient) Registration(login, password string) (string, []byte, error) {
	args := m.Called(login, password)
	return args.String(0), args.Get(1).([]byte), args.Error(2)
}

// Get мокированного метода
func (m *MockAPIClient) Get(endpoint string, headers map[string]string) ([]byte, error) {
	args := m.Called(endpoint, headers)
	return args.Get(0).([]byte), args.Error(1)
}

// Post мокированного метода
func (m *MockAPIClient) Post(endpoint string, data interface{}, headers map[string]string) ([]byte, error) {
	args := m.Called(endpoint, data, headers)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockAPIClient) SetToken(token string) {
}

func (m *MockAPIClient) Ping() error {
	return nil
}

type MockCLIHelper struct {
	mock.Mock
}

// GetLogin мокированного метода
func (m *MockCLIHelper) GetLogin() string {
	args := m.Called()
	return args.String(0)
}

// GetPassword мокированного метода
func (m *MockCLIHelper) GetPassword() string {
	args := m.Called()
	return args.String(0)
}

// GetMetaInfo мокированного метода
func (m *MockCLIHelper) GetMetaInfo() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// EnterText мокированного метода
func (m *MockCLIHelper) EnterText(metaInfo string) (map[string]interface{}, error) {
	args := m.Called(metaInfo)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockCLIHelper) EnterDataID() (uint, error) {
	args := m.Called()
	return args.Get(0).(uint), args.Error(1)
}

// EnterBinary мокированного метода
func (m *MockCLIHelper) EnterBinary(metaInfo string) (map[string]interface{}, error) {
	args := m.Called(metaInfo)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// EnterKeyPas мокированного метода
func (m *MockCLIHelper) EnterKeyPas(metaInfo string) (map[string]interface{}, error) {
	args := m.Called(metaInfo)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// EnterKeyPas мокированного метода
func (m *MockCLIHelper) EnterInfoType() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// Тесты
// TestCommandHandler_AddData тестирует метод AddData
func TestCommandHandler_AddData(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	// Устанавливаем ожидания для вводимых данных
	mockCLIHelper.On("GetMetaInfo").Return("meta", nil)
	mockCLIHelper.On("EnterInfoType").Return("key-pas", nil)
	mockCLIHelper.On("EnterKeyPas", "meta").Return(map[string]interface{}{"login": "user", "password": "pass"}, nil)
	mockAPIClient.On("Post", "add_data", mock.Anything, mock.Anything).Return([]byte(`{"status": "success"}`), nil)

	err := handler.AddData()
	assert.NoError(t, err)

	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

func TestCommandHandler_AddData_Error_GetMetaInfo(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	// Устанавливаем ожидание для GetMetaInfo с ошибкой
	mockCLIHelper.On("GetMetaInfo").Return("", assert.AnError)

	err := handler.AddData()
	assert.Error(t, err) // Проверяем, что ошибка возникла

	// Проверяем, что ожидания выполнены
	mockCLIHelper.AssertExpectations(t)
}

func TestCommandHandler_AddData_Error_EnterInfoType(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	mockCLIHelper.On("GetMetaInfo").Return("test meta", nil)
	mockCLIHelper.On("EnterInfoType").Return("", assert.AnError) // Имитация ошибки

	err := handler.AddData()
	assert.Error(t, err) // Проверяем, что ошибка возникла

	// Проверяем, что ожидания выполнены
	mockCLIHelper.AssertExpectations(t)
}

func TestCommandHandler_AddData_Error_addKeyPas(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	mockCLIHelper.On("GetMetaInfo").Return("test meta", nil)
	mockCLIHelper.On("EnterInfoType").Return("key-pas", nil)
	mockCLIHelper.On("EnterKeyPas", "test meta").Return(map[string]interface{}{}, assert.AnError) // Имитация ошибки

	err := handler.AddData()
	assert.Error(t, err) // Проверяем, что ошибка возникла

	// Проверяем, что ожидания выполнены
	mockCLIHelper.AssertExpectations(t)
}

func TestCommandHandler_AddData_Error_addText(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	mockCLIHelper.On("GetMetaInfo").Return("test meta", nil)
	mockCLIHelper.On("EnterInfoType").Return("text", nil)
	mockCLIHelper.On("EnterText", "test meta").Return(map[string]interface{}{}, assert.AnError) // Имитация ошибки

	err := handler.AddData()
	assert.Error(t, err) // Проверяем, что ошибка возникла

	// Проверяем, что ожидания выполнены
	mockCLIHelper.AssertExpectations(t)
}

func TestCommandHandler_AddData_Error_addBinary(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	mockCLIHelper.On("GetMetaInfo").Return("test meta", nil)
	mockCLIHelper.On("EnterInfoType").Return("binary", nil)
	mockCLIHelper.On("EnterBinary", "test meta").Return(map[string]interface{}{}, assert.AnError) // Имитация ошибки

	err := handler.AddData()
	assert.Error(t, err) // Проверяем, что ошибка возникла

	// Проверяем, что ожидания выполнены
	mockCLIHelper.AssertExpectations(t)
}

// TestCommandHandler_addKeyPas проверяет метод addKeyPas
func TestCommandHandler_addKeyPas(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	metaInfo := "test meta"

	// Устанавливаем ожидания
	mockCLIHelper.On("EnterKeyPas", metaInfo).Return(map[string]interface{}{"login": "user", "password": "pass"}, nil)
	mockAPIClient.On("Post", "add_data", mock.Anything, mock.Anything).Return([]byte(`{"status": "success"}`), nil)

	// Вызываем метод addKeyPas
	err := handler.addKeyPas(mockAPIClient, metaInfo)
	assert.NoError(t, err) // Проверяем, что ошибка не возникла

	// Проверяем, что ожидания выполнены
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

// TestCommandHandler_addKeyPas_Error проверяет обработку ошибки ввода
func TestCommandHandler_addKeyPas_Error_Input(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	metaInfo := "test meta"

	// Устанавливаем ожидания
	mockCLIHelper.On("EnterKeyPas", metaInfo).Return(map[string]interface{}{"login": "user", "password": "pass"}, assert.AnError) // Имитируем ошибку ввода

	err := handler.addKeyPas(mockAPIClient, metaInfo)
	assert.Error(t, err) // Проверяем, что ошибка возникла

	// Проверяем, что ожидания выполнены
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

// TestCommandHandler_addKeyPas_Error_API проверяет обработку ошибки API
func TestCommandHandler_addKeyPas_Error_API(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	metaInfo := "test meta"

	// Устанавливаем ожидания
	mockCLIHelper.On("EnterKeyPas", metaInfo).Return(map[string]interface{}{"login": "user", "password": "pass"}, nil)
	mockAPIClient.On("Post", "add_data", mock.Anything, mock.Anything).Return([]byte{}, assert.AnError) // Имитируем ошибку API

	err := handler.addKeyPas(mockAPIClient, metaInfo)
	assert.Error(t, err) // Проверяем, что ошибка возникла

	// Проверяем, что ожидания выполнены
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

// TestCommandHandler_addText проверяет метод addText
func TestCommandHandler_addText(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	metaInfo := "test meta"

	// Устанавливаем ожидания
	mockCLIHelper.On("EnterText", metaInfo).Return(map[string]interface{}{"text": "Sample text"}, nil)
	mockAPIClient.On("Post", "add_data", mock.Anything, mock.Anything).Return([]byte(`{"status": "success"}`), nil)

	// Вызываем метод addText
	err := handler.addText(mockAPIClient, metaInfo)
	assert.NoError(t, err) // Проверяем, что нет ошибок

	// Проверяем, что ожидания выполнены
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

// TestCommandHandler_addText_InputError проверяет обработку ошибки ввода текста
func TestCommandHandler_addText_InputError(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	metaInfo := "test meta"

	// Устанавливаем ожидания для ошибки ввода
	mockCLIHelper.On("EnterText", metaInfo).Return(map[string]interface{}{}, assert.AnError)

	err := handler.addText(mockAPIClient, metaInfo)
	assert.Error(t, err) // Проверяем, что ошибка возникла

	// Проверяем, что ожидания выполнены
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

// TestCommandHandler_addText_APIError проверяет обработку ошибки API при добавлении текста
func TestCommandHandler_addText_APIError(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	metaInfo := "test meta"

	// Устанавливаем ожидания
	mockCLIHelper.On("EnterText", metaInfo).Return(map[string]interface{}{"text": "Sample text"}, nil)
	mockAPIClient.On("Post", "add_data", mock.Anything, mock.Anything).Return([]byte{}, assert.AnError)

	err := handler.addText(mockAPIClient, metaInfo)
	assert.Error(t, err) // Проверяем, что ошибка возникла

	// Проверяем, что ожидания выполнены
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

// TestCommandHandler_addBinary проверяет метод addBinary
func TestCommandHandler_addBinary(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	metaInfo := "test meta"

	// Устанавливаем ожидания
	mockCLIHelper.On("EnterBinary", metaInfo).Return(map[string]interface{}{"key": "value"}, nil)
	mockAPIClient.On("Post", "add_data", mock.Anything, mock.Anything).Return([]byte(`{"status": "success"}`), nil)

	// Вызываем метод addBinary
	err := handler.addBinary(mockAPIClient, metaInfo)
	assert.NoError(t, err) // Проверяем, что ошибка не возникла

	// Проверяем, что ожидания выполнены
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

// TestCommandHandler_addBinary_InputError проверяет обработку ошибки ввода бинарных данных
func TestCommandHandler_addBinary_InputError(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	metaInfo := "test meta"

	// Устанавливаем ожидания для ошибки ввода
	mockCLIHelper.On("EnterBinary", metaInfo).Return(map[string]interface{}{}, assert.AnError)

	err := handler.addBinary(mockAPIClient, metaInfo)
	assert.Error(t, err) // Проверяем, что ошибка возникла

	// Проверяем, что ожидания выполнены
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

// TestCommandHandler_addBinary_APIError проверяет обработку ошибки API при добавлении бинарных данных
func TestCommandHandler_addBinary_APIError(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	metaInfo := "test meta"

	// Устанавливаем ожидания
	mockCLIHelper.On("EnterBinary", metaInfo).Return(map[string]interface{}{"key": "value"}, nil)
	mockAPIClient.On("Post", "add_data", mock.Anything, mock.Anything).Return([]byte{}, assert.AnError)

	err := handler.addBinary(mockAPIClient, metaInfo)
	assert.Error(t, err) // Проверяем, что ошибка возникла

	// Проверяем, что ожидания выполнены
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

// TestCommandHandler_GetData проверяет успешное получение данных
func TestCommandHandler_GetData_Success(t *testing.T) {
	mockAPIClient := new(MockAPIClient)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
	}

	mockAPIClient.On("Get", "get_data", mock.Anything).Return([]byte(`[{"ID":1,"UserID":1,"UserData":"data1","DataType":"text"},{"ID":2,"UserID":1,"UserData":"data2","DataType":"log-pas"}]`), nil)

	data, err := handler.GetData()
	assert.NoError(t, err) // Проверка, что нет ошибки

	expectedData := []UserData{
		{ID: 1, UserID: 1, UserData: "data1", DataType: "text"},
		{ID: 2, UserID: 1, UserData: "data2", DataType: "log-pas"},
	}
	assert.Equal(t, expectedData, data) // Проверка, что полученные данные совпадают с ожидаемыми

	mockAPIClient.AssertExpectations(t)
}

// TestCommandHandler_GetData_APIError проверяет обработку ошибки при вызове API
func TestCommandHandler_GetData_APIError(t *testing.T) {
	mockAPIClient := new(MockAPIClient)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
	}

	mockAPIClient.On("Get", "get_data", mock.Anything).Return([]byte{}, assert.AnError) // Имитация ошибки API

	data, err := handler.GetData()
	assert.Error(t, err) // Проверьте, что ошибка возникла
	assert.Nil(t, data)  // Данные должны быть nil

	mockAPIClient.AssertExpectations(t)
}

// TestCommandHandler_GetData_JsonError проверяет обработку ошибки при декодировании JSON
func TestCommandHandler_GetData_JsonError(t *testing.T) {
	mockAPIClient := new(MockAPIClient)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
	}

	// Возвращаем неподходящие данные, чтобы вызвать ошибку при декодировании
	mockAPIClient.On("Get", "get_data", mock.Anything).Return([]byte(`invalid json`), nil)

	data, err := handler.GetData()
	assert.Error(t, err) // Проверьте, что ошибка возникла
	assert.Nil(t, data)  // Данные должны быть nil

	mockAPIClient.AssertExpectations(t)
}

// TestCommandHandler_DeleteData тестирует метод DeleteData
func TestCommandHandler_DeleteData(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	// Устанавливаем ожидания
	mockCLIHelper.On("EnterDataID").Return(uint(5), nil)
	mockAPIClient.On("Get", "delete_data?id=5", mock.Anything).Return([]byte(`{"status": "deleted"}`), nil)

	// Вызываем метод DeleteData
	err := handler.DeleteData()
	assert.NoError(t, err) // Указываем, что ошибка не должна возникнуть

	// Проверяем, что ожидания выполнены
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

func TestCommandHandler_DeleteData_InputError(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	// Устанавливаем ошибку для ввода ID
	mockCLIHelper.On("EnterDataID").Return(uint(0), assert.AnError)

	// Вызываем метод DeleteData
	err := handler.DeleteData()
	assert.Error(t, err) // Ошибка должна возникнуть из-за неверного ввода ID

	// Проверяем, что ожидания выполнены
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

func TestCommandHandler_DeleteData_APIError(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	// Устанавливаем ожидания
	mockCLIHelper.On("EnterDataID").Return(uint(5), nil)
	mockAPIClient.On("Get", "delete_data?id=5", mock.Anything).Return([]byte{}, assert.AnError)

	// Вызываем метод DeleteData
	err := handler.DeleteData()
	assert.Error(t, err) // Ошибка должна возникнуть из-за проблем с API

	// Проверяем, что ожидания выполнены
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

// TestCommandHandler_EditData проверяет метод EditData
func TestCommandHandler_EditData(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	// Настройка поведения для GetData
	mockAPIClient.On("Get", "get_data", mock.Anything).
		Return([]byte(`[{"ID":1,"UserID":1,"UserData":"data1","DataType":"text"}]`), nil)

	// Настройка поведения для ввода ID
	mockCLIHelper.On("EnterDataID").Return(uint(1), nil)

	// Настройка поведения для ввода метаинформации
	mockCLIHelper.On("GetMetaInfo").Return("test meta", nil)

	// Настройка поведения для ввода данных
	mockCLIHelper.On("EnterText", "test meta").Return(map[string]interface{}{"data": "new data"}, nil)

	// Настройка поведения для метода Post
	mockAPIClient.On("Post", "edit_data", mock.Anything, mock.Anything).
		Return([]byte(`{"status": "edited"}`), nil)

	// Вызываем метод EditData
	err := handler.EditData()
	assert.NoError(t, err) // Проверяем, что ошибок не возникло

	// Проверяем, что ожидания выполнены
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

// TestCommandHandler_EditData_Error_GetData проверяет обработку ошибки при получении данных
func TestCommandHandler_EditData_Error_GetData(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	// Устанавливаем ожидание для ошибки получения данных
	mockAPIClient.On("Get", "get_data", mock.Anything).Return([]byte{}, assert.AnError)

	err := handler.EditData()
	assert.Error(t, err) // Проверяем, что ошибка возникает

	// Проверяем, что ожидания выполнены
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

// TestCommandHandler_EditData_Error_EnterDataID проверяет обработку ошибки при вводе ID
func TestCommandHandler_EditData_Error_EnterDataID(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	// Настраиваем ожидаемое поведение
	mockAPIClient.On("Get", "get_data", mock.Anything).
		Return([]byte(`[{"ID":1,"UserID":1,"UserData":"data1","DataType":"text"}]`), nil)

	// Устанавливаем ожидания для ошибки ввода ID
	mockCLIHelper.On("EnterDataID").Return(uint(0), assert.AnError)

	err := handler.EditData()
	assert.Error(t, err) // Проверяем, что ошибка возникает

	// Проверяем, что ожидания выполнены
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

// TestCommandHandler_EditData_Error_EnterText проверяет обработку ошибки при вводе текста
func TestCommandHandler_EditData_Error_EnterText(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	// Настраиваем поведение получения данных
	mockAPIClient.On("Get", "get_data", mock.Anything).
		Return([]byte(`[{"ID":1,"UserID":1,"UserData":"data1","DataType":"text"}]`), nil)

	// Устанавливаем ожидание для ввода ID
	mockCLIHelper.On("EnterDataID").Return(uint(1), nil)

	// Устанавливаем ожидаемое поведение для ошибки ввода текста
	mockCLIHelper.On("GetMetaInfo").Return("test meta", nil)
	mockCLIHelper.On("EnterText", "test meta").Return(map[string]interface{}{}, assert.AnError)

	err := handler.EditData()
	assert.Error(t, err) // Проверяем, что ошибка возникает

	// Проверяем, что ожидания выполнены
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

// TestCommandHandler_EditData_Error_Post проверяет обработку ошибки при вызове Post
func TestCommandHandler_EditData_Error_Post(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	// Настраиваем поведение получения данных
	mockAPIClient.On("Get", "get_data", mock.Anything).
		Return([]byte(`[{"ID":1,"UserID":1,"UserData":"data1","DataType":"text"}]`), nil)

	// Устанавливаем ожидание для ввода ID
	mockCLIHelper.On("EnterDataID").Return(uint(1), nil)

	// Устанавливаем метаинформацию и ввод текста
	mockCLIHelper.On("GetMetaInfo").Return("test meta", nil)
	mockCLIHelper.On("EnterText", "test meta").Return(map[string]interface{}{"data": "new data"}, nil)

	// Устанавливаем поведение метода Post на возврат ошибки
	mockAPIClient.On("Post", "edit_data", mock.Anything, mock.Anything).Return([]byte{}, assert.AnError)

	err := handler.EditData()
	assert.Error(t, err) // Проверяем, что ошибка возникает

	// Проверяем, что ожидания выполнены
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

// TestCommandHandler_Authenticate тестирует метод Authenticate
func TestCommandHandler_Authenticate(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	// Устанавливаем ожидания для тестового логина и пароля
	mockCLIHelper.On("GetLogin").Return("testUser")
	mockCLIHelper.On("GetPassword").Return("testPassword")
	mockAPIClient.On("Authenticate", "testUser", "testPassword").Return("token123", []byte("Login successful"), nil)

	err := handler.Authenticate()
	assert.NoError(t, err) // Проверяем, что ошибка не возникла

	// Проверяем, что устанавливается токен
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)

	// Вы можете также добавить проверку вывода при необходимости, но обычно это делается с помощью системных тестов.
}

// TestCommandHandler_Authenticate_Error тестирует обработку ошибки при аутентификации
func TestCommandHandler_Authenticate_Error(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	// Устанавливаем ожидания для тестового логина и пароля
	mockCLIHelper.On("GetLogin").Return("testUser")
	mockCLIHelper.On("GetPassword").Return("testPassword")
	mockAPIClient.On("Authenticate", "testUser", "testPassword").Return("", []byte{}, assert.AnError)

	err := handler.Authenticate()
	assert.Error(t, err) // Проверяем, что ошибка возникла

	// Проверяем, что ожидаемое поведение было выполнено
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

// TestCommandHandler_Register тестирует метод Register
func TestCommandHandler_Register(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	// Устанавливаем ожидания для тестового логина и пароля
	mockCLIHelper.On("GetLogin").Return("newUser")
	mockCLIHelper.On("GetPassword").Return("newPassword")
	mockAPIClient.On("Registration", "newUser", "newPassword").Return("token456", []byte("Registration successful"), nil)

	err := handler.Register()
	assert.NoError(t, err) // Проверяем, что ошибка не возникла

	// Проверяем, что ожидания выполнены
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

// TestCommandHandler_Register_Error тестирует обработку ошибок при регистрации
func TestCommandHandler_Register_Error(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	// Устанавливаем ожидания для тестового логина и пароля
	mockCLIHelper.On("GetLogin").Return("newUser")
	mockCLIHelper.On("GetPassword").Return("newPassword")
	mockAPIClient.On("Registration", "newUser", "newPassword").Return("", []byte{}, assert.AnError)

	err := handler.Register()
	assert.Error(t, err) // Проверяем, что ошибка возникла

	// Проверяем, что ожидаемое поведение было выполнено
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}

// TestCommandHandler_Register_ResponseBodyError тестирует обработку ошибок с пустым ответом
func TestCommandHandler_Register_ResponseBodyError(t *testing.T) {
	mockAPIClient := new(MockAPIClient)
	mockCLIHelper := new(MockCLIHelper)

	handler := &CommandHandler{
		APIClient: mockAPIClient,
		CLIHelper: mockCLIHelper,
	}

	// Устанавливаем ожидания для тестового логина и пароля
	mockCLIHelper.On("GetLogin").Return("newUser")
	mockCLIHelper.On("GetPassword").Return("newPassword")
	mockAPIClient.On("Registration", "newUser", "newPassword").Return("token456", []byte{}, errors.New("test"))

	err := handler.Register()
	assert.Error(t, err) // Проверяем, что ошибка возникла, потому что responseBody пустой

	// Проверяем, что ожидаемое поведение было выполнено
	mockAPIClient.AssertExpectations(t)
	mockCLIHelper.AssertExpectations(t)
}
