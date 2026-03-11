package boba

import "image/color"

type Theme struct {
	Primary    color.Color
	Secondary  color.Color
	Accent     color.Color
	Muted      color.Color
	Subtle     color.Color
	Danger     color.Color
	Warning    color.Color
	Success    color.Color
	Text       color.Color
	Background color.Color
}

//============ Default Themes ============//

func DarkTheme() Theme {
	return Theme{
		Primary:    NewColor("205"),
		Secondary:  NewColor("86"),
		Accent:     NewColor("205"),
		Muted:      NewColor("240"),
		Subtle:     NewColor("236"),
		Danger:     NewColor("9"),
		Warning:    NewColor("214"),
		Success:    NewColor("10"),
		Text:       NewColor("255"),
		Background: NewColor("0"),
	}
}

func LightTheme() Theme {
	return Theme{
		Primary:    NewColor("205"),
		Secondary:  NewColor("26"),
		Accent:     NewColor("205"),
		Muted:      NewColor("250"),
		Subtle:     NewColor("254"),
		Danger:     NewColor("9"),
		Warning:    NewColor("214"),
		Success:    NewColor("2"),
		Text:       NewColor("0"),
		Background: NewColor("255"),
	}
}

func DefaultTheme() Theme {
	return DarkTheme()
}

type TomlTheme struct {
	Primary    string `toml:"primary"`
	Secondary  string `toml:"secondary"`
	Accent     string `toml:"accent"`
	Muted      string `toml:"muted"`
	Subtle     string `toml:"subtle"`
	Danger     string `toml:"danger"`
	Warning    string `toml:"warning"`
	Success    string `toml:"success"`
	Text       string `toml:"text"`
	Background string `toml:"background"`
}

func tomlToTheme(t TomlTheme) Theme {
	return Theme{
		Primary:    NewColor(t.Primary),
		Secondary:  NewColor(t.Secondary),
		Accent:     NewColor(t.Accent),
		Muted:      NewColor(t.Muted),
		Subtle:     NewColor(t.Subtle),
		Danger:     NewColor(t.Danger),
		Warning:    NewColor(t.Warning),
		Success:    NewColor(t.Success),
		Text:       NewColor(t.Text),
		Background: NewColor(t.Background),
	}
}
