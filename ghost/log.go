// Copyright 2017 Jonathan Pincas

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ghost

import (
	"fmt"
	"log"

	"testing"

	"github.com/spf13/viper"
	"github.com/wsxiaoys/terminal/color"
)

//Log outputs a custom ghost log entry
func Log(module string, isOk bool, message string, err error) {

	var entryType string
	if isOk {
		entryType = color.Sprint("@gOK")
	} else {
		entryType = color.Sprint("@rERROR")
	}

	errorText := ""
	if err != nil {
		errorText = color.Sprint(" | @rError: " + err.Error())
	}

	logText := fmt.Sprintf("%s | %s | %s %s", module, entryType, message, errorText)
	log.Println(logText)
}

//LogFatal outputs a custom ghost log entry and exits
func LogFatal(module string, isOk bool, message string, err error) {

	errorText := ""
	if err != nil {
		errorText = color.Sprint(" | @rError: " + err.Error())
	}

	entryType := color.Sprint("@{!r}FATAL")
	logText := fmt.Sprintf("%s | %s | %s %s", module, entryType, message, errorText)
	log.Fatalln(logText)
}

//LogDebug outputs a custom ghost log entry if debugmode is on
func LogDebug(module string, isOk bool, message string, err error) {

	errorText := ""
	if err != nil {
		errorText = color.Sprint(" | @rError: " + err.Error())
	}

	if viper.GetBool("debug") {

		entryType := color.Sprint("@yDEBUG")
		logText := fmt.Sprintf("%s | %s | %s %s", module, entryType, message, errorText)
		log.Println(logText)

	}

}

func TestErrorFatal(t *testing.T, description string, got string, expected string) {

	logText := fmt.Sprintf(`
	Test: %s 
	Got : %s 
	Exp : %s`,
		description, got, expected)
	t.Errorf(logText)

}
