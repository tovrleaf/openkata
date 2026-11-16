package templates

// DocsLink is a titled link to another docs page.
type DocsLink struct {
	Title string
	URL   string
}

// DocsNavPage is one entry in the docs sidebar.
type DocsNavPage struct {
	DocsLink
	Current bool
}

// DocsNavSection is a section and its pages in the docs sidebar.
type DocsNavSection struct {
	DocsNavPage
	Pages []DocsNavPage
}

// DocsUpdated is a page's last-updated date. ISO is empty when the date
// is unknown, and the footer then shows a placeholder.
type DocsUpdated struct {
	ISO   string // 2026-10-25
	Label string // 25 October 2026
}

// DocsPageModel is everything DocsPage needs. Body is trusted HTML
// rendered from repository Markdown.
type DocsPageModel struct {
	Title   string
	Body    string
	Crumbs  []DocsLink // breadcrumb trail; the last item is the current page
	Nav     []DocsNavSection
	Prev    *DocsLink
	Next    *DocsLink
	Pages   []DocsCard // the section's pages, on a section index only
	Updated DocsUpdated
}

// DocsCard is one card in a list of sections or pages.
type DocsCard struct {
	DocsLink
	Description string
}

// DocsHomeModel is everything DocsHome needs.
type DocsHomeModel struct {
	Title    string
	Body     string
	Sections []DocsCard
	Updated  DocsUpdated
}
