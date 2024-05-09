package redirect

import (
	"proxy/modules/cache/redirect"
	"proxy/modules/db"
	"proxy/modules/log"
)

var (
	redirectCache redirect.Cache
)

func init() {
	redirectCache = redirect.NewInMemoryCache() // TODO allow config to choose which caching needs to be used
}

// BuildRules builds the additional structures for the redirect rules
func BuildRules() {
	BuildRegexRules()
	BuildHashRules()
}

// FindRedirect finds the first (hash then regex) rule that matches the util
func FindRedirect(query string, hash string) string {
	// first search in cache
	cachedServer, foundInCache := redirectCache.Find(hash)
	if foundInCache {
		return cachedServer
	}

	// search in hash rules
	hashRule, hashRuleHit := FindHashRule(hash)
	if hashRuleHit {
		log.Logger.Tracef("Hash rule found for util: %s", query)
		redirectCache.Add(hash, hashRule.TargetGroup)
		return hashRule.TargetGroup
	}

	// if none of the hash rules match, then check the regex rules
	regexRule, regexRuleHit := FindRegexRule(query)
	if regexRuleHit {
		log.Logger.Tracef("Regex rule found for util: %s", query)
		redirectCache.Add(hash, regexRule.TargetGroup)
		return regexRule.TargetGroup
	}

	// add hash to cache
	log.Logger.Tracef("No rule found for util: %s, use default server", query)
	redirectCache.Add(hash, db.DbPool.DefaultServer.Config.ServerGroup)

	// if none of the rules matched then return the default db
	return db.DbPool.DefaultServer.Config.ServerGroup
}
