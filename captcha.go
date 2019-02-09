package captcha

import (
	"bytes"
	"encoding/base64"
	"html/template"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/fogleman/gg.v1"
)

type Captcha struct {
	ImageURL   template.URL
	Value      string
	Identifier string
}

var (
	letters = [...]byte{
		'B', 'C', 'D', 'F', 'G', 'H', 'J', 'K', 'L', 'M',
		'N', 'P', 'Q', 'R', 'S', 'T', 'W', 'X', 'Y', 'Z',
	}
	numbers = [...]byte{
		'1', '2', '3', '4', '5', '6', '7', '8', '9',
	}
)

func init() {
	rand.Seed(int64(time.Now().UTC().UnixNano()))
}

type newBody func(*Options) (*Captcha, error)

func New(opts Options) (*Captcha, error) {
	opts.SetDefaults()

	var body newBody

	if opts.UseConcurrency {
		body = newConcurrentBody
	} else {
		body = newSequentialBody
	}

	return body(&opts)
}

func newConcurrentBody(opts *Options) (*Captcha, error) {
	c := Captcha{}

	var waitGroup sync.WaitGroup

	if opts.UseIdentifier {
		waitGroup.Add(1)

		go func() {
			c.generateIdentifier()
			waitGroup.Done()
		}()
	}

	errChan := make(chan error)

	{
		waitGroup.Add(1)

		go func(errChan chan<- error) {
			if err := c.generateImage(opts); err != nil {
				errChan <- err
			}
			waitGroup.Done()
		}(errChan)
	}

	waitGroup.Wait()

	select {
	case err := <-errChan:
		return nil, err
	default:
		return &c, nil
	}
}

func newSequentialBody(opts *Options) (*Captcha, error) {
	c := Captcha{}

	if opts.UseIdentifier {
		c.generateIdentifier()
	}

	if err := c.generateImage(opts); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Captcha) CheckValue(value string) bool {
	return strings.ToUpper(value) == c.Value
}

func (c *Captcha) generateImage(opts *Options) error {
	width := float64(opts.Width)
	height := float64(opts.Height)
	area := width * height
	halfWidth := width / 2
	halfHeight := height / 2
	backgroundColor := opts.BackgroundColor
	textColor := opts.TextColor
	fontSize := opts.FontSize
	characterCount := opts.CharacterCount

	dc := gg.NewContext(opts.Width, opts.Height)

	r, g, b := float64(backgroundColor.R)/255, float64(backgroundColor.G)/255, float64(backgroundColor.B)/255

	switch opts.BackgroundType {
	case BackgroundCirclesType:
		{
			dc.SetRGB(1, 1, 1)
			dc.Clear()

			for x := float64(0); x < width; x += float64(rand.Intn(int(width/5)+1)) + width/80 {
				alpha := 255 - rand.Intn(129)
				dc.SetRGBA255(int(backgroundColor.R), int(backgroundColor.G), int(backgroundColor.B), alpha)
				r := float64(rand.Intn(int(area/1e3)+1)) + area/600
				y := (float64(rand.Intn(21)-10)*DefaultHeight)/height + halfHeight
				dc.DrawCircle(x, y, r)
				dc.Fill()
			}
		}
	case BackgroundFillType:
		{
			dc.SetRGB(r, g, b)
			dc.Clear()
		}
	}

	font, err := loadFont(fontSize)
	if err != nil {
		return err
	}
	dc.SetFontFace(font)

	for i := 0; i < characterCount; i++ {
		var s string

		if f := rand.Float64(); f <= opts.LetterProportion {
			i := rand.Intn(len(letters))
			s = string(letters[i])
		} else {
			i := rand.Intn(len(numbers))
			s = string(numbers[i])
		}

		c.Value += s

		alpha := 255 - rand.Intn(65)
		dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), alpha)

		w, h := dc.MeasureString(s)
		x := width/float64(characterCount)*(float64(i)+0.5) - w/2
		y := halfHeight + h/4
		r := float64(rand.Intn(65)-32) / 384

		dc.RotateAbout(r, halfWidth, halfHeight)
		dc.DrawString(s, x, y)
		dc.RotateAbout(-r, halfWidth, halfHeight)
	}

	buffer := bytes.NewBuffer(nil)
	err = dc.EncodePNG(buffer)
	if err != nil {
		return err
	}

	{
		var imageURLBuilder strings.Builder

		imageURLBuilder.WriteString("data:image/png;base64,")
		imageURLBuilder.WriteString(base64.StdEncoding.EncodeToString(buffer.Bytes()))

		c.ImageURL = template.URL(imageURLBuilder.String())
	}

	return nil
}

func (c *Captcha) generateIdentifier() {
	var builder strings.Builder

	builder.Grow(160)

	for i := 0; i < 4; i++ {
		builder.WriteString(strconv.FormatInt(rand.Int63(), 36))
		builder.WriteByte('-')
	}

	builder.WriteString(strconv.FormatInt(int64(time.Now().UTC().UnixNano()), 36))

	for i := 0; i < 5; i++ {
		builder.WriteByte('-')
		builder.WriteString(strconv.FormatInt(rand.Int63(), 36))
	}

	c.Identifier = builder.String()
}

func CheckValues(expectedValue, receivedValue string) bool {
	c := Captcha{Value: expectedValue}
	return c.CheckValue(receivedValue)
}
