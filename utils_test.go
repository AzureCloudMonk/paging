package paging

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateLimitOffset(t *testing.T) {
	a := assert.New(t)

	a.True(ValidateLimitOffset(20, 10))
	a.True(ValidateLimitOffset(1, 1))
	a.True(ValidateLimitOffset(1, 1))
	a.True(ValidateLimitOffset(1, 0))
	a.True(ValidateLimitOffset(0, 1))
	a.False(ValidateLimitOffset(-1, 1))
	a.False(ValidateLimitOffset(1, -1))
}

func TestGetLimitFromRequest(t *testing.T) {
	a := assert.New(t)

	options := NewOptions()

	// We define a default limit...
	options.DefaultLimit = 40
	request, _ := http.NewRequest("GET", "http://example.com", nil)
	a.Equal(int64(40), GetLimitFromRequest(request, options))

	// But even we define a default limit, user must be able to override it
	// by setting it in the URL.
	options.DefaultLimit = 10
	request, _ = http.NewRequest("GET", "http://example.com/?limit=33", nil)
	a.Equal(int64(33), GetLimitFromRequest(request, options))

	// If user use a different query string key, we must set the default one.
	options.DefaultLimit = 20
	options.LimitKeyName = "perpage"
	request, _ = http.NewRequest("GET", "http://example.com/?limit=33", nil)
	a.Equal(int64(20), GetLimitFromRequest(request, options))

	// Now, let's check with a good query string key.
	options.DefaultLimit = 20
	options.LimitKeyName = "perpage"
	request, _ = http.NewRequest("GET", "http://example.com/?perpage=56", nil)
	a.Equal(int64(56), GetLimitFromRequest(request, options))
}

func TestGetOffsetFromRequest(t *testing.T) {
	a := assert.New(t)

	options := NewOptions()

	// No offset in URL, we should have 0.
	request, _ := http.NewRequest("GET", "http://example.com", nil)
	a.Equal(int64(0), GetOffsetFromRequest(request, options))

	// Offset in URL, let's get it.
	request, _ = http.NewRequest("GET", "http://example.com?offset=42", nil)
	a.Equal(int64(42), GetOffsetFromRequest(request, options))

	// If user use a different query string key, we must set the default one.
	request, _ = http.NewRequest("GET", "http://example.com?offshore=90", nil)
	a.Equal(int64(0), GetOffsetFromRequest(request, options))

	// Now, let's check with a good query string key.
	options.OffsetKeyName = "mayday"
	request, _ = http.NewRequest("GET", "http://example.com?mayday=90", nil)
	a.Equal(int64(90), GetOffsetFromRequest(request, options))
}

func TestGenerateURI(t *testing.T) {
	a := assert.New(t)

	options := NewOptions()
	a.Equal("?limit=10&offset=40", GenerateURI(int64(10), int64(40), options))

	options.LimitKeyName = "l"
	options.OffsetKeyName = "o"
	a.Equal("?l=14&o=60", GenerateURI(int64(14), int64(60), options))
}