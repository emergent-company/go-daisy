package ui

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/devmode"
)

// ButtonWithBoundary wraps Button with a dev-mode component boundary annotation.
// gallery:token variant,size,typ,shape,icon,loading
// gallery:hint href:default(#)
func ButtonWithBoundary(href string, variant ButtonVariant, size ButtonSize, typ ButtonType, shape ButtonShape, icon string, loading bool) templ.Component {
	return devmode.ComponentBoundary("Button", map[string]any{
		"href":    href,
		"variant": string(variant),
		"size":    string(size),
		"type":    string(typ),
		"shape":   string(shape),
		"icon":    icon,
		"loading": loading,
	}, Button(href, variant, size, typ, shape, icon, loading))
}

// BadgeWithBoundary wraps Badge with a dev-mode component boundary annotation.
// gallery:token variant,style,size,dot,icon,label
// gallery:hint label:default(Active)
func BadgeWithBoundary(variant BadgeIntent, style BadgeStyle, size BadgeSize, dot bool, icon string, label string) templ.Component {
	props := BadgeProps{Label: label, Variant: variant, Style: style, Size: size, Dot: dot, Icon: icon}
	return devmode.ComponentBoundary("Badge", props, Badge(props))
}

// StatusBadgeWithBoundary wraps StatusBadge with a dev-mode component boundary annotation.
// gallery:token status
// gallery:hint status:default(active)
func StatusBadgeWithBoundary(status string) templ.Component {
	return devmode.ComponentBoundary("StatusBadge", map[string]any{"status": status}, StatusBadge(status))
}

// AvatarWithBoundary wraps Avatar with a dev-mode component boundary annotation.
// gallery:token name,icon,size
// gallery:hint name:default(Jane Smith)
// gallery:hint icon:default()
func AvatarWithBoundary(name string, src string, icon string, size AvatarSize) templ.Component {
	return devmode.ComponentBoundary("Avatar", map[string]any{
		"name": name,
		"src":  src,
		"icon": icon,
		"size": string(size),
	}, Avatar(name, src, icon, size))
}

// CardWithBoundary wraps Card with a dev-mode component boundary annotation.
// gallery:token title
// gallery:hint title:default(Card Title)
func CardWithBoundary(title string) templ.Component {
	return devmode.ComponentBoundary("Card", map[string]any{"title": title}, Card(title))
}

// AlertWithBoundary wraps Alert with a dev-mode component boundary annotation.
// gallery:token typ,icon,message
// gallery:hint message:default(Operation completed successfully.)
// gallery:hint icon:default(lucide--circle-check)
func AlertWithBoundary(typ AlertType, icon string, message string) templ.Component {
	return devmode.ComponentBoundary("Alert", map[string]any{
		"type":    string(typ),
		"icon":    icon,
		"message": message,
	}, Alert(typ, icon, message))
}

// AlertWithIconBoundary is a backwards-compatible alias for AlertWithBoundary.
// Deprecated: use AlertWithBoundary directly.
func AlertWithIconBoundary(typ AlertType, icon string, message string) templ.Component {
	return AlertWithBoundary(typ, icon, message)
}

// ToastWithBoundary wraps Toast with a dev-mode component boundary annotation.
// gallery:token typ,message
// gallery:hint message:default(Action completed successfully.)
func ToastWithBoundary(typ ToastType, message string) templ.Component {
	return devmode.ComponentBoundary("Toast", map[string]any{
		"type":    string(typ),
		"message": message,
	}, Toast(typ, message))
}

// PaginationWithBoundary wraps Pagination with a dev-mode component boundary annotation.
// gallery:token currentPage,totalPages
// gallery:hint currentPage:range(1,20,1)
// gallery:hint totalPages:range(1,20,1)
func PaginationWithBoundary(currentPage int, totalPages int, baseURL string, targetID string) templ.Component {
	return devmode.ComponentBoundary("Pagination", map[string]any{
		"currentPage": currentPage,
		"totalPages":  totalPages,
		"baseURL":     baseURL,
		"targetID":    targetID,
	}, Pagination(currentPage, totalPages, baseURL, targetID))
}

// StatCardWithBoundary wraps StatCard with a dev-mode component boundary annotation.
func StatCardWithBoundary(p StatCardProps) templ.Component {
	return devmode.ComponentBoundary("StatCard", p, StatCard(p))
}

// EmptyWithBoundary wraps Empty with a dev-mode component boundary annotation.
// gallery:token title,description
// gallery:hint title:default(Nothing here yet)
// gallery:hint description:default(Add some items to get started.)
func EmptyWithBoundary(icon string, title string, description string) templ.Component {
	return devmode.ComponentBoundary("Empty", map[string]any{
		"icon":        icon,
		"title":       title,
		"description": description,
	}, Empty(icon, title, description))
}

