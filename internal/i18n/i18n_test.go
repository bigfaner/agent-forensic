package i18n

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLocale(t *testing.T) {
	reset()
	assert.Equal(t, "zh", CurrentLocale())
}

func TestSetLocaleToEN(t *testing.T) {
	reset()
	err := SetLocale("en")
	assert.NoError(t, err)
	assert.Equal(t, "en", CurrentLocale())
}

func TestSetLocaleToZH(t *testing.T) {
	reset()
	// Switch to en first
	_ = SetLocale("en")
	err := SetLocale("zh")
	assert.NoError(t, err)
	assert.Equal(t, "zh", CurrentLocale())
}

func TestSetLocaleUnknown(t *testing.T) {
	reset()
	err := SetLocale("fr")
	assert.Error(t, err)
	// Should remain at default
	assert.Equal(t, "zh", CurrentLocale())
}

func TestTLookupChinese(t *testing.T) {
	reset()
	// Default locale is zh
	val := T("panel.sessions.title")
	assert.NotEqual(t, "panel.sessions.title", val) // should find a real translation
	assert.NotEmpty(t, val)
}

func TestTLookupEnglish(t *testing.T) {
	reset()
	_ = SetLocale("en")
	val := T("panel.sessions.title")
	assert.NotEqual(t, "panel.sessions.title", val)
	assert.NotEmpty(t, val)
}

func TestTFallbackToKey(t *testing.T) {
	reset()
	val := T("nonexistent.key.that.does.not.exist")
	assert.Equal(t, "nonexistent.key.that.does.not.exist", val)
}

func TestTAfterSwitch(t *testing.T) {
	reset()
	zhVal := T("status.loading")
	_ = SetLocale("en")
	enVal := T("status.loading")
	assert.NotEqual(t, zhVal, enVal)
	assert.NotEmpty(t, zhVal)
	assert.NotEmpty(t, enVal)
}

func TestLanguageSwitchInstant(t *testing.T) {
	reset()
	_ = SetLocale("en")
	assert.Equal(t, "en", CurrentLocale())
	_ = SetLocale("zh")
	assert.Equal(t, "zh", CurrentLocale())
	// No restart required - immediate effect
	val := T("status.loading")
	assert.NotEmpty(t, val)
}

func TestConcurrentAccess(t *testing.T) {
	reset()
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(n int) {
			locale := "zh"
			if n%2 == 0 {
				locale = "en"
			}
			_ = SetLocale(locale)
			_ = T("status.loading")
			done <- true
		}(i)
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestLocaleYAMLFilesExist(t *testing.T) {
	localesDir := filepath.Join("locales")
	_, errZh := os.Stat(filepath.Join(localesDir, "zh.yaml"))
	_, errEn := os.Stat(filepath.Join(localesDir, "en.yaml"))
	assert.NoError(t, errZh, "locales/zh.yaml should exist")
	assert.NoError(t, errEn, "locales/en.yaml should exist")
}

func TestAllKeysMatchBetweenLocales(t *testing.T) {
	reset()
	zhKeys := allKeys("zh")
	enKeys := allKeys("en")
	assert.Equal(t, len(zhKeys), len(enKeys), "zh and en should have the same number of keys")
	for _, k := range zhKeys {
		assert.Contains(t, enKeys, k, "key %s exists in zh but not in en", k)
	}
}

func TestKeyNamingConvention(t *testing.T) {
	reset()
	zhKeys := allKeys("zh")
	for _, k := range zhKeys {
		assert.Contains(t, k, ".", "key %q should use dot-separated paths", k)
	}
}

// reset resets the global state for test isolation
func reset() {
	mu.Lock()
	defer mu.Unlock()
	currentLocale = "zh"
	if translations == nil {
		translations = make(map[string]map[string]string)
	}
	// Force reload on next access
	loaded = false
}
