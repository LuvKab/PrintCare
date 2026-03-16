package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Enabled      bool      `json:"enabled"`
	PrinterName  string    `json:"printerName"`
	PaperSource  int       `json:"paperSource"`
	IntervalDays int       `json:"intervalDays"`
	ImagePath    string    `json:"imagePath"`
	Language     string    `json:"language"`
	LastPrint    time.Time `json:"lastPrintTime"`
}

func defaultConfig() *Config {
	return &Config{
		Enabled:      false,
		IntervalDays: 7,
		Language:     "en",
	}
}

func configPath() string {
	return filepath.Join(getAppDataDir(), "config.json")
}

func LoadConfig(logger *Logger) *Config {
	cfg := defaultConfig()
	data, err := os.ReadFile(configPath())
	if err != nil {
		logger.Info("No config file found, using defaults")
		return cfg
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		logger.Error("Failed to parse config: %v", err)
		return defaultConfig()
	}
	logger.Info("Config loaded from %s", configPath())
	return cfg
}

func saveConfig(cfg *Config, logger *Logger) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(configPath())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(configPath(), data, 0644); err != nil {
		logger.Error("Failed to save config: %v", err)
		return err
	}
	logger.Info("Config saved")
	return nil
}
