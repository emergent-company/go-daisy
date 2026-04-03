package galleryruntime

import (
	"regexp"
	"strconv"
	"strings"
)

// TokenGroup is a named group of design tokens used for rendering grouped sections
// in the token panel.
type TokenGroup struct {
	Name   string
	Tokens []DesignToken
}

// TokenGroups returns the tokens for a component organised into ordered groups,
// preserving the first-seen order of group names.
func TokenGroups(tokens []DesignToken) []TokenGroup {
	order := []string{}
	byName := map[string][]DesignToken{}
	for _, t := range tokens {
		g := t.Group
		if g == "" {
			g = "General"
		}
		if _, ok := byName[g]; !ok {
			order = append(order, g)
		}
		byName[g] = append(byName[g], t)
	}
	groups := make([]TokenGroup, 0, len(order))
	for _, name := range order {
		groups = append(groups, TokenGroup{Name: name, Tokens: byName[name]})
	}
	return groups
}

var reNonAlphaNum = regexp.MustCompile(`[^a-zA-Z0-9]+`)

// SanitizeID converts a CSS variable name like "--btn-padding-x" to a safe HTML
// id suffix like "btn-padding-x" that can be used in element ids.
func SanitizeID(cssVar string) string {
	s := reNonAlphaNum.ReplaceAllString(cssVar, "-")
	s = strings.Trim(s, "-")
	return s
}

// RangeDefaultValue extracts the numeric portion of a token's default value so
// it can be set as the initial value of an <input type="range">.
// e.g. "0.75rem" → "0.75", "16px" → "16", "50" → "50".
func RangeDefaultValue(tok DesignToken) string {
	raw := tok.Default
	// Strip unit suffix
	if tok.Unit != "" {
		raw = strings.TrimSuffix(raw, tok.Unit)
	}
	// Validate it's a number; fall back to Min if not
	if _, err := strconv.ParseFloat(strings.TrimSpace(raw), 64); err != nil {
		return strconv.FormatFloat(tok.Min, 'f', -1, 64)
	}
	return strings.TrimSpace(raw)
}

// BuildCategoryGroups organises a flat list of components into CategoryGroups
// with nested SubcategoryGroups, preserving first-seen order.
func BuildCategoryGroups(all []GalleryComponent) []CategoryGroup {
	catOrder := []Category{}
	catMap := map[Category]*CategoryGroup{}

	for _, c := range all {
		cg, exists := catMap[c.Category]
		if !exists {
			catOrder = append(catOrder, c.Category)
			catMap[c.Category] = &CategoryGroup{Name: c.Category}
			cg = catMap[c.Category]
		}

		// find or create subcategory
		var sub *SubcategoryGroup
		for i := range cg.Subcategories {
			if cg.Subcategories[i].Name == c.Subcategory {
				sub = &cg.Subcategories[i]
				break
			}
		}
		if sub == nil {
			cg.Subcategories = append(cg.Subcategories, SubcategoryGroup{Name: c.Subcategory})
			sub = &cg.Subcategories[len(cg.Subcategories)-1]
		}
		sub.Components = append(sub.Components, c)
		cg.Components = append(cg.Components, c)
	}

	result := make([]CategoryGroup, 0, len(catOrder))
	for _, cat := range catOrder {
		result = append(result, *catMap[cat])
	}
	return result
}

// ComponentBySlug looks up a component by slug from a registry list.
func ComponentBySlug(components []GalleryComponent, slug string) (GalleryComponent, bool) {
	for _, c := range components {
		if c.Slug == slug {
			return c, true
		}
	}
	return GalleryComponent{}, false
}

// StoryByName finds a named story within a component's effective variants.
// Falls back to the first story if name is empty or not found.
func StoryByName(c GalleryComponent, name string) GalleryStory {
	variants := c.EffectiveVariants()
	if name == "" {
		return variants[0]
	}
	for _, v := range variants {
		if slugify(v.Name) == name || v.Name == name {
			return v
		}
	}
	return variants[0]
}

// SlugifyStoryName converts a story name to a URL-safe slug.
// e.g. "Loading State" → "loading-state"
func SlugifyStoryName(name string) string {
	return slugify(name)
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = reNonAlphaNum.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
