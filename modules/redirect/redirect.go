package redirect

import (
	"go-proxy/modules/cache"
	"go-proxy/modules/db"
	"go-proxy/modules/log"
	"go.uber.org/zap"
)

// BuildRules builds the additional structures for the redirect rules
func BuildRules() {
	BuildRegexRules()
	BuildHashRules()
}

// FindRedirect finds the first (hash then regex) rule that matches the util
func FindRedirect(query string, hash string) string {
	// first search in cache
	cachedServer, foundInCache := cache.GetCache().Get(hash)
	if foundInCache {
		return cachedServer
	}

	// search in hash rules
	hashRule, hashRuleHit := FindHashRule(hash)
	if hashRuleHit {
		log.Logger.Debug("Hash rule found", zap.String("query", query))
		cache.GetCache().Set(hash, hashRule.TargetGroup)
		return hashRule.TargetGroup
	}

	// if none of the hash rules match, then check the regex rules
	regexRule, regexRuleHit := FindRegexRule(query)
	if regexRuleHit {
		log.Logger.Debug("Regex rule found", zap.String("query", query))
		cache.GetCache().Set(hash, regexRule.TargetGroup)
		return regexRule.TargetGroup
	}

	// add hash to cache
	log.Logger.Debug("No rule found, use default server", zap.String("query", query))
	cache.GetCache().Set(hash, db.DbPool.DefaultServer.Config.ServerGroup)

	// if none of the rules matched then return the default db
	return db.DbPool.DefaultServer.Config.ServerGroup
}
