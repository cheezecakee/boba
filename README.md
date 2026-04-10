# Boba
Boba is a mini [Bubble Tea](https://github.com/charmbracelet/bubbletea/tree/main) framework designed to facilitate screen routing and component composition.

# Features 

- Screen routing with navigation history (push / back)
- Screen factories for lazy screen creation
- Composable layouts with Compose
- Reusable UI blocks (Block)
- Directional navigation between blocks
- Automatic focus management
- Flexible selection system (single, multi, or none)
- Action-based items (Redirect, Cmd, Display)
- Model adapters for integrating Bubble Tea models
- Built-in keyboard navigation helpers
- Centralized keymap

# Tutorial 
Boba has some unique features that require minimal setup, below you can see how each of them are used.

### Create a Screen
Using boba is very simple, you start by making the usual bubble model for the screen, the only difference is Update will return (boba.Screen, tea.Cmd):

```go
package screens

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"github.com/cheezecakee/boba"
)

type Menu struct {}

func NewMenu() Screen {}

func (s *Menu) Update(msg tea.Msg) (boba.Screen, tea.Cmd) {}

func (s *Menu) View() tea.View {
	return tea.NewView("")
}
```
### Create Items 
Items represent the data displayed, but they can also have Actions that can be executed by a block.
*Note: All UI data becomes an Item inside of boba, everything is handled internally all you need to do is provide the data*

```go
items := boba.Items{
	boba.Redirect("Packs", NewPacks),
	boba.Cmd("Quit", tea.Quit),
}
```
Available item types:

- Redirect — navigate to another screen
- Cmd — run a Bubble Tea command
- Display — render a static item with no action

### Create a Block 
Blocks render items and manage selection. Think of blocks as component. Blocks can also hold pre-existing bubbles models.
```go 
block := boba.NewBlock[tea.Model](
	items,
	40, 10,
	&boba.SingleSelection{},
)
```
Blocks handle:
cursor movement
selection
item rendering

### Compose Layout 
Compose is used to arrange blocks and define navigation between them. You can vidually see how blocks aranged when composing it. Two blocks of the same type touching each other horizontally or veritcally will created a spanned blocked.

```go
	box1 := NewBlock[tea.Model]("box1", 30, 5, NoSelection)
	bax1.Items = Items{
        Display("item1"),
    }
	bar.Build("col")

	box2 := NewBlock[tea.Model]("box2", 30, 5, NoSelection)
	box2.Items = Items{
        Display("item2"),
    }
	bar.Build("col")

	bar := NewBlock[tea.Model]("bar", 60, 5, NoSelection)
	bar.Items = Items{Redirect("import", NewMenu)}
	bar.Build("bar")

	c := NewCompose()
	c.Row(bar, bar)
	c.Row(box1, box2)
```

### Execute Item 
When the user presses the submit key, execute the selected item.
```go
if item := s.compose.Focused().Selected(); item != nil {
	return boba.ExecItem(s, item)
}```

For more complete applications, see the `examples/` directory.

### Run the App

Boba includes bubbles component with minimal setup without hidenring configuration. 

# License

[MIT](https://github.com/cheezecakee/boba/blob/main/LICENSE)
