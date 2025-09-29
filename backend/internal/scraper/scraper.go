package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

type ScrapingConfig struct {
	MaxDepth       int           `json:"maxDepth"`
	Delay          time.Duration `json:"delay"`
	FollowExternal bool          `json:"followExternal"`
	UserAgent      string        `json:"userAgent"`
}

type ScrapedPage struct {
	URL      string `json:"url"`
	Title    string `json:"title"`
	Markdown string `json:"markdown"`
	Depth    int    `json:"depth"`
	Error    string `json:"error,omitempty"`
}

type Scraper struct {
	config         ScrapingConfig
	visited        map[string]bool
	baseHost       string
	converter      *md.Converter
	client         *http.Client
	duplicateCount int
}

const (
	DEFAULT_MAX_DEPTH  = 3
	DEFAULT_DELAY      = time.Second * 1
	DEFAULT_USER_AGENT = "Website-Markdown-Converter/1.0"
)

func NewScraper(config *ScrapingConfig) *Scraper {
	if config == nil {
		config = &ScrapingConfig{
			MaxDepth:       DEFAULT_MAX_DEPTH,
			Delay:          DEFAULT_DELAY,
			FollowExternal: false,
			UserAgent:      DEFAULT_USER_AGENT,
		}
	}

	converter := md.NewConverter("", true, nil)

	return &Scraper{
		config:         *config,
		visited:        make(map[string]bool),
		converter:      converter,
		duplicateCount: 0,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *Scraper) ScrapeWebsite(startURL string) ([]*ScrapedPage, error) {
	parsedURL, err := url.Parse(startURL)
	if err != nil {
		return nil, fmt.Errorf("🚫 Invalid URL: %v", err)
	}

	s.baseHost = parsedURL.Host
	var results []*ScrapedPage

	// Normalize the starting URL
	normalizedStartURL := s.normalizeURL(startURL)
	fmt.Printf("🚀 Starting recursive scrape of %s (max depth: %d)\n", normalizedStartURL, s.config.MaxDepth)

	s.scrapeRecursive(normalizedStartURL, 0, &results)

	if s.duplicateCount > 0 {
		fmt.Printf("✅ Scraping completed! Found %d unique pages (skipped %d duplicates)\n", len(results), s.duplicateCount)
	} else {
		fmt.Printf("✅ Scraping completed! Found %d pages\n", len(results))
	}
	return results, nil
}

func (s *Scraper) scrapeRecursive(pageURL string, depth int, results *[]*ScrapedPage) {
	// Normalize URL to prevent duplicates
	normalizedURL := s.normalizeURL(pageURL)

	// Check if already visited or max depth reached
	if s.visited[normalizedURL] {
		s.duplicateCount++
		if pageURL != normalizedURL {
			fmt.Printf("⏭️  Skipping duplicate (normalized): %s -> %s\n", pageURL, normalizedURL)
		} else {
			fmt.Printf("⏭️  Skipping already visited: %s\n", normalizedURL)
		}
		return
	}

	if depth > s.config.MaxDepth {
		fmt.Printf("🛑 Max depth reached (%d), skipping: %s\n", s.config.MaxDepth, normalizedURL)
		return
	}

	s.visited[normalizedURL] = true

	if pageURL != normalizedURL {
		fmt.Printf("📄 Scraping (depth %d): %s (normalized from %s)\n", depth, normalizedURL, pageURL)
	} else {
		fmt.Printf("📄 Scraping (depth %d): %s\n", depth, normalizedURL)
	}

	// Add delay between requests
	if depth > 0 {
		time.Sleep(s.config.Delay)
	}

	page, links := s.scrapePage(normalizedURL, depth)
	*results = append(*results, page)

	// Recursively scrape found links
	for _, link := range links {
		normalizedLink := s.normalizeURL(link)
		if !s.visited[normalizedLink] {
			s.scrapeRecursive(normalizedLink, depth+1, results)
		}
	}
}

func (s *Scraper) scrapePage(pageURL string, depth int) (*ScrapedPage, []string) {
	page := &ScrapedPage{
		URL:   pageURL,
		Depth: depth,
	}

	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		page.Error = fmt.Sprintf("Failed to create request: %v", err)
		return page, nil
	}

	req.Header.Set("User-Agent", s.config.UserAgent)

	resp, err := s.client.Do(req)
	if err != nil {
		page.Error = fmt.Sprintf("Failed to fetch page: %v", err)
		return page, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		page.Error = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)
		return page, nil
	}

	// Check if content is HTML
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		page.Error = "Not an HTML page"
		return page, nil
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		page.Error = fmt.Sprintf("Failed to parse HTML: %v", err)
		return page, nil
	}

	// Extract title
	page.Title = doc.Find("title").First().Text()
	if page.Title == "" {
		page.Title = pageURL
	}

	// Convert to markdown
	html, _ := doc.Html()
	markdown, err := s.converter.ConvertString(html)
	if err != nil {
		page.Error = fmt.Sprintf("Failed to convert to markdown: %v", err)
		return page, nil
	}

	page.Markdown = s.cleanMarkdown(markdown)

	// Extract links for recursive scraping
	var links []string
	if depth < s.config.MaxDepth {
		links = s.extractLinks(doc, pageURL)
	}

	return page, links
}

