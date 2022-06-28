package testhelper

import "os"

func CreateTestEnvironment(variables map[string]string) func() {
	oldVariables := map[string]string{}

	for name, value := range variables {
		if oldValue, ok := os.LookupEnv(name); ok {
			oldVariables[name] = oldValue
		}

		os.Setenv(name, value)
	}

	return func() {
		for name := range variables {
			oldValue, ok := oldVariables[name]
			if ok {
				os.Setenv(name, oldValue)
			} else {
				os.Unsetenv(name)
			}
		}
	}
}
