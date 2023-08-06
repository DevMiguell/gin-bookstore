package shared

import "fmt"

func Info(message string) {
	fmt.Printf("[INFO] %s\n", message)
}

func Error(message string) {
	fmt.Printf("[ERROR] %s\n", message)
}
