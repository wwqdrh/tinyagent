package srv

// env
// REDIS_PASSWORD=password123
// ALLOW_EMPTY_PASSWORD=yes
//
type BitnamiRedisOpt struct {
	BaseSrvOpt
}

func (o *BitnamiRedisOpt) Start() error {
	if o.Password == "" {
		o.Envs = append(o.Envs, "ALLOW_EMPTY_PASSWORD=yes")
	} else {
		o.Envs = append(o.Envs, "REDIS_PASSWORD="+o.Password)
	}
	if o.Image == "" {
		o.Image = "bitnami/redis:6.2"
	}
	return o.BaseSrvOpt.Start()
}
