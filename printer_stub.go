//go:build !windows

package main

import "fmt"

func EnumPrinters() ([]string, error) {
	return nil, fmt.Errorf("printing only supported on Windows")
}

func PrintImage(printerName, imagePath string, paperSource int) error {
	return fmt.Errorf("printing only supported on Windows")
}
