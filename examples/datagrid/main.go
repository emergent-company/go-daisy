package main

import (
	"context"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/emergent-company/go-daisy/components/table"
	"github.com/emergent-company/go-daisy/components/ui"
	"github.com/emergent-company/go-daisy/render"
	"github.com/emergent-company/go-daisy/staticfs"

	"github.com/a-h/templ"
)

// User is a mock user record.
type User struct {
	ID     int
	Name   string
	Email  string
	Role   string
	Status string
}

var (
	usersMu sync.RWMutex
	users   = []User{
		{1, "Alice Johnson", "alice@example.com", "Admin", "Active"},
		{2, "Bob Smith", "bob@example.com", "Editor", "Active"},
		{3, "Charlie Brown", "charlie@example.com", "Viewer", "Inactive"},
		{4, "Diana Prince", "diana@example.com", "Admin", "Active"},
		{5, "Eve Wilson", "eve@example.com", "Editor", "Active"},
		{6, "Frank Castle", "frank@example.com", "Viewer", "Active"},
		{7, "Grace Hopper", "grace@example.com", "Admin", "Inactive"},
		{8, "Henry Ford", "henry@example.com", "Editor", "Active"},
		{9, "Iris West", "iris@example.com", "Viewer", "Active"},
		{10, "Jack Sparrow", "jack@example.com", "Editor", "Inactive"},
		{11, "Kate Bishop", "kate@example.com", "Viewer", "Active"},
		{12, "Leo Messi", "leo@example.com", "Admin", "Active"},
		{13, "Mia Wallace", "mia@example.com", "Editor", "Active"},
		{14, "Neo Anderson", "neo@example.com", "Viewer", "Inactive"},
		{15, "Olivia Benson", "olivia@example.com", "Admin", "Active"},
		{16, "Peter Parker", "peter@example.com", "Editor", "Active"},
		{17, "Quinn Fabray", "quinn@example.com", "Viewer", "Active"},
		{18, "Rachel Green", "rachel@example.com", "Editor", "Inactive"},
		{19, "Steve Rogers", "steve@example.com", "Admin", "Active"},
		{20, "Tony Stark", "tony@example.com", "Admin", "Active"},
	}
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	staticHandler := staticfs.Handler("/static")
	e.GET("/static/*", echo.WrapHandler(staticHandler))

	e.GET("/", handleIndex)
	e.GET("/datagrid", handleDataGrid)
	e.POST("/datagrid/batch-delete", handleBatchDelete)

	fmt.Println("DataGrid demo at http://localhost:8090")
	e.Logger.Fatal(e.Start(":8090"))
}

func handleIndex(c echo.Context) error {
	gridProps := buildGridProps(c)
	pageComp := datagridPage("Users", gridProps)
	render.RenderPage(c.Response().Writer, c.Request(), pageComp)
	return nil
}

func handleDataGrid(c echo.Context) error {
	gridProps := buildGridProps(c)
	render.RenderPartial(c.Response().Writer, c.Request(), table.DataGrid(gridProps))
	return nil
}

func handleBatchDelete(c echo.Context) error {
	idsStr := c.FormValue("ids")
	if idsStr == "" {
		render.RedirectAfterMutation(c.Response().Writer, c.Request(), "/")
		return nil
	}

	parts := strings.Split(idsStr, ",")
	deleteIDs := make(map[int]bool)
	for _, p := range parts {
		id, err := strconv.Atoi(strings.TrimSpace(p))
		if err == nil {
			deleteIDs[id] = true
		}
	}

	usersMu.Lock()
	filtered := make([]User, 0, len(users))
	for _, u := range users {
		if !deleteIDs[u.ID] {
			filtered = append(filtered, u)
		}
	}
	users = filtered
	usersMu.Unlock()

	w := c.Response().Writer
	r := c.Request()
	render.AppendToast(w, "success", fmt.Sprintf("Deleted %d user(s)", len(deleteIDs)))
	render.RedirectAfterMutation(w, r, "/")
	return nil
}

