package envs

import (
	"os"
	"testing"
)

// logger *log.Logger, envName, defaultValue string, mandatory bool)
type checkdata struct {
	envName                             string
	envValue                            string
	defaultValue                        string
	expectedEnvValue                    string
	expectedFatalfExecutedWhenMandatory bool
}

type loggerMock struct {
	fatalfExecuted bool
}

func (l *loggerMock) Fatalf(format string, v ...interface{}) {
	l.fatalfExecuted = true
}

func TestGetenv(t *testing.T) {
	lm := loggerMock{fatalfExecuted: false}
	checkdatas := []checkdata{{
		envName:                             "TESTING",
		envValue:                            "",
		defaultValue:                        "",
		expectedEnvValue:                    "",
		expectedFatalfExecutedWhenMandatory: true,
	}, {
		envName:                             "TESTING",
		envValue:                            "",
		defaultValue:                        "a",
		expectedEnvValue:                    "a",
		expectedFatalfExecutedWhenMandatory: false,
	}, {
		envName:                             "TESTING",
		envValue:                            "a",
		defaultValue:                        "",
		expectedEnvValue:                    "a",
		expectedFatalfExecutedWhenMandatory: false,
	}, {
		envName:                             "TESTING",
		envValue:                            "a",
		defaultValue:                        "b",
		expectedEnvValue:                    "a",
		expectedFatalfExecutedWhenMandatory: false,
	}}
	for idx, checkdata := range checkdatas {
		if checkdata.envValue != "" {
			os.Setenv(checkdata.envName, checkdata.envValue)
		} else {
			os.Unsetenv(checkdata.envName)
		}
		for _, mandatory := range []bool{false, true} {
			lm.fatalfExecuted = false
			envValue := Getenv(&lm, checkdata.envName, checkdata.defaultValue, mandatory)
			if envValue != checkdata.expectedEnvValue {
				t.Errorf("Got environment variable '%v', expected '%v'. Testdata for check %v is: mandatory = %v, %#v",
					envValue, checkdata.expectedEnvValue, idx, mandatory, checkdata)
			}
			if mandatory == false && lm.fatalfExecuted {
				t.Errorf("Fatalf executed for non-mandatory environment variable. Testdata for check %v is: mandatory = %v, %#v",
					idx, mandatory, checkdata)
			}
			if mandatory && lm.fatalfExecuted != checkdata.expectedFatalfExecutedWhenMandatory {
				t.Errorf("Fatalf executed = %v, but should have been %v. Testdata for check %v is: mandatory = %v, %#v",
					lm.fatalfExecuted, checkdata.expectedFatalfExecutedWhenMandatory, idx, mandatory, checkdata)
			}
		}
	}
}