// LoaderWithBoundary wraps Loader with a dev-mode component boundary annotation.
// gallery:token variant
func LoaderWithBoundary(variant LoaderVariant) templ.Component {
	return devmode.ComponentBoundary("Loader", map[string]any{"variant": string(variant)}, Loader(variant))
}

// ActionMenuWithBoundary wraps ActionMenu with a dev-mode component boundary annotation.
// gallery:token items
// gallery:hint items:slice(3)
func ActionMenuWithBoundary(items []ActionMenuItem) templ.Component {
	return devmode.ComponentBoundary("ActionMenu", map[string]any{"itemCount": len(items)}, ActionMenu(items))
}

// FilterCardWithBoundary wraps FilterCard with a dev-mode component boundary annotation.
func FilterCardWithBoundary(props FilterCardProps) templ.Component {
	return devmode.ComponentBoundary("FilterCard", props, FilterCard(props))
}

// ProgressWithBoundary wraps Progress with a dev-mode component boundary annotation.
// gallery:token color,value,max
// gallery:hint value:range(0,100,1)
// gallery:hint value:default(70)
// gallery:hint max:range(1,200,1)
// gallery:hint max:default(100)
func ProgressWithBoundary(color ProgressColor, value int, max int) templ.Component {
	return devmode.ComponentBoundary("Progress", map[string]any{
		"color": string(color),
		"value": value,
		"max":   max,
	}, Progress(color, value, max))
}

// SkeletonWithBoundary wraps Skeleton with a dev-mode component boundary annotation.
// gallery:token classes
// gallery:hint classes:default(h-4 w-full)
func SkeletonWithBoundary(classes string) templ.Component {
	return devmode.ComponentBoundary("Skeleton", map[string]any{"classes": classes}, Skeleton(classes))
}

// SectionHeaderWithBoundary wraps SectionHeader with a dev-mode component boundary annotation.
// gallery:token title
// gallery:hint title:default(Personal Information)
func SectionHeaderWithBoundary(title string) templ.Component {
	return devmode.ComponentBoundary("SectionHeader", map[string]any{"title": title}, SectionHeader(title))
}

// NoPermissionsWithBoundary wraps NoPermissions with a dev-mode component boundary annotation.
func NoPermissionsWithBoundary() templ.Component {
	return devmode.ComponentBoundary("NoPermissions", nil, NoPermissions())
}

// StatusDotWithBoundary wraps StatusDot with a dev-mode component boundary annotation.
// gallery:token color,animate
// gallery:hint color:default(status-success)
func StatusDotWithBoundary(color StatusColor, animate bool) templ.Component {
	return devmode.ComponentBoundary("StatusDot", map[string]any{
		"color":   string(color),
		"animate": animate,
	}, StatusDot(color, animate))
}

// DividerWithBoundary wraps Divider with a dev-mode component boundary annotation.
// gallery:token color,vertical
// gallery:hint color:default()
func DividerWithBoundary(color DividerColor, vertical bool, label string) templ.Component {
	child := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, label)
		return err
	})
	inner := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return Divider(color, vertical).Render(templ.WithChildren(ctx, child), w)
	})
	return devmode.ComponentBoundary("Divider", map[string]any{
		"color":    string(color),
		"vertical": vertical,
		"label":    label,
	}, inner)
}

// KbdWithBoundary wraps Kbd with a dev-mode component boundary annotation.
// gallery:token size,key
// gallery:hint key:default(⌘K)
func KbdWithBoundary(size KbdSize, key string) templ.Component {
	child := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, key)
		return err
	})
	inner := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return Kbd(size).Render(templ.WithChildren(ctx, child), w)
	})
	return devmode.ComponentBoundary("Kbd", map[string]any{
		"size": string(size),
		"key":  key,
	}, inner)
}

// CountdownWithBoundary wraps Countdown with a dev-mode component boundary annotation.
// gallery:token days,hours,minutes,seconds
// gallery:hint days:range(0,99,1)
// gallery:hint hours:range(0,23,1)
// gallery:hint minutes:range(0,59,1)
// gallery:hint seconds:range(0,59,1)
// gallery:hint days:default(2)
// gallery:hint hours:default(10)
// gallery:hint minutes:default(24)
// gallery:hint seconds:default(45)
func CountdownWithBoundary(days, hours, minutes, seconds int) templ.Component {
	return devmode.ComponentBoundary("Countdown", map[string]any{
		"days":    days,
		"hours":   hours,
		"minutes": minutes,
		"seconds": seconds,
	}, Countdown(days, hours, minutes, seconds))
}

