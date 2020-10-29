package cmd

import (
	"errors"

	"github.com/manifoldco/promptui"
)

// errors
var (
	errEmptyString = errors.New("empty string")
)

// promptText formats a prompt and returns its result
func promptText(promptLabel, promptDefault string, validateFunc promptui.ValidateFunc) (string, error) {

	var (
		prompt promptui.Prompt
	)

	prompt.Label = promptLabel
	if promptDefault != "" {
		prompt.Default = promptDefault
	}

	if validateFunc != nil {
		prompt.Validate = validateFunc
	}

	return prompt.Run()

}

// promptPassword formats a prompt for password request, returns its result and validation error if any
func promptPassword(promptLabel string, validateFunc promptui.ValidateFunc) (string, error) {

	var (
		prompt promptui.Prompt
	)

	prompt.Label = promptLabel
	prompt.Mask = '*'
	if validateFunc != nil {
		prompt.Validate = validateFunc
	}

	return prompt.Run()

}

// promptTrueFalseBool returns a prompt for true, false questions
// -  promptLabel: Text to use as label
// -   trueString: Text for true
// -  falseString: Text for false
// - defaultValue: default value
func promptTrueFalseBool(promptLabel, trueString, falseString string, defaultValue bool) (bool, error) {

	var (
		items  []string
		result string
		err    error
	)

	if defaultValue {
		items = []string{trueString, falseString}
	} else {
		items = []string{falseString, trueString}
	}

	prompt := promptui.Select{
		Label: promptLabel,
		Items: items,
	}

	_, result, err = prompt.Run()
	if err != nil {
		return false, err
	}

	if result == trueString {
		return true, err
	}
	return false, err

}

// ValidationNotEmptyString validates the string is not empty
func validationNotEmptyString(input string) error {

	if input == "" {
		return errEmptyString
	}

	return nil

}
