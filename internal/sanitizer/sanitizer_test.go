package sanitizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitize_APIKey(t *testing.T) {
	input := `api_key=sk-abc123def`
	output, masked := Sanitize(input)
	assert.True(t, masked)
	assert.Contains(t, output, "api_key=")
	assert.Contains(t, output, "***")
	assert.NotContains(t, output, "sk-abc123def")
}

func TestSanitize_APISecret(t *testing.T) {
	input := `SECRET=my_secret_value`
	output, masked := Sanitize(input)
	assert.True(t, masked)
	assert.Contains(t, output, "SECRET=")
	assert.Contains(t, output, "***")
	assert.NotContains(t, output, "my_secret_value")
}

func TestSanitize_Token(t *testing.T) {
	input := `token: ghp_abcdef123456`
	output, masked := Sanitize(input)
	assert.True(t, masked)
	assert.Contains(t, output, "token:")
	assert.Contains(t, output, "***")
	assert.NotContains(t, output, "ghp_abcdef123456")
}

func TestSanitize_Password(t *testing.T) {
	input := `password="supersecretpass"`
	output, masked := Sanitize(input)
	assert.True(t, masked)
	assert.Contains(t, output, `password="`)
	assert.Contains(t, output, "***")
	assert.NotContains(t, output, "supersecretpass")
}

func TestSanitize_CaseInsensitive(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"API_KEY uppercase", `API_KEY=sk-abc123`},
		{"Secret mixed", `Secret=value123`},
		{"TOKEN all caps", `TOKEN=abc`},
		{"Password mixed", `Password="mypass"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, masked := Sanitize(tt.input)
			assert.True(t, masked)
			assert.Contains(t, output, "***")
		})
	}
}

func TestSanitize_NoMatch(t *testing.T) {
	input := `this is a normal log line with no secrets`
	output, masked := Sanitize(input)
	assert.False(t, masked)
	assert.Equal(t, input, output)
}

func TestSanitize_EmptyString(t *testing.T) {
	output, masked := Sanitize("")
	assert.False(t, masked)
	assert.Equal(t, "", output)
}

func TestSanitize_MultipleMatches(t *testing.T) {
	input := `api_key=sk-abc token=ghp_def password="s3cret!"`
	output, masked := Sanitize(input)
	assert.True(t, masked)
	assert.NotContains(t, output, "sk-abc")
	assert.NotContains(t, output, "ghp_def")
	assert.NotContains(t, output, "s3cret!")
}

func TestSanitize_CJKContent_NoFalsePositive(t *testing.T) {
	input := `这是一个中文日志行，包含密码两个字但不应该被匹配`
	output, masked := Sanitize(input)
	assert.False(t, masked)
	assert.Equal(t, input, output)
}

func TestSanitize_CJKContent_WithSecret(t *testing.T) {
	input := `设置 api_key=sk-abc123 完成`
	output, masked := Sanitize(input)
	assert.True(t, masked)
	assert.Contains(t, output, "api_key=")
	assert.Contains(t, output, "***")
	assert.NotContains(t, output, "sk-abc123")
}

func TestSanitize_QuotedValue(t *testing.T) {
	input := `api_key="sk-proj-abc123xyz"`
	output, masked := Sanitize(input)
	assert.True(t, masked)
	assert.Contains(t, output, `api_key="`)
	assert.NotContains(t, output, "sk-proj-abc123xyz")
}

func TestSanitize_SingleQuotedValue(t *testing.T) {
	input := `secret='my_secret'`
	output, masked := Sanitize(input)
	assert.True(t, masked)
	assert.Contains(t, output, "secret='")
	assert.NotContains(t, output, "my_secret")
}

func TestSanitize_SpaceSeparator(t *testing.T) {
	input := `api_key sk-abc123`
	output, masked := Sanitize(input)
	assert.True(t, masked)
	assert.Contains(t, output, "***")
	assert.NotContains(t, output, "sk-abc123")
}

func TestSanitize_ColonSeparator(t *testing.T) {
	input := `token:ghp_abc`
	output, masked := Sanitize(input)
	assert.True(t, masked)
	assert.Contains(t, output, "***")
	assert.NotContains(t, output, "ghp_abc")
}

func TestSanitize_KeepsKeyName(t *testing.T) {
	input := `api_key=sk-abc123`
	output, _ := Sanitize(input)
	assert.Contains(t, output, "api_key=")
}

func TestSanitize_SecretInMiddle(t *testing.T) {
	input := `config: api_key=sk-abc done`
	output, masked := Sanitize(input)
	assert.True(t, masked)
	assert.Contains(t, output, "config: ")
	assert.Contains(t, output, " done")
	assert.NotContains(t, output, "sk-abc")
}

func TestSanitize_JapaneseContent_NoFalsePositive(t *testing.T) {
	input := `パスワードは秘密ですがここには何もありません`
	output, masked := Sanitize(input)
	assert.False(t, masked)
	assert.Equal(t, input, output)
}

func TestSanitize_MixedContent(t *testing.T) {
	input := `ユーザー設定: token=abc123 を保存しました`
	output, masked := Sanitize(input)
	assert.True(t, masked)
	assert.Contains(t, output, "***")
	assert.NotContains(t, output, "abc123")
}
