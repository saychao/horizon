package hal

import (
	"net/url"
)

// BasePage represents the simplest page: one with no links and only embedded records.
// Can be used to build custom page-like resources
type BasePage struct {
	BaseURL  *url.URL `json:"-"`
	Embedded struct {
		Meta    *PageMeta  `json:"meta,omitempty"`
		Records []Pageable `json:"records"`
	} `json:"_embedded"`
}

// Add appends the provided record onto the page
func (p *BasePage) Add(rec Pageable) {
	p.Embedded.Records = append(p.Embedded.Records, rec)
}

// Init initialized the Records slice.  This ensures that an empty page
// renders its records as an empty array, rather than `null`
func (p *BasePage) Init() {
	if p.Embedded.Records == nil {
		p.Embedded.Records = make([]Pageable, 0, 1)
	}
}

// Page represents the common page configuration (i.e. has self, next, and prev
// links) and has a helper method `PopulateLinks` to automate their
// initialization.
type Page struct {
	Links struct {
		Self Link `json:"self"`
		Next Link `json:"next"`
		Prev Link `json:"prev"`
	} `json:"_links"`

	BasePage
	BasePath string            `json:"-"`
	Order    string            `json:"-"`
	Limit    uint64            `json:"-"`
	Cursor   string            `json:"-"`
	Page     uint64            `json:"-"`
	Filters  map[string]string `json:"-"`
}

// PopulateLinks sets the common links for a page.
func (p *Page) PopulateLinks() {
	p.Init()
	fmts := p.BasePath + "?order=%s&limit=%d&cursor=%s"
	lb := LinkBuilder{p.BaseURL}

	p.Links.Self = lb.Linkf(p.Filters, fmts, p.Order, p.Limit, p.Cursor)
	rec := p.Embedded.Records

	if len(rec) > 0 {
		p.Links.Next = lb.Linkf(p.Filters, fmts, p.Order, p.Limit, rec[len(rec)-1].PagingToken())
		p.Links.Prev = lb.Linkf(p.Filters, fmts, p.InvertedOrder(), p.Limit, rec[0].PagingToken())
	} else {
		p.Links.Next = lb.Linkf(p.Filters, fmts, p.Order, p.Limit, p.Cursor)
		p.Links.Prev = lb.Linkf(p.Filters, fmts, p.InvertedOrder(), p.Limit, p.Cursor)
	}
}

// PopulateLinksV2 sets the common links for the page for limit/offset-based pagination
func (p *Page) PopulateLinksV2() {
	p.Init()
	fmts := p.BasePath + "?page=%d&limit=%d"

	lb := LinkBuilder{p.BaseURL}

	p.Links.Self = lb.Linkf(p.Filters, fmts, p.Page, p.Limit)
	rec := p.Embedded.Records

	if len(rec) > 0 {
		p.Links.Next = lb.Linkf(p.Filters, fmts, p.Page+1, p.Limit)
		if p.Page > 0 {
			p.Links.Prev = lb.Linkf(p.Filters, fmts, p.Page-1, p.Limit)
		}
	} else {
		p.Links.Next = lb.Linkf(p.Filters, fmts, p.Page, p.Limit)
		p.Links.Prev = lb.Linkf(p.Filters, fmts, p.Page, p.Limit)
	}
}

// InvertedOrder returns the inversion of the page's current order. Used to
// populate the prev link
func (p *Page) InvertedOrder() string {
	switch p.Order {
	case "asc":
		return "desc"
	case "desc":
		return "asc"
	default:
		return "asc"
	}
}