// TagWithBoundary wraps Tag with a dev-mode component boundary annotation.
// gallery:token label
// gallery:hint label:default(Contract Law)
func TagWithBoundary(label string, removeHref string) templ.Component {
	return devmode.ComponentBoundary("Tag", map[string]any{
		"label":      label,
		"removeHref": removeHref,
	}, Tag(label, removeHref))
}

// ChatBubbleWithBoundary wraps ChatBubble with a dev-mode component boundary annotation.
// gallery:token sent,author,timestamp,message
// gallery:hint author:default(Alice)
// gallery:hint timestamp:default(10:32 AM)
// gallery:hint message:default(Hey! How are you doing?)
func ChatBubbleWithBoundary(sent bool, author, timestamp, bubbleClass, message string) templ.Component {
	child := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, message)
		return err
	})
	inner := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return ChatBubble(sent, author, timestamp, bubbleClass).Render(templ.WithChildren(ctx, child), w)
	})
	return devmode.ComponentBoundary("ChatBubble", map[string]any{
		"sent":        sent,
		"author":      author,
		"timestamp":   timestamp,
		"bubbleClass": bubbleClass,
		"message":     message,
	}, inner)
}

// MockupBrowserWithBoundary wraps MockupBrowser with a dev-mode component boundary annotation.
// gallery:token url
// gallery:hint url:default(https://go-daisy.dev)
func MockupBrowserWithBoundary(url string) templ.Component {
	inner := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return MockupBrowser(url).Render(templ.WithChildren(ctx, MockupBrowserPlaceholder()), w)
	})
	return devmode.ComponentBoundary("MockupBrowser", map[string]any{"url": url}, inner)
}

// MockupPhoneWithBoundary wraps MockupPhone with a dev-mode component boundary annotation.
func MockupPhoneWithBoundary() templ.Component {
	inner := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return MockupPhone().Render(templ.WithChildren(ctx, MockupPhonePlaceholder()), w)
	})
	return devmode.ComponentBoundary("MockupPhone", nil, inner)
}

// MockupWindowWithBoundary wraps MockupWindow with a dev-mode component boundary annotation.
func MockupWindowWithBoundary() templ.Component {
	inner := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return MockupWindow().Render(templ.WithChildren(ctx, MockupWindowPlaceholder()), w)
	})
	return devmode.ComponentBoundary("MockupWindow", nil, inner)
}

// AccordionWithBoundary wraps Accordion + AccordionItem with a dev-mode component boundary annotation.
func AccordionWithBoundary(items []AccordionItemProps) templ.Component {
	children := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, item := range items {
			it := item
			inner := templ.ComponentFunc(func(ctx2 context.Context, w2 io.Writer) error {
				return AccordionItem(it.Title, it.Open).Render(templ.WithChildren(ctx2, it.Content), w2)
			})
			if err := inner.Render(ctx, w); err != nil {
				return err
			}
		}
		return nil
	})
	outer := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return Accordion().Render(templ.WithChildren(ctx, children), w)
	})
	return devmode.ComponentBoundary("Accordion", map[string]any{"itemCount": len(items)}, outer)
}

// AccordionItemProps holds props for a single accordion item.
type AccordionItemProps struct {
	Title   string
	Content templ.Component
	Open    bool
}

// StepsWithBoundary wraps Steps + Step with a dev-mode component boundary annotation.
// gallery:token steps
func StepsWithBoundary(steps []StepProps) templ.Component {
	children := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, s := range steps {
			if err := Step(s.Label, s.Done).Render(ctx, w); err != nil {
				return err
			}
		}
		return nil
	})
	outer := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return Steps().Render(templ.WithChildren(ctx, children), w)
	})
	return devmode.ComponentBoundary("Steps", map[string]any{"stepCount": len(steps)}, outer)
}

// StepProps holds props for a single step.
type StepProps struct {
	Label string
	Done  bool
}

// SwapWithBoundary wraps Swap with a dev-mode component boundary annotation.
// gallery:token rotate
func SwapWithBoundary(rotate bool, onContent templ.Component, offContent templ.Component) templ.Component {
	return devmode.ComponentBoundary("Swap", map[string]any{
		"rotate": rotate,
	}, Swap(rotate, onContent, offContent))
}

