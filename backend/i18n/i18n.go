package i18n

import (
	"doo-store/backend/config"
	"embed"
	"strings"

	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func GetMsgWithMap(ctx *gin.Context, key string, maps map[string]any) string {
	content := ""
	if maps == nil {
		content = ginI18n.MustGetMessage(ctx, &i18n.LocalizeConfig{
			MessageID: key,
		})
	} else {
		content = ginI18n.MustGetMessage(ctx, &i18n.LocalizeConfig{
			MessageID:    key,
			TemplateData: maps,
		})
	}
	content = strings.ReplaceAll(content, ": <no value>", "")
	if content == "" {
		return key
	} else {
		return content
	}
}

func GetErrMsg(ctx *gin.Context, key string, maps map[string]any) string {
	content := ""
	if maps == nil {
		content = GetMsgWithMap(ctx, key, maps)
	} else {
		content = GetMsgWithMap(ctx, key, maps)
	}
	return content
}

//go:embed lang/*
var fs embed.FS

func GinI18nLocalize() gin.HandlerFunc {

	acceptLangs := make([]language.Tag, len(config.Language))
	for i, lang := range config.Language {
		acceptLangs[i] = language.Make(lang)
	}

	return ginI18n.Localize(
		ginI18n.WithBundle(&ginI18n.BundleCfg{
			RootPath:         "./lang",
			AcceptLanguage:   acceptLangs,
			DefaultLanguage:  language.Chinese,
			FormatBundleFile: "yaml",
			UnmarshalFunc:    yaml.Unmarshal,
			Loader:           &ginI18n.EmbedLoader{FS: fs},
		}),
		ginI18n.WithGetLngHandle(
			func(context *gin.Context, defaultLng string) string {
				lng := context.GetHeader("language")
				if lng == "" {
					lng = context.Query("language")
				}
				if lng == "" {
					return defaultLng
				}
				return lng
			},
		))
}
