package utils

import (
	"fmt"
	"strings"
)

func ValidateValuesInList(savedList []string, value string) error {
	if len(savedList) == 0 {
		return fmt.Errorf("list is empty")
	}

	for i, _ := range savedList {
		if savedList[i] == value {
			return nil
		}
	}

	return fmt.Errorf("value should be in this [%s]", strings.Join(savedList, ", "))
}
