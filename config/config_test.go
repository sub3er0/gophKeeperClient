package config

import (
	"os"
	"testing"
)

var cfg ConfigData

func TestInitConfig_FileNotFound(t *testing.T) {
	os.Setenv("CONFIG", "nonexistent.json")
	defer os.Unsetenv("CONFIG")

	config := &Configuration{}
	_, err := config.InitConfig(&cfg)
	if err != nil {
		t.Log("Expected warning for nonexistent file, proceeding.")
	}

	expectedAddress := "http://localhost:8080"
	if cfg.ServerAddress != expectedAddress {
		t.Errorf("Expected server address %s, got %s", expectedAddress, cfg.ServerAddress)
	}
}

func TestInitConfig_FallbackToDefault(t *testing.T) {
	os.Unsetenv("CONFIG")
	os.Unsetenv("SERVER_ADDRESS")

	config := &Configuration{}
	_, err := config.InitConfig(&cfg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedAddress := "http://localhost:8080"
	if cfg.ServerAddress != expectedAddress {
		t.Errorf("Expected server address %s, got %s", expectedAddress, cfg.ServerAddress)
	}
}

func TestInitConfig_FromFile(t *testing.T) {
	// Создаем временный конфигурационный файл
	configContent := `{"server_address": "http://localhost:9090"}`
	tmpFile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()
	os.Setenv("CONFIG", tmpFile.Name())
	defer os.Unsetenv("CONFIG")

	config := &Configuration{}
	_, err = config.InitConfig(&cfg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedAddress := "http://localhost:9090"
	if cfg.ServerAddress != expectedAddress {
		t.Errorf("Expected server address %s, got %s", expectedAddress, cfg.ServerAddress)
	}
}

func TestInitConfig_FromEnv(t *testing.T) {
	os.Setenv("SERVER_ADDRESS", "http://localhost:8081")
	defer os.Unsetenv("SERVER_ADDRESS")

	config := &Configuration{}
	_, err := config.InitConfig(&cfg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedAddress := "http://localhost:8081"
	if cfg.ServerAddress != expectedAddress {
		t.Errorf("Expected server address %s, got %s", expectedAddress, cfg.ServerAddress)
	}
}

func TestInitConfig_InvalidJson(t *testing.T) {
	configContent := `{"server_address": "invalid"`
	tmpFile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()
	os.Setenv("CONFIG", tmpFile.Name())
	defer os.Unsetenv("CONFIG")

	config := &Configuration{}
	_, err = config.InitConfig(&cfg)
	if err == nil {
		t.Error("Expected error due to invalid JSON, got none")
	}
}
