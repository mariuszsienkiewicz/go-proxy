package redirect

import (
	"proxy/modules/cache/redirect"
	"proxy/modules/config"
)

var (
	ServerMap     map[string]config.Server
	redirectCache redirect.Cache
)

func init() {
	ServerMap = make(map[string]config.Server)
	redirectCache = redirect.NewInMemoryCache() // TODO move it to the configuration
}

// BuildRules builds the additional structures for the redirect rules
func BuildRules() {
	buildRedirectMap()
	BuildRegexRules()
	BuildHashRules()
}

func buildRedirectMap() {
	for _, server := range config.Config.Proxy.Servers {
		ServerMap[server.Id] = server
	}
}

// FindRedirect finds the first (hash then regex) rule that matches the query
func FindRedirect(query string, hash string) config.Server {
	// first search in cache
	cachedServer, foundInCache := redirectCache.Find(hash)
	if foundInCache {
		return cachedServer
	}

	// search in hash rules
	hashRule, hashRuleHit := FindHashRule(hash)
	if hashRuleHit {
		redirectCache.Add(hash, hashRule.Target)
		return hashRule.Target
	}

	// if none of the hash rules match, then check the regex rules
	regexRule, regexRuleHit := FindRegexRule(query)
	if regexRuleHit {
		redirectCache.Add(hash, regexRule.Target)
		return regexRule.Target
	}

	// add hash to cache
	redirectCache.Add(hash, *config.Config.Proxy.DefaultServer)

	// if none of the rules matched then return the default server
	return *config.Config.Proxy.DefaultServer
}
