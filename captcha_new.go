package captcha

type newBody func(*Options) (*Captcha, error)

func New(opts Options) (captcha *Captcha, err error) {
	opts.SetDefaults()

	captcha = &Captcha{
		random: randomNew(),
	}

	if opts.UseIdentifier {
		captcha.generateIdentifier()
	}

	if err = captcha.generateImage(&opts); err != nil {
		return
	}

	return
}
