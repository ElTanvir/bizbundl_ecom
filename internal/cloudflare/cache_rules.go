package cloudflare

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/cache"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/cloudflare/cloudflare-go/v6/rulesets"
	"github.com/gofiber/fiber/v2/log"
)

type CacheRulesClient struct {
	client *cloudflare.Client
	zoneID string
}

// NewCacheRulesClientWithToken creates client with API Token (recommended)
func NewCacheRulesClientWithToken(apiToken, zoneID string) (*CacheRulesClient, error) {
	client := cloudflare.NewClient(
		option.WithAPIToken(apiToken),
	)
	return &CacheRulesClient{
		client: client,
		zoneID: zoneID,
	}, nil
}

// NewCacheRulesClient creates client with Global API Key
func NewCacheRulesClient(apiKey, email, zoneID string) (*CacheRulesClient, error) {
	client := cloudflare.NewClient(
		option.WithAPIKey(apiKey),
		option.WithAPIEmail(email),
	)
	return &CacheRulesClient{
		client: client,
		zoneID: zoneID,
	}, nil
}

// findCacheRuleset finds existing cache ruleset for the phase
func (c *CacheRulesClient) findCacheRuleset() (string, error) {
	ctx := context.TODO()

	page, err := c.client.Rulesets.List(ctx, rulesets.RulesetListParams{
		ZoneID: cloudflare.F(c.zoneID),
	})

	if err != nil {
		return "", fmt.Errorf("failed to list rulesets: %w", err)
	}

	for _, rs := range page.Result {
		if rs.Phase == rulesets.PhaseHTTPRequestCacheSettings {
			log.Infof("Found existing cache ruleset: %s", rs.ID)
			return rs.ID, nil
		}
	}

	return "", nil
}

