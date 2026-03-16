package main

import (
	"embed"
	"encoding/json"
	"strings"
)

//go:embed locales/*.json
var localeFS embed.FS

type LocaleInfo struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type I18nManager struct {
	locales map[string]map[string]interface{}
	infos   []LocaleInfo
}

func NewI18nManager() *I18nManager {
	m := &I18nManager{
		locales: make(map[string]map[string]interface{}),
	}
	m.loadAll()
	return m
}

func (m *I18nManager) loadAll() {
	entries, err := localeFS.ReadDir("locales")
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		code := strings.TrimSuffix(entry.Name(), ".json")
		data, err := localeFS.ReadFile("locales/" + entry.Name())
		if err != nil {
			continue
		}
		var locale map[string]interface{}
		if err := json.Unmarshal(data, &locale); err != nil {
			continue
		}
		m.locales[code] = locale

		name, _ := locale["_name"].(string)
		if name == "" {
			name = code
		}
		m.infos = append(m.infos, LocaleInfo{Code: code, Name: name})
	}
}

func (m *I18nManager) GetLocale(lang string) map[string]interface{} {
	if l, ok := m.locales[lang]; ok {
		return l
	}
	if l, ok := m.locales["en"]; ok {
		return l
	}
	return map[string]interface{}{}
}

func (m *I18nManager) GetAvailableLocales() []LocaleInfo {
	return m.infos
}
