package srv

type ZipkinOpt struct {
	BaseSrvOpt
}

func (o *ZipkinOpt) Start() error {
	if o.Image == "" {
		o.Image = "openzipkin/zipkin:latest"
	}

	return o.BaseSrvOpt.Start()
}
