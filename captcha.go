package captcha

import (
	"html/template"
)

type Captcha struct {
	ImageURL   template.URL
	Value      string
	Identifier string
}
