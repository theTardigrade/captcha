package captcha

import (
	"html/template"
	"math/rand"
)

type Captcha struct {
	ImageURL   template.URL
	Value      string
	Identifier string
	random     *rand.Rand
}
