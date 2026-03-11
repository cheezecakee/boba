package boba

import (
	"os"
	"strings"

	"charm.land/bubbles/v2/key"
	"github.com/BurntSushi/toml"
)

const (
	configFile = "config.toml"
	themeFile  = "theme.toml"
	keysFile   = "keys.toml"
)

type Config struct {
	App   AppConfig
	Keys  KeysConfig
	Style StyleConfig
	Theme ThemeConfig
}

type AppConfig struct {
	AltScreen bool   `toml:"alt_screen"`
	Title     string `toml:"title"`
}

type KeyEntry struct {
	Keys []string `toml:"keys"`
	Help []string `toml:"help"`
}

type KeysConfig struct {
	Global     map[string]KeyEntry `toml:"global"`
	Navigation map[string]KeyEntry `toml:"navigation"`
	Focus      map[string]KeyEntry `toml:"focus"`
	Custom     map[string]KeyEntry `toml:"custom"`
}

type StyleConfig struct{}

type ThemeConfig struct {
	Active string               `toml:"active"`
	Themes map[string]TomlTheme `toml:"themes"`
}

func defaultConfig() Config {
	k := DefaultKeyMap()
	Keys = &k
	Keys.index = Keys.buildIndex()

	return Config{
		App: AppConfig{
			AltScreen: true,
			Title:     "",
		},
		Theme: ThemeConfig{
			Active: "dark",
		},
		Keys: keysToConfig(Keys),
	}
}

func (c *Config) load() {
	c.loadFile(configFile, &c.App)
	c.loadFile(themeFile, &c.Theme)
	c.loadFile(keysFile, &c.Keys)
}

func (c *Config) save() {
	c.saveFile(configFile, &c.App)
	c.saveFile(themeFile, &c.Theme)
	c.saveFile(keysFile, &c.Keys)
}

func (c *Config) applyKeys() {
	for action, entry := range c.Keys.Global {
		if len(entry.Keys) > 0 {
			Keys.bind(action, entry.Keys...)
		}
	}
	for action, entry := range c.Keys.Navigation {
		if len(entry.Keys) > 0 {
			Keys.bind("navigation."+action, entry.Keys...)
		}
	}
	for action, entry := range c.Keys.Focus {
		if len(entry.Keys) > 0 {
			Keys.bind("focus."+action, entry.Keys...)
		}
	}
	for action, entry := range c.Keys.Custom {
		if len(entry.Keys) > 0 {
			var help string
			if len(entry.Help) > 1 {
				help = entry.Help[1]
			}
			Keys.Custom[action] = key.NewBinding(
				key.WithKeys(entry.Keys...),
				key.WithHelp(strings.Join(entry.Keys, "/"), help),
			)
		}
	}
}

func (c *Config) apply() {
	switch c.Theme.Active {
	case "light":
		SetTheme(LightTheme())
	case "dark":
		SetTheme(DarkTheme())
	default:
		t, ok := c.Theme.Themes[c.Theme.Active]
		if !ok {
			SetTheme(DarkTheme())
			return
		}
		SetTheme(tomlToTheme(t))
	}
	c.applyKeys()
}

func (c *Config) loadFile(path string, v any) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) || info.Size() == 0 {
		c.saveFile(path, v)
		return
	}
	toml.DecodeFile(path, v)
}

func (c *Config) saveFile(path string, v any) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	toml.NewEncoder(f).Encode(v)
}
