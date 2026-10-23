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

// DocsPageModel is everything DocsPage needs. Body is trusted HTML
// rendered from repository Markdown.
type DocsPageModel struct {
	Title  string
	Body   string
	Crumbs []DocsLink // breadcrumb trail; the last item is the current page
	Nav    []DocsNavSection
	Prev   *DocsLink
	Next   *DocsLink
}

// DocsHomeSection is one section card on the docs home page.
type DocsHomeSection struct {
	DocsLink
	Description string
}

// DocsHomeModel is everything DocsHome needs.
type DocsHomeModel struct {
	Title    string
	Body     string
	Sections []DocsHomeSection
}
