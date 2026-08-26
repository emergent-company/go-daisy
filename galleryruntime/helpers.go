package galleryruntime

import (
	"regexp"
	"sort"
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

// CategoryOrder defines the explicit display order of category groups in the
// gallery sidebar. Categories not in this slice fall to the end in first-seen order.
var CategoryOrder = []Category{
	CategoryFoundation,
	CategoryComponents,
	CategoryNavigation,
	CategoryLayout,
	CategoryPageTemplates,
	CategoryFeedback,
	CategoryForms,
	CategoryDataDisplay,
	CategoryOverlays,
}

// SubcategoryOrder defines the explicit display order of subcategories within each
// category. Subcategories not listed fall to the end in first-seen order.
var SubcategoryOrder = map[Category][]string{
	CategoryFoundation:    {"Typography", "Layout", "Effects"},
	CategoryComponents:    {"Buttons", "Badges", "Tag", "Avatars", "Display", "Cards", "Chat", "Links", "FAB"},
	CategoryNavigation:    {"Menus", "Tabs", "Headers", "Filters", "Page Title", "Pagination", "Notifications", "Search", "Profile Menus", "Misc"},
	CategoryLayout:        {"Sidebars", "Topbars", "Footers", "Page Titles", "Drawers", "Builder", "Hero", "Mockups"},
	CategoryPageTemplates: {"Pages"},
	CategoryFeedback:      {"Alerts", "States", "Loading", "Toasts", "Notifications", "Indicators", "Progress"},
	CategoryForms:         {"Inputs", "Toggles", "Filters", "Layout", "Prompt Bar", "Wizard"},
	CategoryDataDisplay:   {"Display", "Lists", "Tables"},
	CategoryOverlays:      {"Dropdowns", "Modals", "Popovers", "Panels"},
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
// with nested SubcategoryGroups, ordered by CategoryOrder (explicit rank).
// Categories not in CategoryOrder fall to the end, preserving first-seen order.
func BuildCategoryGroups(all []GalleryComponent) []CategoryGroup {
	rank := map[Category]int{}
	for i, cat := range CategoryOrder {
		rank[cat] = i
	}

	catMap := map[Category]*CategoryGroup{}

	for _, c := range all {
		if _, exists := catMap[c.Category]; !exists {
			catMap[c.Category] = &CategoryGroup{Name: c.Category}
		}
		cg := catMap[c.Category]

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

	// sort subcategories within each category group by SubcategoryOrder
	for _, cg := range catMap {
		subRank := map[string]int{}
		for i, name := range SubcategoryOrder[cg.Name] {
			subRank[name] = i
		}
		sort.Slice(cg.Subcategories, func(i, j int) bool {
			ri, rj := subRank[cg.Subcategories[i].Name], subRank[cg.Subcategories[j].Name]
			// ranked before unranked; among unranked, preserve order
			_, iRanked := subRank[cg.Subcategories[i].Name]
			_, jRanked := subRank[cg.Subcategories[j].Name]
			if iRanked != jRanked {
				return iRanked
			}
			return ri < rj
		})
	}

	// collect seen categories preserving seed.go order for unranked cats
	seen := []Category{}
	seenSet := map[Category]bool{}
	for _, c := range all {
		if seenSet[c.Category] {
			continue
		}
		seenSet[c.Category] = true
		seen = append(seen, c.Category)
	}

	// split into ranked (by CategoryOrder) and unranked (first-seen order)
	ranked := []Category{}
	unranked := []Category{}
	for _, cat := range seen {
		if _, ok := rank[cat]; ok {
			ranked = append(ranked, cat)
		} else {
			unranked = append(unranked, cat)
		}
	}
	sort.Slice(ranked, func(i, j int) bool { return rank[ranked[i]] < rank[ranked[j]] })

	catOrder := append(ranked, unranked...)
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

// Slugify converts any string to a URL-safe slug (lowercase, hyphens).
func Slugify(s string) string {
	return slugify(s)
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = reNonAlphaNum.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

// categorySlugMap maps URL-safe category slugs to Category constants.
var categorySlugMap map[string]Category

func initCategorySlugMap() {
	categorySlugMap = make(map[string]Category, len(CategoryOrder))
	for _, cat := range CategoryOrder {
		categorySlugMap[slugify(string(cat))] = cat
	}
}

// SlugifyCategory converts a category name to a URL-safe slug.
// e.g. "Data Display" → "data-display"
func SlugifyCategory(name string) string {
	return slugify(name)
}

// CategoryBySlug looks up a Category by its URL-safe slug.
// e.g. "data-display" → CategoryDataDisplay, true
func CategoryBySlug(slug string) (Category, bool) {
	if categorySlugMap == nil {
		initCategorySlugMap()
	}
	cat, ok := categorySlugMap[slug]
	return cat, ok
}

// ComponentsByCategory filters a component list to those belonging to the given category.
func ComponentsByCategory(components []GalleryComponent, category Category) []GalleryComponent {
	var out []GalleryComponent
	for _, c := range components {
		if c.Category == category {
			out = append(out, c)
		}
	}
	return out
}
