package config

type Provider interface {
	Load()
}

type YmlProvider struct {
	file string
}

func (ymlProvider YmlProvider) Load() {

}
