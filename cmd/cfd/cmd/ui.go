package cmd

import (
	"errors"
	"strconv"

	"github.com/manifoldco/promptui"
)

// errors
var (
	errRequired   = errors.New("required")
	errNotInteger = errors.New("not an integer")
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

// retuns a selector of `items`
// selectios are 'text to show' string array
// values are 'text to return'
//
// defaultValue is the position on the map/array that will appear as selected
func promptSelection(promptLabel string, items, values []string, defaultValue int) (result string, err error) {

	var (
		selected string
	)

	prompt := promptui.Select{
		Label:     promptLabel,
		Items:     items,
		CursorPos: defaultValue,
	}

	_, selected, err = prompt.Run()
	if err != nil {
		return
	}

	result = values[getPositionOnSlice(selected, items)]

	return

}

func getPositionOnSlice(item string, items []string) int {

	for idx, i := range items {
		if i == item {
			return idx
		}
	}

	return -1

}

// validationRequired validates the string is not empty
func validationRequired(input string) error {

	if input == "" {
		return errRequired
	}

	return nil

}

// validationInteger validates the string is not empty
func validationInteger(input string) error {

	if _, err := strconv.ParseUint(input, 10, 64); err != nil {
		return errNotInteger
	}

	return nil

}
