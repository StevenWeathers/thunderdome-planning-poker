package email

import (
	"os"
	"testing"
)

func TestGetEnvDefault(t *testing.T) {
	var TestEnv string
	TestEnv = GetEnv("TESTING_ENV", "test")

	if TestEnv != "test" {
		t.Error("Expected test, got ", TestEnv)
	}
}

func TestGetEnvNonDefault(t *testing.T) {
	os.Setenv("TESTING_ENV_NON_DEFAULT", "success")
	var TestEnv string
	TestEnv = GetEnv("TESTING_ENV_NON_DEFAULT", "test")

	if TestEnv != "success" {
		t.Error("Expected success, got ", TestEnv)
	}
}

func TestGetBoolEnvDefault(t *testing.T) {
	var TestEnv bool
	TestEnv = GetBoolEnv("TESTING_BOOL_ENV", true)

	if TestEnv != true {
		t.Error("Expected true, got ", TestEnv)
	}
}

func TestGetBoolEnvNonDefault(t *testing.T) {
	os.Setenv("TESTING_BOOL_ENV_NON_DEFAULT", "true")
	var TestEnv bool
	TestEnv = GetBoolEnv("TESTING_BOOL_ENV_NON_DEFAULT", false)

	if TestEnv != true {
		t.Error("Expected true, got ", TestEnv)
	}
}
