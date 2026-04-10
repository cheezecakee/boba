package boba

import (
	"image/color"
	"os"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
	"github.com/BurntSushi/toml"
)

const (
	configFile = "config.toml"
	themeFile  = "theme.toml"
	keysFile   = "keys.toml"
	styleFile  = "style.toml"
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
	Component  map[string]KeyEntry `toml:"component"`
	Custom     map[string]KeyEntry `toml:"custom"`
}

type ThemeConfig struct {
	Active string               `toml:"active"`
	Themes map[string]TomlTheme `toml:"themes"`
}

func defaultConfig() Config {
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
	c.loadFile(styleFile, &c.Style)
}

func (c *Config) save() {
	c.saveFile(configFile, &c.App)
	c.saveFile(themeFile, &c.Theme)
	c.saveFile(keysFile, &c.Keys)
	c.saveFile(styleFile, &c.Style)
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

func (c *Config) applyStyle() {
	t := style.Theme
	if t == (Theme{}) {
		t = DefaultTheme()
	}
	style = NewStyle(style.Size.Width, style.Size.Height, t)
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
	c.applyStyle()
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

func (c *Config) Empty() bool {
	for _, path := range []string{configFile, keysFile, themeFile, styleFile} {
		info, err := os.Stat(path)
		if os.IsNotExist(err) || info.Size() == 0 {
			return true
		}
	}
	return false
}

func keysToConfig(k *KeyMap) KeysConfig {
	global := map[string]KeyEntry{}
	navigation := map[string]KeyEntry{}
	focus := map[string]KeyEntry{}
	component := map[string]KeyEntry{}

	for action, binding := range k.index {
		b := key.Binding(*binding)
		entry := KeyEntry{
			Keys: b.Keys(),
			Help: []string{b.Help().Key, b.Help().Desc},
		}
		switch {
		case strings.HasPrefix(action, "navigation."):
			navigation[strings.TrimPrefix(action, "navigation.")] = entry
		case strings.HasPrefix(action, "focus."):
			focus[strings.TrimPrefix(action, "focus.")] = entry
		case strings.HasPrefix(action, "component."):
			component[strings.TrimPrefix(action, "component.")] = entry
		default:
			global[action] = entry
		}
	}

	return KeysConfig{
		Global:     global,
		Navigation: navigation,
		Focus:      focus,
		Component:  component,
		Custom:     map[string]KeyEntry{},
	}
}

func (e StyleEntry) ToLipgloss(theme Theme) lipgloss.Style {
	s := lipgloss.NewStyle()

	switch e.Border {
	case "rounded":
		s = s.Border(lipgloss.RoundedBorder())
	case "normal":
		s = s.Border(lipgloss.NormalBorder())
	case "double":
		s = s.Border(lipgloss.DoubleBorder())
	case "hidden":
		s = s.Border(lipgloss.HiddenBorder())
	}

	if e.BorderColor != "" {
		s = s.BorderForeground(resolveColor(e.BorderColor, theme))
	}
	if e.Foreground != "" {
		s = s.Foreground(resolveColor(e.Foreground, theme))
	}
	if e.Background != "" {
		s = s.Background(resolveColor(e.Background, theme))
	}
	if e.Bold {
		s = s.Bold(true)
	}

	switch e.Align {
	case "center":
		s = s.Align(lipgloss.Center)
	case "right":
		s = s.Align(lipgloss.Right)
	case "left":
		s = s.Align(lipgloss.Left)
	}

	switch len(e.Padding) {
	case 1:
		s = s.Padding(e.Padding[0])
	case 2:
		s = s.Padding(e.Padding[0], e.Padding[1])
	case 4:
		s = s.Padding(e.Padding[0], e.Padding[1], e.Padding[2], e.Padding[3])
	}

	return s
}

func resolveColor(key string, theme Theme) color.Color {
	switch key {
	case "primary":
		return theme.Primary
	case "subtle":
		return theme.Subtle
	case "accent":
		return theme.Accent
	case "background":
		return theme.Background
	case "text":
		return theme.Text
	case "muted":
		return theme.Muted
	case "danger":
		return theme.Danger
	case "success":
		return theme.Success
	case "warning":
		return theme.Warning
	default:
		return lipgloss.Color(key) // raw hex fallback
	}
}

func (s *StyleConfig) ResolveBlock(screenName, blockName string, theme Theme) lipgloss.Style {
	base := s.Components.Container.ToLipgloss(theme)

	// global block override
	if entry, ok := s.Blocks[blockName]; ok {
		base = base.Inherit(entry.ToLipgloss(theme))
	}

	// scoped screen block override
	if screen, ok := s.Screens[screenName]; ok {
		if entry, ok := screen.Blocks[blockName]; ok {
			base = base.Inherit(entry.ToLipgloss(theme))
		}
	}

	return base
}

//======= STYLE CONFIG ========//

type StyleConfig struct {
	Global     GlobalConfig                 `toml:"global"`
	Sections   SectionsConfig               `toml:"sections"`
	Elements   ElementsConfig               `toml:"elements"`
	Components ComponentsConfig             `toml:"components"`
	Blocks     map[string]StyleEntry        `toml:"blocks"`  // global named blocks
	Screens    map[string]ScreenStyleConfig `toml:"screens"` // scoped overrides
}

type StyleEntry struct {
	Border      string `toml:"border"`       // "rounded", "normal", "double", "hidden", "none"
	BorderColor string `toml:"border_color"` // theme key: "primary", "subtle", "accent" or hex
	Foreground  string `toml:"foreground"`   // theme key or hex
	Background  string `toml:"background"`   // theme key or hex
	Padding     []int  `toml:"padding"`      // [all], [v, h], [t, r, b, l]
	Margin      []int  `toml:"margin"`
	Align       string `toml:"align"` // "left", "center", "right"
	Bold        bool   `toml:"bold"`
	Width       int    `toml:"width"`
	Height      int    `toml:"height"`
}

type GlobalConfig struct {
	CursorLeft  string `toml:"cursor_left"`
	CursorRight string `toml:"cursor_right"`
}

type SectionsConfig struct {
	Header StyleEntry `toml:"header"`
	Main   StyleEntry `toml:"main"`
	Footer StyleEntry `toml:"footer"`
}

type ElementsConfig struct {
	Title   StyleEntry `toml:"title"`
	Text    StyleEntry `toml:"text"`
	Label   StyleEntry `toml:"label"`
	Badge   StyleEntry `toml:"badge"`
	Divider StyleEntry `toml:"divider"`
}

type ComponentsConfig struct {
	Container        StyleEntry `toml:"container"`
	ContainerFocused StyleEntry `toml:"container_focused"`
	Item             StyleEntry `toml:"item"`
	ItemSelected     StyleEntry `toml:"item_selected"`
	Blank            StyleEntry `toml:"blank"`
	Content          StyleEntry `toml:"content"`
}

type ScreenStyleConfig struct {
	Composite StyleEntry            `toml:"composite"`
	Blocks    map[string]StyleEntry `toml:"blocks"`
}
