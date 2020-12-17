package http

type ClientFactoryMethodKey struct {
	insecure bool
	redirect bool
}

var (
	TLSClientKey = ClientFactoryMethodKey{
		insecure: false,
		redirect: false,
	}
	InsecureClientKey = ClientFactoryMethodKey{
		insecure: true,
		redirect: false,
	}
	TLSRedirectClientKey = ClientFactoryMethodKey{
		insecure: false,
		redirect: true,
	}
	InsecureRedirectClientKey = ClientFactoryMethodKey{
		insecure: true,
		redirect: true,
	}
)

type ClientFactory struct {
	initializers map[ClientFactoryMethodKey]func() Client
}

func NewClientFactory() *ClientFactory {
	return &ClientFactory{
		initializers: map[ClientFactoryMethodKey]func() Client{
			TLSClientKey:              NewTLSClient,
			InsecureClientKey:         NewInsecureClient,
			TLSRedirectClientKey:      NewTLSRedirectClient,
			InsecureRedirectClientKey: NewInsecureRedirectClient,
		},
	}
}

func (f ClientFactory) Create(insecure, redirect bool) Client {
	init, ok := f.initializers[ClientFactoryMethodKey{insecure: insecure, redirect: redirect}]
	if !ok {
		return NewInsecureClient()
	}

	return init()
}
