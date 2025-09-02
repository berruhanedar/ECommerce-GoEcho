package pkg

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type SkipRule struct {
	PathPrefix string
	Method     string
	ExactMatch bool
}

var skipRules = []SkipRule{
	{PathPrefix: "/customer/login", Method: http.MethodPost, ExactMatch: true},
	{PathPrefix: "/customer/create", Method: http.MethodPost, ExactMatch: true},
	{PathPrefix: "/swagger", Method: "", ExactMatch: false},
}

func Skipper(c echo.Context) bool {
	path := c.Path()
	method := c.Request().Method

	for _, rule := range skipRules {
		if rule.ExactMatch {
			if path == rule.PathPrefix && (rule.Method == "" || method == rule.Method) {
				return true
			}
		} else {
			if strings.HasPrefix(path, rule.PathPrefix) && (rule.Method == "" || method == rule.Method) {
				return true
			}
		}
	}

	return false
}