func (c *CacheRulesClient) CreateETagAwareCacheRules(hostname string) error {
	ctx := context.TODO()

	// Define expressions - FREE PLAN: Only eq, contains, in operators
	adminExpr := `(http.request.uri.path contains "/admin" or http.request.uri.path contains "/backend" or http.request.uri.path contains "/adm" or http.request.uri.path contains "/api/admin")`
	methodExpr := `(http.request.method in {"POST" "PUT" "PATCH" "DELETE"})`
	publicExpr := fmt.Sprintf(`(http.host contains "%s" and http.request.method eq "GET")`, hostname)

	// Static assets using "contains" instead of "ends_with" or "matches"
	staticExpr := fmt.Sprintf(`(http.host contains "%s" and (http.request.uri.path contains ".css" or http.request.uri.path contains ".js" or http.request.uri.path contains ".jpg" or http.request.uri.path contains ".jpeg" or http.request.uri.path contains ".png" or http.request.uri.path contains ".webp" or http.request.uri.path contains ".svg" or http.request.uri.path contains ".ico" or http.request.uri.path contains ".woff" or http.request.uri.path contains ".woff2"))`, hostname)

	// Build rules - FREE PLAN compatible
	rules := []rulesets.RulesetUpdateParamsRuleUnion{
		// Rule 1: Bypass admin and mutating routes
		rulesets.SetCacheSettingsRuleParam{
			Action:      cloudflare.F(rulesets.SetCacheSettingsRuleActionSetCacheSettings),
			Expression:  cloudflare.F(fmt.Sprintf("(%s or %s)", adminExpr, methodExpr)),
			Description: cloudflare.F("Bypass cache for admin and mutating routes"),
			Enabled:     cloudflare.F(true),
			ActionParameters: cloudflare.F(rulesets.SetCacheSettingsRuleActionParametersParam{
				Cache: cloudflare.F(false),
				BrowserTTL: cloudflare.F(rulesets.SetCacheSettingsRuleActionParametersBrowserTTLParam{
					Mode: cloudflare.F(rulesets.SetCacheSettingsRuleActionParametersBrowserTTLModeBypassByDefault),
				}),
			}),
		},
		// Rule 2: Cache static assets (CSS, JS, images, fonts)
		rulesets.SetCacheSettingsRuleParam{
			Action:      cloudflare.F(rulesets.SetCacheSettingsRuleActionSetCacheSettings),
			Expression:  cloudflare.F(staticExpr),
			Description: cloudflare.F("Cache static assets (CSS, JS, images, fonts)"),
			Enabled:     cloudflare.F(true),
			ActionParameters: cloudflare.F(rulesets.SetCacheSettingsRuleActionParametersParam{
				Cache: cloudflare.F(true),
				EdgeTTL: cloudflare.F(rulesets.SetCacheSettingsRuleActionParametersEdgeTTLParam{
					Mode:    cloudflare.F(rulesets.SetCacheSettingsRuleActionParametersEdgeTTLModeOverrideOrigin),
					Default: cloudflare.F(int64(31536000)), // 1 year
				}),
				BrowserTTL: cloudflare.F(rulesets.SetCacheSettingsRuleActionParametersBrowserTTLParam{
					Mode:    cloudflare.F(rulesets.SetCacheSettingsRuleActionParametersBrowserTTLModeOverrideOrigin),
					Default: cloudflare.F(int64(31536000)), // 1 year
				}),
			}),
		},
		// Rule 3: Cache GET HTML pages with ETag support
		rulesets.SetCacheSettingsRuleParam{
			Action:      cloudflare.F(rulesets.SetCacheSettingsRuleActionSetCacheSettings),
			Expression:  cloudflare.F(fmt.Sprintf("%s and not (%s) and not (%s)", publicExpr, adminExpr, staticExpr)),
			Description: cloudflare.F("Cache HTML pages with ETag support"),
			Enabled:     cloudflare.F(true),
			ActionParameters: cloudflare.F(rulesets.SetCacheSettingsRuleActionParametersParam{
				Cache: cloudflare.F(true),
				EdgeTTL: cloudflare.F(rulesets.SetCacheSettingsRuleActionParametersEdgeTTLParam{
					Mode:    cloudflare.F(rulesets.SetCacheSettingsRuleActionParametersEdgeTTLModeOverrideOrigin),
					Default: cloudflare.F(int64(3600)), // 1 hour
					StatusCodeTTL: cloudflare.F([]rulesets.SetCacheSettingsRuleActionParametersEdgeTTLStatusCodeTTLParam{
						{
							StatusCodeRange: cloudflare.F(rulesets.SetCacheSettingsRuleActionParametersEdgeTTLStatusCodeTTLStatusCodeRangeParam{
								From: cloudflare.F(int64(200)),
								To:   cloudflare.F(int64(299)),
							}),
							Value: cloudflare.F(int64(3600)),
						},
						{
							StatusCode: cloudflare.F(int64(304)),
							Value:      cloudflare.F(int64(3600)),
						},
					}),
				}),
				BrowserTTL: cloudflare.F(rulesets.SetCacheSettingsRuleActionParametersBrowserTTLParam{
					Mode: cloudflare.F(rulesets.SetCacheSettingsRuleActionParametersBrowserTTLModeRespectOrigin),
				}),
				ServeStale: cloudflare.F(rulesets.SetCacheSettingsRuleActionParametersServeStaleParam{
					DisableStaleWhileUpdating: cloudflare.F(false),
				}),
			}),
		},
	}

	log.Infof("Creating/updating cache ruleset with %d rules (Free plan - basic operators)", len(rules))

	rulesetID, err := c.findCacheRuleset()
	if err != nil {
		return err
	}

	if rulesetID != "" {
		log.Infof("Updating existing ruleset: %s", rulesetID)
		result, err := c.client.Rulesets.Update(
			ctx,
			rulesetID,
			rulesets.RulesetUpdateParams{
				ZoneID: cloudflare.F(c.zoneID),
				Rules:  cloudflare.F(rules),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to update cache rules: %w", err)
		}
		fmt.Printf("✅ Cache ruleset updated: %s (version %s)\n", result.ID, result.Version)
	} else {
		log.Infof("Creating new cache ruleset")

		newRules := make([]rulesets.RulesetNewParamsRuleUnion, len(rules))
		for i, rule := range rules {
			newRules[i] = rule.(rulesets.RulesetNewParamsRuleUnion)
		}

		result, err := c.client.Rulesets.New(
			ctx,
			rulesets.RulesetNewParams{
				Name:        cloudflare.F("Cache Ruleset"),
				Kind:        cloudflare.F(rulesets.KindZone),
				Phase:       cloudflare.F(rulesets.PhaseHTTPRequestCacheSettings),
				ZoneID:      cloudflare.F(c.zoneID),
				Description: cloudflare.F("Zone-level Cache Rules (Free plan compatible)"),
				Rules:       cloudflare.F(newRules),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to create cache rules: %w", err)
		}
		fmt.Printf("✅ Cache ruleset created: %s (version %s)\n", result.ID, result.Version)
	}

	fmt.Println("✅ ETag-aware cache rules deployed (Free plan) - admin bypassed, static assets cached 1yr, HTML 1hr")
	return nil
}

// PurgeCache purges all cache
func (c *CacheRulesClient) PurgeCache() error {
	ctx := context.Background()
	params := cache.CachePurgeParams{
		ZoneID: cloudflare.F(c.zoneID),
		Body: cache.CachePurgeParamsBodyCachePurgeEverything{
			PurgeEverything: cloudflare.F(true),
		},
	}
	_, err := c.client.Cache.Purge(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to purge cache: %w", err)
	}
	fmt.Println("✅ Cache purged successfully")
	return nil
}

// PurgeCacheByURL purges specific URLs
func (c *CacheRulesClient) PurgeCacheByURL(urls []string) error {
	ctx := context.Background()
	params := cache.CachePurgeParams{
		ZoneID: cloudflare.F(c.zoneID),
		Body: cache.CachePurgeParamsBodyCachePurgeSingleFile{
			Files: cloudflare.F(urls),
		},
	}

	_, err := c.client.Cache.Purge(ctx, params)

	if err != nil {
		return fmt.Errorf("failed to purge cache by URL: %w", err)
	}

	fmt.Printf("✅ Purged %d URLs from cache\n", len(urls))
	return nil
}
