package tui

import (
	"context"
	"fmt"
	"strings"

	"outline-cli/internal/api"
	"outline-cli/internal/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewState int

const (
	viewCollections viewState = iota
	viewDocuments
	viewDocumentDetail
	viewSearch
)

type Model struct {
	state      viewState
	client     *api.Client
	width      int
	height     int
	err        error
	loading    bool
	breadcrumb []string

	// Collections view
	collections []models.Collection
	collCursor  int

	// Documents view
	currentCollection *models.Collection
	docNodes          []models.NavigationNode
	flatDocs          []flatDoc
	docCursor         int

	// Document detail view
	currentDoc    *models.Document
	docScroll     int
	docLines      []string

	// Search view
	searchQuery   string
	searchResults []models.SearchResult
	searchCursor  int
	searchActive  bool
}

type flatDoc struct {
	Node  models.NavigationNode
	Depth int
}

type collectionsLoadedMsg struct {
	collections []models.Collection
}

type documentsLoadedMsg struct {
	nodes []models.NavigationNode
}

type documentLoadedMsg struct {
	doc *models.Document
}

type searchResultsMsg struct {
	results []models.SearchResult
}

type errMsg struct {
	err error
}

func NewModel(client *api.Client) Model {
	return Model{
		client:     client,
		state:      viewCollections,
		loading:    true,
		breadcrumb: []string{"Collections"},
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadCollections()
}

func (m Model) loadCollections() tea.Cmd {
	return func() tea.Msg {
		collections, _, err := m.client.Collections.List(context.Background(), models.CollectionListParams{
			PaginationParams: models.PaginationParams{Limit: 100},
		})
		if err != nil {
			return errMsg{err}
		}
		return collectionsLoadedMsg{collections}
	}
}

func (m Model) loadDocuments(collectionID string) tea.Cmd {
	return func() tea.Msg {
		nodes, err := m.client.Collections.Documents(context.Background(), collectionID)
		if err != nil {
			return errMsg{err}
		}
		return documentsLoadedMsg{nodes}
	}
}

func (m Model) loadDocument(id string) tea.Cmd {
	return func() tea.Msg {
		doc, err := m.client.Documents.Info(context.Background(), id)
		if err != nil {
			return errMsg{err}
		}
		return documentLoadedMsg{doc}
	}
}

func (m Model) searchDocuments(query string) tea.Cmd {
	return func() tea.Msg {
		results, _, err := m.client.Documents.Search(context.Background(), models.SearchParams{
			Query: query,
			PaginationParams: models.PaginationParams{Limit: 25},
		})
		if err != nil {
			return errMsg{err}
		}
		return searchResultsMsg{results}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.searchActive {
			return m.updateSearch(msg)
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			return m, nil
		case "/":
			m.searchActive = true
			m.searchQuery = ""
			m.state = viewSearch
			m.breadcrumb = []string{"Search"}
			return m, nil
		}

		switch m.state {
		case viewCollections:
			return m.updateCollections(msg)
		case viewDocuments:
			return m.updateDocuments(msg)
		case viewDocumentDetail:
			return m.updateDocumentDetail(msg)
		}

	case collectionsLoadedMsg:
		m.collections = msg.collections
		m.loading = false
		return m, nil

	case documentsLoadedMsg:
		m.docNodes = msg.nodes
		m.flatDocs = flattenNodes(msg.nodes, 0)
		m.docCursor = 0
		m.loading = false
		return m, nil

	case documentLoadedMsg:
		m.currentDoc = msg.doc
		m.docLines = strings.Split(msg.doc.Text, "\n")
		m.docScroll = 0
		m.loading = false
		return m, nil

	case searchResultsMsg:
		m.searchResults = msg.results
		m.searchCursor = 0
		m.loading = false
		return m, nil

	case errMsg:
		m.err = msg.err
		m.loading = false
		return m, nil
	}

	return m, nil
}

func (m Model) updateCollections(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.collCursor > 0 {
			m.collCursor--
		}
	case "down", "j":
		if m.collCursor < len(m.collections)-1 {
			m.collCursor++
		}
	case "enter":
		if len(m.collections) > 0 {
			coll := m.collections[m.collCursor]
			m.currentCollection = &coll
			m.state = viewDocuments
			m.loading = true
			m.breadcrumb = []string{"Collections", coll.Name}
			return m, m.loadDocuments(coll.ID)
		}
	}
	return m, nil
}

func (m Model) updateDocuments(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.docCursor > 0 {
			m.docCursor--
		}
	case "down", "j":
		if m.docCursor < len(m.flatDocs)-1 {
			m.docCursor++
		}
	case "enter":
		if len(m.flatDocs) > 0 {
			node := m.flatDocs[m.docCursor].Node
			m.state = viewDocumentDetail
			m.loading = true
			m.breadcrumb = append(m.breadcrumb, node.Title)
			return m, m.loadDocument(node.ID)
		}
	case "esc", "backspace":
		m.state = viewCollections
		m.breadcrumb = []string{"Collections"}
	}
	return m, nil
}

func (m Model) updateDocumentDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.docScroll > 0 {
			m.docScroll--
		}
	case "down", "j":
		if m.docScroll < len(m.docLines)-1 {
			m.docScroll++
		}
	case "esc", "backspace":
		m.state = viewDocuments
		if len(m.breadcrumb) > 2 {
			m.breadcrumb = m.breadcrumb[:len(m.breadcrumb)-1]
		}
	}
	return m, nil
}

