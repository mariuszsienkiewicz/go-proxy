package redirect

import "proxy/modules/config"

var (
	ServerMap map[string]config.Server
)

func BuildRedirectMap() {
	ServerMap = make(map[string]config.Server)
	for _, server := range config.Config.Proxy.Servers {
		ServerMap[server.Id] = server
	}
}

// FindRedirect finds the first rule that matches the query
// TODO: change return value from string to new connection with target database or at least Server
func FindRedirect(query string) config.Server {
	for _, rule := range RegexRules {
		if rule.match(query) {
			return ServerMap[rule.Rule.Target]
		}
	}

	return *config.Config.Proxy.DefaultServer
}
