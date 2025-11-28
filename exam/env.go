package exam

import "os"

func ClearEnv(e E, key string) {
	e.Helper()
	if oldVal, ok := os.LookupEnv(key); ok {
		e.Cleanup(func() {
			os.Setenv(key, oldVal)
		})
	}
	os.Unsetenv(key)
}

func SetEnv(e E, key, value string) {
	e.Helper()
	if oldVal, ok := os.LookupEnv(key); ok {
		e.Cleanup(func() {
			os.Setenv(key, oldVal)
		})
	} else {
		e.Cleanup(func() {
			os.Unsetenv(key)
		})
	}
	os.Setenv(key, value)
}
