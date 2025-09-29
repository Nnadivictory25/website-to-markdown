package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"website-markdown/internal/scraper"

	"github.com/spf13/cobra"
)

var (
	maxDepth       int
	delay          int
	followExternal bool
	output         string
	format         string
	userAgent      string
)

var rootCmd = &cobra.Command{
	Use:   "website-markdown [URL]",
	Short: "🔄 Convert websites to markdown recursively",
	Long: `A powerful tool to convert websites to markdown format.
Supports recursive scraping with configurable depth and delays.

Examples:
  website-markdown https://example.com
  website-markdown https://example.com --depth 2 --output ./docs
  website-markdown https://example.com --format json --external`,
	Args: cobra.ExactArgs(1),
	RunE: runScraper,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().IntVarP(&maxDepth, "depth", "d", 3, "Maximum depth for recursive scraping")
	rootCmd.Flags().IntVar(&delay, "delay", 1000, "Delay between requests in milliseconds")
	rootCmd.Flags().BoolVar(&followExternal, "external", false, "Follow external links")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Output directory (default: current directory)")
	rootCmd.Flags().StringVarP(&format, "format", "f", "files", "Output format: files, json, single")
	rootCmd.Flags().StringVar(&userAgent, "user-agent", "Website-Markdown-Converter/1.0", "User agent string")
}

func runScraper(cmd *cobra.Command, args []string) error {
	url := args[0]

	fmt.Printf("🚀 Starting website to markdown conversion\n")
	fmt.Printf("📍 URL: %s\n", url)
	fmt.Printf("📊 Max Depth: %d\n", maxDepth)
	fmt.Printf("⏱️  Delay: %dms\n", delay)
	fmt.Printf("🌐 Follow External: %t\n", followExternal)

	config := &scraper.ScrapingConfig{
		MaxDepth:       maxDepth,
		Delay:          time.Duration(delay) * time.Millisecond,
		FollowExternal: followExternal,
		UserAgent:      userAgent,
	}

	s := scraper.NewScraper(config)
	pages, err := s.ScrapeWebsite(url)
	if err != nil {
		return fmt.Errorf("❌ Scraping failed: %v", err)
	}

	if len(pages) == 0 {
		fmt.Println("⚠️  No pages were scraped")
		return nil
	}

	return saveOutput(pages, url)
}

func saveOutput(pages []*scraper.ScrapedPage, baseURL string) error {
	if output == "" {
		output = "."
	}

	switch format {
	case "json":
		return saveAsJSON(pages, baseURL)
	case "single":
		return saveAsSingleFile(pages, baseURL)
	default:
		return saveAsFiles(pages, baseURL)
	}
}

func saveAsJSON(pages []*scraper.ScrapedPage, baseURL string) error {
	filename := filepath.Join(output, generateFilename(baseURL, "json"))

	data, err := json.MarshalIndent(pages, "", "  ")
	if err != nil {
		return fmt.Errorf("❌ Failed to marshal JSON: %v", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("❌ Failed to write JSON file: %v", err)
	}

	fmt.Printf("💾 JSON output saved to: %s\n", filename)
	return nil
}

func saveAsSingleFile(pages []*scraper.ScrapedPage, baseURL string) error {
	filename := filepath.Join(output, generateFilename(baseURL, "md"))

	var content strings.Builder
	content.WriteString(fmt.Sprintf("# Website Content: %s\n\n", baseURL))
	content.WriteString(fmt.Sprintf("*Scraped on %s*\n\n", time.Now().Format("2006-01-02 15:04:05")))
	content.WriteString("---\n\n")

	for i, page := range pages {
		if page.Error != "" {
			content.WriteString(fmt.Sprintf("## ❌ Error: %s\n\n", page.URL))
			content.WriteString(fmt.Sprintf("**Error:** %s\n\n", page.Error))
			continue
		}

		content.WriteString(fmt.Sprintf("## 📄 Page %d: %s\n\n", i+1, page.Title))
		content.WriteString(fmt.Sprintf("**URL:** %s  \n", page.URL))
		content.WriteString(fmt.Sprintf("**Depth:** %d\n\n", page.Depth))
		content.WriteString("---\n\n")
		content.WriteString(page.Markdown)
		content.WriteString("\n\n")

		if i < len(pages)-1 {
			content.WriteString("---\n\n")
		}
	}

	err := os.WriteFile(filename, []byte(content.String()), 0644)
	if err != nil {
		return fmt.Errorf("❌ Failed to write markdown file: %v", err)
	}

	fmt.Printf("💾 Single markdown file saved to: %s\n", filename)
	return nil
}

func saveAsFiles(pages []*scraper.ScrapedPage, baseURL string) error {
	// Create output directory
	if err := os.MkdirAll(output, 0755); err != nil {
		return fmt.Errorf("❌ Failed to create output directory: %v", err)
	}

	successCount := 0
	errorCount := 0

	for i, page := range pages {
		if page.Error != "" {
			errorCount++
			fmt.Printf("⚠️  Error on %s: %s\n", page.URL, page.Error)
			continue
		}

		filename := fmt.Sprintf("page-%03d-%s.md", i+1, sanitizeFilename(page.Title))
		filepath := filepath.Join(output, filename)

		var content strings.Builder
		content.WriteString(fmt.Sprintf("# %s\n\n", page.Title))
		content.WriteString(fmt.Sprintf("**URL:** %s  \n", page.URL))
		content.WriteString(fmt.Sprintf("**Depth:** %d  \n", page.Depth))
		content.WriteString(fmt.Sprintf("**Scraped:** %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
		content.WriteString("---\n\n")
		content.WriteString(page.Markdown)

		err := os.WriteFile(filepath, []byte(content.String()), 0644)
		if err != nil {
			fmt.Printf("⚠️  Failed to write %s: %v\n", filename, err)
			errorCount++
			continue
		}

		successCount++
	}

	fmt.Printf("✅ Successfully saved %d files to: %s\n", successCount, output)
	if errorCount > 0 {
		fmt.Printf("⚠️  %d pages had errors\n", errorCount)
	}

	return nil
}

func generateFilename(websiteURL, extension string) string {
	parsedURL, err := url.Parse(websiteURL)
	if err != nil {
		// Fallback if URL parsing fails
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		return fmt.Sprintf("website_%s.%s", timestamp, extension)
	}

	// Get hostname and remove www. prefix
	siteName := parsedURL.Hostname()
	siteName = strings.TrimPrefix(siteName, "www.")

	// Clean up the site name to be filesystem-safe
	siteName = strings.ReplaceAll(siteName, ".", "-")
	siteName = sanitizeFilename(siteName)

	// Generate timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")

	return fmt.Sprintf("%s_%s.%s", siteName, timestamp, extension)
}

func sanitizeFilename(title string) string {
	// Remove/replace invalid filename characters
	title = strings.ReplaceAll(title, "/", "-")
	title = strings.ReplaceAll(title, "\\", "-")
	title = strings.ReplaceAll(title, ":", "-")
	title = strings.ReplaceAll(title, "*", "-")
	title = strings.ReplaceAll(title, "?", "-")
	title = strings.ReplaceAll(title, "\"", "-")
	title = strings.ReplaceAll(title, "<", "-")
	title = strings.ReplaceAll(title, ">", "-")
	title = strings.ReplaceAll(title, "|", "-")

	// Truncate if too long
	if len(title) > 50 {
		title = title[:50]
	}

	// Remove multiple dashes and trim
	title = strings.ReplaceAll(title, "--", "-")
	title = strings.Trim(title, "-")

	if title == "" {
		title = "untitled"
	}

	return title
}