func buildGridProps(c echo.Context) table.DataGridProps {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	search := c.QueryParam("search")
	role := c.QueryParam("role")
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	if pageSize < 1 {
		pageSize = 5
	}

	usersMu.RLock()
	defer usersMu.RUnlock()

	// Filter
	filtered := make([]User, 0)
	for _, u := range users {
		if search != "" {
			searchLower := strings.ToLower(search)
			if !strings.Contains(strings.ToLower(u.Name), searchLower) &&
				!strings.Contains(strings.ToLower(u.Email), searchLower) {
				continue
			}
		}
		if role != "" {
			if !strings.EqualFold(u.Role, role) {
				continue
			}
		}
		filtered = append(filtered, u)
	}

	totalItems := len(filtered)
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
	if totalPages < 1 {
		totalPages = 1
	}

	// Paginate
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > totalItems {
		start = totalItems
	}
	if end > totalItems {
		end = totalItems
	}
	paginated := filtered[start:end]

	// Build rows
	rows := make([]table.DataGridRow, 0, len(paginated))
	for _, u := range paginated {
		u := u // capture
		row := table.DataGridRow{
			ID: strconv.Itoa(u.ID),
			Cells: []templ.Component{
				// Name — bold text
				templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
					_, err := io.WriteString(w, "<span class=\"font-medium\">"+u.Name+"</span>")
					return err
				}),
				// Email
				templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
					_, err := io.WriteString(w, u.Email)
					return err
				}),
				// Role — badge
				roleBadge(u.Role),
				// Status — badge
				statusBadge(u.Status),
				// Actions — action menu
				actionsMenu(u),
			},
		}
		rows = append(rows, row)
	}

	// Build base URL for pagination/search
	qp := []string{}
	if search != "" {
		qp = append(qp, "search="+search)
	}
	if role != "" {
		qp = append(qp, "role="+role)
	}
	qp = append(qp, "pageSize="+strconv.Itoa(pageSize))
	baseURL := "/datagrid?" + strings.Join(qp, "&")

	// Role filter dropdown
	roleFilter := roleFilterComponent(baseURL)

	return table.DataGridProps{
		ID:               "users",
		Columns: []table.DataGridColumn{
			{Key: "name", Label: "Name", Sortable: false},
			{Key: "email", Label: "Email", Sortable: false},
			{Key: "role", Label: "Role", Sortable: false},
			{Key: "status", Label: "Status", Sortable: false},
			{Key: "actions", Label: "", Sortable: false, Width: "w-20"},
		},
		Rows:              rows,
		CurrentPage:       page,
		TotalPages:        totalPages,
		TotalItems:        totalItems,
		PageSize:          pageSize,
		BaseURL:           baseURL,
		SearchName:        "search",
		SearchValue:       search,
		SearchPlaceholder: "Search users...",
		FilterChildren:    roleFilter,
		Selectable:        true,
		RowIDField:        "ids",
		BatchActions: []table.BatchAction{
			{Label: "Delete", Icon: "lucide:trash", Method: "post", Href: "/datagrid/batch-delete", Confirm: "Delete selected users?", Danger: true},
		},
		Striped:  true,
		Bordered: false,
		Compact:  false,
	}
}

// roleBadge returns a Badge with colour based on role.
func roleBadge(role string) templ.Component {
	intent := ui.BadgeGhost
	switch role {
	case "Admin":
		intent = ui.BadgePrimary
	case "Editor":
		intent = ui.BadgeSecondary
	case "Viewer":
		intent = ui.BadgeGhost
	}
	return ui.Badge(ui.BadgeProps{Label: role, Variant: intent, Size: ui.BadgeSizeSM})
}

// statusBadge returns a Badge with colour based on status.
func statusBadge(status string) templ.Component {
	intent := ui.StatusIntentFor(status)
	return ui.Badge(ui.BadgeProps{Label: status, Variant: intent, Size: ui.BadgeSizeSM})
}

// actionsMenu renders an action menu for a user row.
func actionsMenu(u User) templ.Component {
	return ui.ActionMenu([]ui.ActionMenuItem{
		{Label: "Edit", Href: "#", Icon: "lucide--pencil"},
		{Label: "Delete", Href: "#", Icon: "lucide--trash-2", Danger: true},
	})
}

// roleFilterComponent renders the role filter <select> for the DataGrid toolbar.
func roleFilterComponent(baseURL string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<select
			name="role"
			class="select select-sm w-36"
			hx-get="`+baseURL+`"
			hx-target="#datagrid-users-body"
			hx-trigger="change"
			hx-include="closest .card"
		>
			<option value="">All Roles</option>
			<option value="Admin">Admin</option>
			<option value="Editor">Editor</option>
			<option value="Viewer">Viewer</option>
		</select>`)
		return err
	})
}