func (m Model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.searchActive = false
		m.state = viewCollections
		m.breadcrumb = []string{"Collections"}
		return m, nil
	case "enter":
		if m.searchQuery != "" {
			m.loading = true
			return m, m.searchDocuments(m.searchQuery)
		}
	case "backspace":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
		}
	case "up":
		if m.searchCursor > 0 {
			m.searchCursor--
		}
	case "down":
		if m.searchCursor < len(m.searchResults)-1 {
			m.searchCursor++
		}
	case "tab":
		if len(m.searchResults) > 0 {
			result := m.searchResults[m.searchCursor]
			m.searchActive = false
			m.state = viewDocumentDetail
			m.loading = true
			m.breadcrumb = []string{"Search", result.Document.Title}
			return m, m.loadDocument(result.Document.ID)
		}
	default:
		if len(msg.String()) == 1 {
			m.searchQuery += msg.String()
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var content string

	// Breadcrumb
	bc := breadcrumbStyle.Render(strings.Join(m.breadcrumb, " > "))

	if m.loading {
		content = "\n  Loading..."
	} else if m.err != nil {
		content = errorStyle.Render(fmt.Sprintf("\n  Error: %s", m.err))
	} else {
		switch m.state {
		case viewCollections:
			content = m.viewCollections()
		case viewDocuments:
			content = m.viewDocuments()
		case viewDocumentDetail:
			content = m.viewDocumentDetail()
		case viewSearch:
			content = m.viewSearch()
		}
	}

	// Status bar
	help := helpStyle.Render("j/k: navigate  enter: select  esc: back  /: search  q: quit")
	statusBar := statusBarStyle.Width(m.width).Render(help)

	return lipgloss.JoinVertical(lipgloss.Left,
		bc,
		content,
		statusBar,
	)
}

func (m Model) viewCollections() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Collections"))
	b.WriteString("\n\n")
	for i, c := range m.collections {
		if i == m.collCursor {
			b.WriteString(selectedStyle.Render(fmt.Sprintf("  %s", c.Name)))
		} else {
			b.WriteString(normalStyle.Render(fmt.Sprintf("  %s", c.Name)))
		}
		b.WriteString(dimStyle.Render(fmt.Sprintf("  %s", c.ID[:8])))
		b.WriteString("\n")
	}
	if len(m.collections) == 0 {
		b.WriteString(dimStyle.Render("  No collections found"))
	}
	return b.String()
}

func (m Model) viewDocuments() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render(m.currentCollection.Name))
	b.WriteString("\n\n")

	visibleHeight := m.height - 6
	start := 0
	if m.docCursor >= visibleHeight {
		start = m.docCursor - visibleHeight + 1
	}

	for i := start; i < len(m.flatDocs) && i < start+visibleHeight; i++ {
		fd := m.flatDocs[i]
		indent := strings.Repeat("  ", fd.Depth)
		prefix := "  "
		if len(fd.Node.Children) > 0 {
			prefix = "  "
		}
		line := fmt.Sprintf("%s%s%s", indent, prefix, fd.Node.Title)
		if i == m.docCursor {
			b.WriteString(selectedStyle.Render(line))
		} else {
			b.WriteString(normalStyle.Render(line))
		}
		b.WriteString("\n")
	}
	if len(m.flatDocs) == 0 {
		b.WriteString(dimStyle.Render("  No documents found"))
	}
	return b.String()
}

func (m Model) viewDocumentDetail() string {
	if m.currentDoc == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(titleStyle.Render(m.currentDoc.Title))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render(fmt.Sprintf("Updated: %s  Rev: %d", m.currentDoc.UpdatedAt.Format("2006-01-02 15:04"), m.currentDoc.Revision)))
	b.WriteString("\n\n")

	visibleHeight := m.height - 8
	end := m.docScroll + visibleHeight
	if end > len(m.docLines) {
		end = len(m.docLines)
	}
	for i := m.docScroll; i < end; i++ {
		b.WriteString(docContentStyle.Render(m.docLines[i]))
		b.WriteString("\n")
	}
	return b.String()
}

func (m Model) viewSearch() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Search"))
	b.WriteString("\n\n")
	b.WriteString(searchInputStyle.Render(fmt.Sprintf("  > %s_", m.searchQuery)))
	b.WriteString("\n\n")

	for i, r := range m.searchResults {
		line := fmt.Sprintf("  %s", r.Document.Title)
		if i == m.searchCursor {
			b.WriteString(selectedStyle.Render(line))
		} else {
			b.WriteString(normalStyle.Render(line))
		}
		b.WriteString("\n")
	}
	if len(m.searchResults) == 0 && m.searchQuery != "" && !m.loading {
		b.WriteString(dimStyle.Render("  No results"))
	}
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("  enter: search  tab: open  esc: back"))
	return b.String()
}

func flattenNodes(nodes []models.NavigationNode, depth int) []flatDoc {
	var result []flatDoc
	for _, n := range nodes {
		result = append(result, flatDoc{Node: n, Depth: depth})
		if len(n.Children) > 0 {
			result = append(result, flattenNodes(n.Children, depth+1)...)
		}
	}
	return result
}
