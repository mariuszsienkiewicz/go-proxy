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

// FindRedirect finds the first (regex or hash) rule that matches the query
func FindRedirect(query string, hash string) config.Server {
	// TODO check if the query hits hash rule

	// check regex rules
	for _, rule := range RegexRules {
		if rule.match(query) {
			return ServerMap[rule.Rule.Target]
		}
	}

	return *config.Config.Proxy.DefaultServer
}
