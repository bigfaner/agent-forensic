package i18n

import (
	"embed"
	"fmt"
	"sync"

	"gopkg.in/yaml.v3"
)

//go:embed locales/*.yaml
var localeFS embed.FS

var (
	mu            sync.RWMutex
	currentLocale = "zh"
	translations  map[string]map[string]string // locale → key → value
	loaded        bool
)

// loadAll loads all locale YAML files from the embedded filesystem.
func loadAll() error {
	mu.Lock()
	defer mu.Unlock()

	if loaded {
		return nil
	}

	translations = make(map[string]map[string]string)

	files := []struct {
		locale string
		path   string
	}{
		{"zh", "locales/zh.yaml"},
		{"en", "locales/en.yaml"},
	}

	for _, f := range files {
		data, err := localeFS.ReadFile(f.path)
		if err != nil {
			return fmt.Errorf("failed to read locale file %s: %w", f.path, err)
		}

		m := make(map[string]string)
		if err := yaml.Unmarshal(data, &m); err != nil {
			return fmt.Errorf("failed to parse locale file %s: %w", f.path, err)
		}
		translations[f.locale] = m
	}

	loaded = true
	return nil
}

// ensureLoaded lazily loads locale data on first access.
// Must NOT be called while mu is held — loadAll acquires mu.Lock().
func ensureLoaded() {
	mu.RLock()
	already := loaded
	mu.RUnlock()
	if !already {
		_ = loadAll()
	}
}

// T looks up a translation key in the current locale.
// Falls back to key itself if not found.
func T(key string) string {
	ensureLoaded()

	mu.RLock()
	defer mu.RUnlock()

	localeMap, ok := translations[currentLocale]
	if !ok {
		return key
	}
	if val, found := localeMap[key]; found {
		return val
	}
	return key
}

// SetLocale switches the active locale (zh or en).
// Returns error for unknown locale codes.
func SetLocale(code string) error {
	ensureLoaded()

	mu.Lock()
	defer mu.Unlock()

	if _, ok := translations[code]; !ok {
		return fmt.Errorf("unknown locale: %s", code)
	}
	currentLocale = code
	return nil
}

// CurrentLocale returns the active locale code.
func CurrentLocale() string {
	mu.RLock()
	defer mu.RUnlock()
	return currentLocale
}

// allKeys returns all keys for a given locale (for testing).
func allKeys(locale string) []string {
	ensureLoaded()

	mu.RLock()
	defer mu.RUnlock()

	localeMap, ok := translations[locale]
	if !ok {
		return nil
	}
	keys := make([]string, 0, len(localeMap))
	for k := range localeMap {
		keys = append(keys, k)
	}
	return keys
}
