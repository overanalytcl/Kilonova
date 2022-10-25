package web

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"

	"github.com/KiloProjects/kilonova/internal/config"
)

//go:generate /bin/sh -c "cd .. && /usr/bin/python scripts/toml_gen.py"

type Translation map[string]string
type Translations map[string]Translation

//go:embed _translations.json
var keys []byte

var translations Translations

func hasTranslationKey(line string) bool {
	_, ok := translations[line]
	return ok
}

func getText(lang, line string, args ...any) string {
	if _, ok := translations[line]; !ok {
		log.Printf("Invalid translation key %q\n", line)
		return "ERR"
	}
	if _, ok := translations[line][lang]; !ok {
		return translations[line][config.Common.DefaultLang]
	}
	return fmt.Sprintf(translations[line][lang], args...)
}

func recurse(prefix string, val map[string]any) {
	for name, val := range val {
		if str, ok := val.(string); ok {
			if _, ok = translations[prefix]; !ok {
				translations[prefix] = make(Translation)
			}
			translations[prefix][name] = str
		} else if deeper, ok := val.(map[string]any); ok {
			recurse(prefix+"."+name, deeper)
		} else {
			panic("Wtf")
		}
	}
}

func init() {
	translations = make(Translations)
	var elems = make(map[string]map[string]any)
	err := json.Unmarshal(keys, &elems)
	if err != nil {
		panic(err)
	}
	for name, children := range elems {
		recurse(name, children)
	}
}
