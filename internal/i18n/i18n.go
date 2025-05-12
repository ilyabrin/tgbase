package i18n

import (
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type I18n struct {
	bundle *i18n.Bundle
}

func NewI18n(localesDir string) (*I18n, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	// Загрузка файлов переводов
	files, err := filepath.Glob(filepath.Join(localesDir, "*.yaml"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		_, err := bundle.LoadMessageFile(file)
		if err != nil {
			return nil, err
		}
	}

	return &I18n{bundle: bundle}, nil
}

func (i *I18n) Localize(lang, messageID string, templateData interface{}) string {
	localizer := i18n.NewLocalizer(i.bundle, lang)
	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
	if err != nil {
		// Fallback на английский, если перевод не найден
		localizer = i18n.NewLocalizer(i.bundle, "en")
		message, _ = localizer.Localize(&i18n.LocalizeConfig{
			MessageID:    messageID,
			TemplateData: templateData,
		})
	}
	return message
}
