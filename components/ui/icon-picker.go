package ui

import "strings"

// iconPickerResolved is the IconPickerProps after defaults are applied and
// derived values (element IDs, active option, trigger label/glyph) are
// computed. It embeds the props so the templ body can read both.
type iconPickerResolved struct {
	IconPickerProps
	TriggerID   string
	PanelID     string
	SearchID    string
	GlyphID     string
	LabelID     string
	EmptyID     string
	IsDefault   bool
	ActiveIndex int
	Label       string
	GlyphClass  string
}

// iconPickerResolve applies zero-value defaults and derives the selection,
// trigger label, and trigger glyph from the caller's options.
func iconPickerResolve(p IconPickerProps) iconPickerResolved {
	if p.ID == "" {
		p.ID = "gd-icon-picker"
	}
	if p.DefaultLabel == "" {
		p.DefaultLabel = "Default"
	}
	if p.DefaultClass == "" {
		p.DefaultClass = "lucide--box"
	}
	if p.SearchPlaceholder == "" {
		p.SearchPlaceholder = "Search icons"
	}
	if p.EmptyLabel == "" {
		p.EmptyLabel = "No icons match"
	}
	if p.ResetLabel == "" {
		p.ResetLabel = "Reset to default"
	}

	value := strings.TrimSpace(p.Value)
	isDefault := value == strings.TrimSpace(p.DefaultValue)

	label := p.DefaultLabel
	glyph := p.DefaultClass
	active := -1
	if !isDefault {
		// Fall back to the raw value as a class for values the caller did not
		// list as options (matches the gateway's "unknown name" behaviour).
		label = value
		glyph = value
		for i, opt := range p.Options {
			if opt.Name == value || (opt.Class != "" && opt.Class == value) {
				active = i
				if opt.Label != "" {
					label = opt.Label
				}
				if opt.Class != "" {
					glyph = opt.Class
				}
				break
			}
		}
		if glyph == "" {
			glyph = p.DefaultClass
		}
	}

	return iconPickerResolved{
		IconPickerProps: p,
		TriggerID:       p.ID + "-trigger",
		PanelID:         p.ID + "-panel",
		SearchID:        p.ID + "-search",
		GlyphID:         p.ID + "-glyph",
		LabelID:         p.ID + "-label",
		EmptyID:         p.ID + "-empty",
		IsDefault:       isDefault,
		ActiveIndex:     active,
		Label:           label,
		GlyphClass:      glyph,
	}
}

// iconPickerTabindex is the roving-tabindex value for a picker tile: 0 for the
// active tile and -1 for the rest, so the grid holds a single tab stop.
func iconPickerTabindex(active bool) string {
	if active {
		return "0"
	}
	return "-1"
}

// iconPickerResetStyle hides the picker's "Reset to default" footer while the
// value is already the default. Inline style (not the hidden attribute) so it
// is never overridden by a display utility such as .btn.
func iconPickerResetStyle(isDefault bool) string {
	if isDefault {
		return "display:none"
	}
	return ""
}
