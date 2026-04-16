package repository

import (
	"os"
	"testing"
)

func TestLoadEnvData_DefaultHost(t *testing.T) {
	prevHost := os.Getenv("MONGO_HOST")
	prevUser := os.Getenv("MONGO_USER")
	prevPassword := os.Getenv("MONGO_PASSWORD")
	prevDB := os.Getenv("MONGO_DB")
	defer func() {
		os.Setenv("MONGO_HOST", prevHost)
		os.Setenv("MONGO_USER", prevUser)
		os.Setenv("MONGO_PASSWORD", prevPassword)
		os.Setenv("MONGO_DB", prevDB)
	}()

	os.Unsetenv("MONGO_HOST")
	os.Unsetenv("MONGO_USER")
	os.Unsetenv("MONGO_PASSWORD")
	os.Unsetenv("MONGO_DB")

	envData := loadEnvData()
	if envData.host != "localhost:27017" {
		t.Fatalf("expected default host localhost:27017, got %q", envData.host)
	}
	if envData.user != "" || envData.password != "" || envData.dbName != "" {
		t.Fatalf("expected empty credentials when env vars are unset, got %+v", envData)
	}
}

func TestConnect_MissingEnvironmentVariables(t *testing.T) {
	prevHost := os.Getenv("MONGO_HOST")
	prevUser := os.Getenv("MONGO_USER")
	prevPassword := os.Getenv("MONGO_PASSWORD")
	prevDB := os.Getenv("MONGO_DB")
	defer func() {
		os.Setenv("MONGO_HOST", prevHost)
		os.Setenv("MONGO_USER", prevUser)
		os.Setenv("MONGO_PASSWORD", prevPassword)
		os.Setenv("MONGO_DB", prevDB)
	}()

	os.Setenv("MONGO_HOST", "localhost:27017")
	os.Unsetenv("MONGO_USER")
	os.Unsetenv("MONGO_PASSWORD")
	os.Unsetenv("MONGO_DB")

	_, err := connect()
	if err == nil {
		t.Fatal("expected error when required environment variables are missing")
	}
}

func TestClose_NoClient(t *testing.T) {
	client = nil
	if err := close(); err != nil {
		t.Fatalf("expected no error when client is nil, got %v", err)
	}
}
