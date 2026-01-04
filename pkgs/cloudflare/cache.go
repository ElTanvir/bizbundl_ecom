package cloudflare

import "github.com/rs/zerolog/log"

func Cache(){
	client,err:=NewCacheRulesClient(
		"fdf115c8851d7affac5fa2e2eb1358faeafe9",
		"aahsan.cs@gmail.com",
		"a37f4ee30e9e6660efeeb75d20b9426c",
	)
	if err!= nil{
		log.Err(err).Msg("Error Creating Client")
	}
	client.CreateETagAwareCacheRules("poshgigx.com")

}