// HeroWithBoundary wraps Hero + HeroContent with a dev-mode component boundary annotation.
// gallery:token title,subtitle,ctaLabel,minHeight
// gallery:hint title:default(go-daisy)
// gallery:hint subtitle:default(Type-safe Templ components styled with DaisyUI for HTMX apps.)
// gallery:hint ctaLabel:default(Get Started)
// gallery:hint minHeight:default(min-h-56)
func HeroWithBoundary(minHeight string, title string, subtitle string, ctaLabel string) templ.Component {
	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return HeroBody(title, subtitle, ctaLabel).Render(ctx, w)
	})
	content := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return HeroContent(true).Render(templ.WithChildren(ctx, body), w)
	})
	outer := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return HeroSection(minHeight).Render(templ.WithChildren(ctx, content), w)
	})
	return devmode.ComponentBoundary("Hero", map[string]any{
		"title":     title,
		"subtitle":  subtitle,
		"ctaLabel":  ctaLabel,
		"minHeight": minHeight,
	}, outer)
}

// TooltipWithBoundary wraps Tooltip with a dev-mode component boundary annotation.
// gallery:token tip,position
// gallery:hint tip:default(Helpful hint)
// gallery:hint position:default()
func TooltipWithBoundary(tip string, position string, trigger templ.Component) templ.Component {
	inner := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return TooltipPositioned(tip, position).Render(templ.WithChildren(ctx, trigger), w)
	})
	return devmode.ComponentBoundary("Tooltip", map[string]any{
		"tip":      tip,
		"position": position,
	}, inner)
}

// DropdownWithBoundary wraps Dropdown with a dev-mode component boundary annotation.
// gallery:token align
// gallery:hint align:default()
func DropdownWithBoundary(align DropdownAlign, trigger templ.Component, items []DropdownItemProps) templ.Component {
	content := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := trigger.Render(ctx, w); err != nil {
			return err
		}
		menu := templ.ComponentFunc(func(ctx2 context.Context, w2 io.Writer) error {
			for _, item := range items {
				if item.Divider {
					if _, err := io.WriteString(w2, `<li class="divider my-0.5"></li>`); err != nil {
						return err
					}
					continue
				}
				it := item
				li := templ.ComponentFunc(func(_ context.Context, w3 io.Writer) error {
					_, err := io.WriteString(w3, it.Label)
					return err
				})
				if err := DropdownItem(false, it.Danger).Render(templ.WithChildren(ctx2, li), w2); err != nil {
					return err
				}
			}
			return nil
		})
		return DropdownMenu().Render(templ.WithChildren(ctx, menu), w)
	})
	outer := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return Dropdown(align).Render(templ.WithChildren(ctx, content), w)
	})
	return devmode.ComponentBoundary("Dropdown", map[string]any{
		"align": string(align),
	}, outer)
}

// DropdownItemProps holds props for a single dropdown menu item.
type DropdownItemProps struct {
	Label   string
	Divider bool
	Danger  bool
}

// JoinWithBoundary wraps Join with a dev-mode component boundary annotation.
// gallery:token vertical
func JoinWithBoundary(vertical bool, children ...templ.Component) templ.Component {
	content := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, c := range children {
			if err := c.Render(ctx, w); err != nil {
				return err
			}
		}
		return nil
	})
	outer := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return Join(vertical).Render(templ.WithChildren(ctx, content), w)
	})
	return devmode.ComponentBoundary("Join", map[string]any{"vertical": vertical}, outer)
}

// IndicatorWithBoundary wraps IndicatorWrapper with a dev-mode component boundary annotation.
func IndicatorWithBoundary(badgeClass string, badgeContent templ.Component, content templ.Component) templ.Component {
	inner := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		badge := templ.ComponentFunc(func(ctx2 context.Context, w2 io.Writer) error {
			return badgeContent.Render(ctx2, w2)
		})
		if err := IndicatorBadge("", badgeClass).Render(templ.WithChildren(ctx, badge), w); err != nil {
			return err
		}
		return content.Render(ctx, w)
	})
	outer := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return IndicatorWrapper().Render(templ.WithChildren(ctx, inner), w)
	})
	return devmode.ComponentBoundary("Indicator", map[string]any{
		"badgeClass": badgeClass,
	}, outer)
}

// StackWithBoundary wraps Stack with a dev-mode component boundary annotation.
func StackWithBoundary(children ...templ.Component) templ.Component {
	content := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, c := range children {
			if err := c.Render(ctx, w); err != nil {
				return err
			}
		}
		return nil
	})
	outer := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return Stack().Render(templ.WithChildren(ctx, content), w)
	})
	return devmode.ComponentBoundary("Stack", nil, outer)
}

