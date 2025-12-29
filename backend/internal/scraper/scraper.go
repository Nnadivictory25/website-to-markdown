package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

type ScrapingConfig struct {
	MaxDepth       int           `json:"maxDepth"`
	Delay          time.Duration `json:"delay"`
	FollowExternal bool          `json:"followExternal"`
	UserAgent      string        `json:"userAgent"`
	Concurrency    int           `json:"concurrency"`
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
	visitedMutex   sync.RWMutex
	baseHost       string
	converter      *md.Converter
	client         *http.Client
	duplicateCount int
}

const (
	DEFAULT_MAX_DEPTH   = 3
	DEFAULT_DELAY       = time.Second * 1
	DEFAULT_USER_AGENT  = "Website-Markdown-Converter/1.0"
	DEFAULT_CONCURRENCY = 5
)

func NewScraper(config *ScrapingConfig) *Scraper {
	if config == nil {
		config = &ScrapingConfig{
			MaxDepth:       DEFAULT_MAX_DEPTH,
			Delay:          DEFAULT_DELAY,
			FollowExternal: false,
			UserAgent:      DEFAULT_USER_AGENT,
			Concurrency:    DEFAULT_CONCURRENCY,
		}
	}

	if config.Concurrency <= 0 {
		config.Concurrency = DEFAULT_CONCURRENCY
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

	// Normalize the starting URL
	normalizedStartURL := s.normalizeURL(startURL)
	fmt.Printf("🚀 Starting level-based scrape of %s (max depth: %d, concurrency: %d)\n", normalizedStartURL, s.config.MaxDepth, s.config.Concurrency)

	results := s.scrapeLevelBFS(normalizedStartURL)

	if s.duplicateCount > 0 {
		fmt.Printf("✅ Scraping completed! Found %d unique pages (skipped %d duplicates)\n", len(results), s.duplicateCount)
	} else {
		fmt.Printf("✅ Scraping completed! Found %d pages\n", len(results))
	}
	return results, nil
}

func (s *Scraper) scrapeLevelBFS(startURL string) []*ScrapedPage {
	var results []*ScrapedPage

	// Start with initial URL
	currentLevel := []string{startURL}
	s.visited[startURL] = true

	// Process each level (depth)
	for depth := 0; depth <= s.config.MaxDepth && len(currentLevel) > 0; depth++ {
		fmt.Printf("📍 Processing depth %d (%d pages)...\n", depth, len(currentLevel))

		// Phase 1: Scrape full content for all pages at current level
		pagesAtThisLevel := s.scrapePagesConcurrent(currentLevel, depth)
		results = append(results, pagesAtThisLevel...)

		// Phase 2: Discover all links from pages at this level
		nextLevelLinks := s.discoverLinksConcurrent(currentLevel)

		// Filter: only keep unvisited links for next level
		currentLevel = s.filterUnvisited(nextLevelLinks)

		if len(currentLevel) > 0 {
			fmt.Printf("🔗 Found %d new pages for depth %d\n", len(currentLevel), depth+1)
		}
	}

	return results
}

func (s *Scraper) scrapePagesConcurrent(urls []string, depth int) []*ScrapedPage {
	var wg sync.WaitGroup
	pagesChan := make(chan *ScrapedPage, len(urls))

	for _, u := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			fmt.Printf("📄 Scraping (depth %d): %s\n", depth, url)
			page, _ := s.scrapePage(url, depth)
			if page != nil {
				pagesChan <- page
			}
		}(u)
	}

	go func() {
		wg.Wait()
		close(pagesChan)
	}()

	var pages []*ScrapedPage
	for page := range pagesChan {
		pages = append(pages, page)
	}

	return pages
}

func (s *Scraper) discoverLinksConcurrent(urls []string) []string {
	var wg sync.WaitGroup
	linksChan := make(chan []string, len(urls))

	for _, u := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			_, links := s.scrapePage(url, 0)
			linksChan <- links
		}(u)
	}

	go func() {
		wg.Wait()
		close(linksChan)
	}()

	var allLinks []string
	seen := make(map[string]bool)

	for links := range linksChan {
		for _, link := range links {
			normalized := s.normalizeURL(link)
			if !seen[normalized] {
				seen[normalized] = true
				allLinks = append(allLinks, normalized)
			}
		}
	}

	return allLinks
}

func (s *Scraper) filterUnvisited(urls []string) []string {
	s.visitedMutex.Lock()
	defer s.visitedMutex.Unlock()

	var unvisited []string
	for _, url := range urls {
		if !s.visited[url] {
			s.visited[url] = true
			unvisited = append(unvisited, url)
		} else {
			s.duplicateCount++
		}
	}

	return unvisited
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

	// Filter out pages with minimal or generic content
	if s.isContentMinimal(page.Title, page.Markdown, pageURL) {
		fmt.Printf("⏭️  Skipping page with minimal content: %s\n", pageURL)
		return nil, nil // Return nil to skip this page
	}

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

func (s *Scraper) isContentMinimal(title, markdown, url string) bool {
	// Skip if title is just the URL (no real title)
	if title == url || strings.TrimSpace(title) == "" {
		return true
	}

	// Skip if content is very short (less than 200 characters)
	if len(strings.TrimSpace(markdown)) < 200 {
		return true
	}

	// Skip if content contains only navigation/header elements without real content
	lines := strings.Split(markdown, "\n")
	contentLines := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Count lines that aren't just navigation, headers, or empty
		if line != "" && !strings.HasPrefix(line, "#") &&
			!strings.Contains(line, "Skip to main content") &&
			!strings.Contains(line, "[") && !strings.Contains(line, "](") {
			contentLines++
		}
	}

	// Skip if less than 3 lines of actual content
	if contentLines < 3 {
		return true
	}

	return false
}

func (s *Scraper) cleanMarkdown(markdown string) string {
	// Remove excessive newlines
	re := regexp.MustCompile(`\n{3,}`)
	markdown = re.ReplaceAllString(markdown, "\n\n")

	// Trim whitespace
	markdown = strings.TrimSpace(markdown)

	return markdown
}