func (s *Scraper) extractLinks(doc *goquery.Document, baseURL string) []string {
	var links []string
	seenLinks := make(map[string]bool)

	doc.Find("a[href]").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if !exists || href == "" {
			return
		}

		// Resolve relative URLs
		parsedBase, _ := url.Parse(baseURL)
		parsedHref, err := url.Parse(href)
		if err != nil {
			return
		}

		resolvedURL := parsedBase.ResolveReference(parsedHref)
		finalURL := resolvedURL.String()

		// Skip non-HTTP(S) URLs
		if !strings.HasPrefix(finalURL, "http://") && !strings.HasPrefix(finalURL, "https://") {
			return
		}

		// Skip if external and not following external links
		if !s.config.FollowExternal && resolvedURL.Host != s.baseHost {
			return
		}

		// Normalize URL to prevent duplicates
		normalizedURL := s.normalizeURL(finalURL)

		// Skip if already seen (using normalized URL)
		if seenLinks[normalizedURL] {
			return
		}
		seenLinks[normalizedURL] = true

		// Skip common file extensions
		if s.isFileLink(normalizedURL) {
			return
		}

		links = append(links, normalizedURL)
	})

	return links
}

func (s *Scraper) isFileLink(url string) bool {
	fileExtensions := []string{
		".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
		".zip", ".rar", ".tar", ".gz", ".7z",
		".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg",
		".mp4", ".avi", ".mov", ".wmv", ".flv",
		".mp3", ".wav", ".flac", ".ogg",
		".css", ".js", ".json", ".xml", ".rss",
	}

	lowerURL := strings.ToLower(url)
	for _, ext := range fileExtensions {
		if strings.HasSuffix(lowerURL, ext) {
			return true
		}
	}
	return false
}

func (s *Scraper) normalizeURL(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL // Return original if can't parse
	}

	// Convert scheme and host to lowercase
	parsedURL.Scheme = strings.ToLower(parsedURL.Scheme)
	parsedURL.Host = strings.ToLower(parsedURL.Host)

	// Remove fragment (everything after #)
	parsedURL.Fragment = ""

	// Remove trailing slash from path (except for root)
	if parsedURL.Path != "/" && strings.HasSuffix(parsedURL.Path, "/") {
		parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/")
	}

	// Remove common tracking parameters
	if parsedURL.RawQuery != "" {
		values := parsedURL.Query()
		trackingParams := []string{
			"utm_source", "utm_medium", "utm_campaign", "utm_term", "utm_content",
			"gclid", "fbclid", "ref", "source", "from",
			"_ga", "_gl", "mc_cid", "mc_eid",
		}

		for _, param := range trackingParams {
			values.Del(param)
		}

		// If no query params left, clear the RawQuery
		if len(values) == 0 {
			parsedURL.RawQuery = ""
		} else {
			parsedURL.RawQuery = values.Encode()
		}
	}

	return parsedURL.String()
}

func (s *Scraper) cleanMarkdown(markdown string) string {
	// Remove excessive newlines
	re := regexp.MustCompile(`\n{3,}`)
	markdown = re.ReplaceAllString(markdown, "\n\n")

	// Trim whitespace
	markdown = strings.TrimSpace(markdown)

	return markdown
}