// DiffWithBoundary wraps Diff with a dev-mode component boundary annotation.
// gallery:token before,after
// gallery:hint before:default(Before: Old content here)
// gallery:hint after:default(After: New content here)
func DiffWithBoundary(before string, after string) templ.Component {
	inner := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := DiffItem1().Render(templ.WithChildren(ctx, DiffItemContent(before, false)), w); err != nil {
			return err
		}
		if err := DiffItem2().Render(templ.WithChildren(ctx, DiffItemContent(after, true)), w); err != nil {
			return err
		}
		return DiffResizer().Render(ctx, w)
	})
	outer := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return DiffContainer().Render(templ.WithChildren(ctx, inner), w)
	})
	return devmode.ComponentBoundary("Diff", map[string]any{
		"before": before,
		"after":  after,
	}, outer)
}

// MaskWithBoundary wraps Mask with a dev-mode component boundary annotation.
// gallery:token shape
// gallery:hint shape:default(mask-squircle)
func MaskWithBoundary(shape MaskShape, content templ.Component) templ.Component {
	outer := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return Mask(shape).Render(templ.WithChildren(ctx, content), w)
	})
	return devmode.ComponentBoundary("Mask", map[string]any{"shape": string(shape)}, outer)
}

// CarouselWithBoundary wraps Carousel with a dev-mode component boundary annotation.
// gallery:token vertical
func CarouselWithBoundary(vertical bool, items []CarouselItemProps) templ.Component {
	children := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, item := range items {
			it := item
			if err := CarouselItem(it.ID).Render(templ.WithChildren(ctx, it.Content), w); err != nil {
				return err
			}
		}
		return nil
	})
	outer := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return Carousel(vertical).Render(templ.WithChildren(ctx, children), w)
	})
	return devmode.ComponentBoundary("Carousel", map[string]any{"vertical": vertical, "itemCount": len(items)}, outer)
}

// CarouselItemProps holds props for a single carousel slide.
type CarouselItemProps struct {
	ID      string
	Content templ.Component
}

// TimelineWithBoundary wraps Timeline with a dev-mode component boundary annotation.
func TimelineWithBoundary(items []TimelineItemProps) templ.Component {
	inner := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for i, item := range items {
			isLast := i == len(items)-1
			if err := TimelineItem(item, i == 0, isLast).Render(ctx, w); err != nil {
				return err
			}
		}
		return nil
	})
	outer := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return Timeline().Render(templ.WithChildren(ctx, inner), w)
	})
	return devmode.ComponentBoundary("Timeline", map[string]any{"itemCount": len(items)}, outer)
}

// MockupCodeWithBoundary wraps MockupCode with a dev-mode component boundary annotation.
func MockupCodeWithBoundary(lines []MockupCodeLineProps) templ.Component {
	children := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, line := range lines {
			text := templ.ComponentFunc(func(_ context.Context, w2 io.Writer) error {
				_, err := io.WriteString(w2, line.Code)
				return err
			})
			if err := MockupCodeLine(line.Prefix, line.ColorClass).Render(templ.WithChildren(ctx, text), w); err != nil {
				return err
			}
		}
		return nil
	})
	outer := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return MockupCode().Render(templ.WithChildren(ctx, children), w)
	})
	return devmode.ComponentBoundary("MockupCode", map[string]any{"lineCount": len(lines)}, outer)
}

// MockupCodeLineProps holds props for a single code mockup line.
type MockupCodeLineProps struct {
	Prefix     string
	Code       string
	ColorClass string
}

// ListWithBoundary wraps List with a dev-mode component boundary annotation.
func ListWithBoundary(items []ListRowProps) templ.Component {
	children := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, item := range items {
			if err := ListRow(item).Render(ctx, w); err != nil {
				return err
			}
		}
		return nil
	})
	outer := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return List().Render(templ.WithChildren(ctx, children), w)
	})
	return devmode.ComponentBoundary("List", map[string]any{"itemCount": len(items)}, outer)
}

// FilterTabsWithBoundary wraps FilterTabs with a dev-mode component boundary annotation.
// gallery:token selected
func FilterTabsWithBoundary(name string, selected string, tabs []string) templ.Component {
	return devmode.ComponentBoundary("FilterTabs", map[string]any{
		"name":     name,
		"selected": selected,
	}, FilterTabs(name, selected, tabs))
}

// FieldsetWithBoundary wraps Fieldset with a dev-mode component boundary annotation.
// gallery:token legend
// gallery:hint legend:default(Account Settings)
func FieldsetWithBoundary(legend string, content templ.Component) templ.Component {
	outer := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return Fieldset(legend).Render(templ.WithChildren(ctx, content), w)
	})
	return devmode.ComponentBoundary("Fieldset", map[string]any{"legend": legend}, outer)
}
