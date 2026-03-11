package boba

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type (
	Key      key.Binding
	keyIndex map[string]*Key
)

type KeyMap struct {
	GlobalKeys
	Navigation NavigationKeys
	Focus      FocusKeys
	Custom     map[string]key.Binding
	index      keyIndex
}

type GlobalKeys struct {
	Quit   Key
	Back   Key
	Submit Key
	Select Key
	Help   Key
}

type NavigationKeys struct {
	Up    Key
	Down  Key
	Left  Key
	Right Key
	Next  Key
	Prev  Key
}

type FocusKeys struct {
	Up     Key
	Down   Key
	Left   Key
	Right  Key
	Toggle Key
}

func (k Key) Match(msg tea.KeyPressMsg) bool {
	return key.Matches(msg, key.Binding(k))
}

func (k *KeyMap) NewBind(action string, binding key.Binding) error {
	if _, ok := k.index[action]; ok {
		return fmt.Errorf("action already exists: %s", action)
	}
	k.Custom[action] = binding
	return nil
}

func (k *KeyMap) bind(action string, keys ...string) {
	binding, ok := k.index[action]
	if !ok {
		return
	}
	b := key.Binding(*binding)
	b.SetKeys(keys...)
	*binding = Key(b)
}

func (k *KeyMap) Is(msg tea.KeyPressMsg, action string) bool {
	if binding, ok := k.index[action]; ok {
		return binding.Match(msg)
	}
	if binding, ok := k.Custom[action]; ok {
		return key.Matches(msg, binding)
	}
	return false
}

func (k *KeyMap) buildIndex() keyIndex {
	return keyIndex{
		"quit":             &k.Quit,
		"back":             &k.Back,
		"submit":           &k.Submit,
		"select":           &k.Select,
		"help":             &k.Help,
		"navigation.up":    &k.Navigation.Up,
		"navigation.down":  &k.Navigation.Down,
		"navigation.left":  &k.Navigation.Left,
		"navigation.right": &k.Navigation.Right,
		"navigation.next":  &k.Navigation.Next,
		"navigation.prev":  &k.Navigation.Prev,
		"focus.up":         &k.Focus.Up,
		"focus.down":       &k.Focus.Down,
		"focus.left":       &k.Focus.Left,
		"focus.right":      &k.Focus.Right,
		"focus.toggle":     &k.Focus.Toggle,
	}
}

func keysToConfig(k *KeyMap) KeysConfig {
	global := map[string]KeyEntry{}
	navigation := map[string]KeyEntry{}
	focus := map[string]KeyEntry{}

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
		default:
			global[action] = entry
		}
	}

	return KeysConfig{
		Global:     global,
		Navigation: navigation,
		Focus:      focus,
		Custom:     map[string]KeyEntry{},
	}
}

func DefaultKeyMap() KeyMap {
	k := KeyMap{
		GlobalKeys: GlobalKeys{
			Quit: Key(key.NewBinding(
				key.WithKeys("ctrl+c", "ctrl+q"),
				key.WithHelp("ctrl+q", "quit"),
			)),
			Back: Key(key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("esc", "back"),
			)),
			Submit: Key(key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "submit"),
			)),
			Select: Key(key.NewBinding(
				key.WithKeys(" "),
				key.WithHelp("space", "select"),
			)),
			Help: Key(key.NewBinding(
				key.WithKeys("?"),
				key.WithHelp("?", "toggle help"),
			)),
		},
		Navigation: NavigationKeys{
			Up: Key(key.NewBinding(
				key.WithKeys("k", "up"),
				key.WithHelp("↑/k", "move up"),
			)),
			Down: Key(key.NewBinding(
				key.WithKeys("j", "down"),
				key.WithHelp("↓/j", "move down"),
			)),
			Left: Key(key.NewBinding(
				key.WithKeys("h", "left"),
				key.WithHelp("←/h", "move left"),
			)),
			Right: Key(key.NewBinding(
				key.WithKeys("l", "right"),
				key.WithHelp("→/l", "move right"),
			)),
			Next: Key(key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "next"),
			)),
			Prev: Key(key.NewBinding(
				key.WithKeys("shift+tab"),
				key.WithHelp("shift+tab", "prev"),
			)),
		},
		Focus: FocusKeys{
			Up: Key(key.NewBinding(
				key.WithKeys("ctrl+k"),
				key.WithHelp("ctrl+k", "focus up"),
			)),
			Down: Key(key.NewBinding(
				key.WithKeys("ctrl+j"),
				key.WithHelp("ctrl+j", "focus down"),
			)),
			Left: Key(key.NewBinding(
				key.WithKeys("ctrl+h"),
				key.WithHelp("ctrl+h", "focus left"),
			)),
			Right: Key(key.NewBinding(
				key.WithKeys("ctrl+l"),
				key.WithHelp("ctrl+l", "focus right"),
			)),
			Toggle: Key(key.NewBinding(
				key.WithHelp("", "toggle focus mode"),
			)),
		},
		Custom: make(map[string]key.Binding),
	}
	k.index = k.buildIndex()
	return k
}
