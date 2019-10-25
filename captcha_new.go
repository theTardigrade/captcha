package captcha

import (
	"sync"
)

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
			defer waitGroup.Done()
			c.generateIdentifier()
		}()
	}

	errChan := make(chan error)

	{
		waitGroup.Add(1)

		go func(errChan chan<- error) {
			defer waitGroup.Done()
			if err := c.generateImage(opts); err != nil {
				errChan <- err
			}
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
