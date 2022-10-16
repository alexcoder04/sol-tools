package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// TODO friendly use
func stringArrayContains(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

// TODO friendly use
func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// TODO friendly add
func Input(prompt string) (string, error) {
	fmt.Print(prompt)
	r := bufio.NewReader(os.Stdin)
	ans, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return ans, nil
}

func Exit(code int, message string) {
	if message != "" {
		fmt.Println(message)
	}
	fmt.Print("Program is done. Press enter to exit")
	r := bufio.NewReader(os.Stdin)
	_, _ = r.ReadString('\n')
	os.Exit(code)
}
