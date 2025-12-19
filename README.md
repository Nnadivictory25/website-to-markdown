# 🔄 Website to Markdown Converter

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)
[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org/)
[![Svelte](https://img.shields.io/badge/Svelte-5-orange.svg)](https://svelte.dev/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind-4-blue.svg)](https://tailwindcss.com/)
[![Bun](https://img.shields.io/badge/Bun-Latest-black.svg)](https://bun.sh/)

A powerful Go-based tool that recursively converts websites to clean markdown format with **intelligent duplicate prevention**. Features both a **CLI tool** for terminal usage and a **REST API** with a beautiful **Svelte + Tailwind CSS** web interface.

> 🎯 **Perfect for**: Documentation archiving, SEO audits, content migration, AI training data, and research projects!

## 📋 Table of Contents

- [✨ Features](#-features)
- [🚀 Quick Start](#-quick-start)
- [🖥️ CLI Usage](#️-cli-usage)
- [🌐 Web Interface](#-web-interface)
- [📁 Output Formats](#-output-formats)
- [🧠 Intelligent Duplicate Prevention](#-intelligent-duplicate-prevention)
- [🔧 API Reference](#-api-reference)
- [⚙️ Configuration](#️-configuration)
- [🛠️ Development](#️-development)
- [📝 Real-World Examples](#-real-world-examples)
 - [🤝 Contributing](#-contributing)
 - [📄 License](#-license)

## ✨ Features

### 🔧 Core Functionality
- **🌐 Recursive Scraping**: Automatically follow links to scrape entire website sections
- **🚫 Duplicate Prevention**: Smart URL normalization prevents infinite loops and duplicate pages
- **📊 Configurable Depth**: Control scraping depth (1-10 levels) with intelligent limits
- **⏱️ Rate Limiting**: Respectful delays between requests (100ms-3000ms)
- **🌍 External Link Support**: Option to follow external links or stay within domain
- **🎯 Smart Filtering**: Auto-skips non-HTML content, files, and minimal pages

### 📄 Output Options
- **📁 Individual Files**: Separate markdown files for each page
- **📄 Single Combined File**: All pages in one markdown document
- **📦 JSON Export**: Structured data format for programmatic use
- **💾 Direct Downloads**: API endpoint returns downloadable .md files
- **🏷️ Smart Naming**: Files named with website + timestamp

### 🖥️ Interfaces
- **🖥️ CLI Tool**: Full-featured command-line interface
- **🌐 REST API**: CORS-enabled API with multiple endpoints
- **📱 Modern Web UI**: Responsive Svelte frontend with real-time progress
- **💾 One-Click Downloads**: Direct download buttons for all formats

## 🚀 Quick Start

### Prerequisites
- **Go 1.19+** for the backend
- **Bun** for the frontend (faster than npm/yarn)
- Modern web browser for the UI

### Installation

1. **Clone and build CLI**
   ```bash
   git clone <repository-url>
   cd website-markdown/backend
   go build -o website-markdown main.go
   ```

2. **Install frontend dependencies**
   ```bash
   cd ../frontend
   bun install
   ```

### Quick Test

```bash
# Test CLI
./website-markdown https://example.com --depth 2

# Start full web interface
cd backend && ./website-markdown --server &
cd frontend && bun run dev
# Visit http://localhost:5173
```

## 🖥️ CLI Usage

### Basic Commands

```bash
# Convert with defaults (3 levels deep)
./website-markdown https://example.com

# Custom depth and output directory
./website-markdown https://example.com --depth 2 --output ./docs

# Follow external links with custom delay
./website-markdown https://example.com --external --delay 2000

# Output formats
./website-markdown https://example.com --format single  # Single .md file
./website-markdown https://example.com --format json    # JSON data
```

### CLI Options

| Flag           | Description                          | Default                        |
| -------------- | ------------------------------------ | ------------------------------ |
| `--depth, -d`  | Maximum scraping depth (1-10)        | 3                              |
| `--delay`      | Delay between requests (ms)          | 1000                           |
| `--external`   | Follow external links                | false                          |
| `--output, -o` | Output directory                     | current directory              |
| `--format, -f` | Output format: `files`, `single`, `json` | files                          |
| `--user-agent` | Custom User-Agent string             | Website-Markdown-Converter/1.0 |

### CLI Examples

```bash
# Scrape a blog with 2 levels depth
./website-markdown https://blog.example.com --depth 2 --output ./blog-content

# Scrape documentation site into single file
./website-markdown https://docs.example.com --format single --output ./docs

# Respectful scraping with longer delays
./website-markdown https://example.com --delay 3000 --depth 1
```

## 🌐 Web Interface

### Starting the Interface

**Option 1: Separate terminals**
```bash
# Terminal 1: API Server
cd backend && ./website-markdown --server

# Terminal 2: Frontend
cd frontend && bun run dev
```

**Option 2: Background mode**
```bash
cd backend && ./website-markdown --server &
cd frontend && bun run dev
# Visit http://localhost:5173
```

### 🎨 Features
- **🌐 Smart URL Input**: Validation with error handling
- **⚙️ Visual Controls**: Depth slider, delay selector, external links toggle
- **📊 Real-time Stats**: Live progress with success/error counts and processing time
- **💾 Direct Downloads**: One-click downloads for markdown and JSON formats
- **📄 Live Preview**: Expandable content preview
- **📱 Responsive Design**: Works on mobile and desktop
- **🎯 Modern Stack**: Svelte 5 + Tailwind CSS 4 + TypeScript

## 📁 Output Formats

All outputs use **smart naming** with website + timestamp (e.g., `github-com_2024-01-15_14-30-25.md`).

### CLI Formats

#### Individual Files (`--format files`)
Creates a directory with separate markdown files for each page:
```
github-com_2024-01-15_14-30-25/
├── page-001-GitHub-Homepage.md
├── page-002-About-GitHub.md
└── page-003-Features.md
```

#### Single Combined File (`--format single`)
One comprehensive markdown file with all pages:
```markdown
# Website Content: https://github.com

_Scraped on 2024-01-15 14:30:25_
_Found 15 unique pages_

---

## 📄 Page 1: GitHub Homepage
**URL:** https://github.com | **Depth:** 0

[Content...]

## 📄 Page 2: About GitHub
**URL:** https://github.com/about | **Depth:** 1

[Content...]
```

#### JSON Export (`--format json`)
Structured data file for programmatic use:
```json
[
  {
    "url": "https://github.com",
    "title": "GitHub Homepage",
    "markdown": "[converted content]",
    "depth": 0,
    "error": ""
  }
]
```

### API Endpoints

#### POST `/scrape`
Returns JSON response with pages array and stats.

#### GET `/download/markdown`
Returns downloadable `.md` file with combined content and table of contents. Query parameters: `url` (required), `depth`, `delay`, `external`.

## 🧠 Intelligent Duplicate Prevention

### 🚫 **How It Works**

The scraper now includes smart URL normalization to prevent infinite loops and duplicate pages:

**URL Normalization:**

- ✅ `https://Example.com/Page/` → `https://example.com/page`
- ✅ Removes tracking parameters (`utm_*`, `fbclid`, `gclid`, etc.)
- ✅ Removes URL fragments (everything after `#`)
- ✅ Handles trailing slashes consistently

**Console Feedback:**

```bash
📄 Scraping (depth 1): https://github.com/features
⏭️  Skipping duplicate (normalized): https://github.com/features/ → https://github.com/features
⏭️  Skipping already visited: https://github.com/about
🛑 Max depth reached (3), skipping: https://github.com/deep/page
✅ Scraping completed! Found 15 unique pages (skipped 8 duplicates)
```

### 🎯 **Benefits**

- 🚀 **Faster scraping** - No wasted time on duplicates
- 🔄 **No infinite loops** - Smart cycle detection
- 📊 **Clean data** - Only unique content
- 📝 **Transparent logging** - See exactly what's happening

## 🔧 API Reference

### POST `/scrape`
Scrape a website and return JSON results.

**Request:**
```json
{
  "url": "https://example.com",
  "maxDepth": 3,
  "delay": 1000,
  "followExternal": false
}
```

**Response:**
```json
{
  "success": true,
  "message": "Successfully scraped 5 pages",
  "pages": [...],
  "stats": {
    "totalPages": 5,
    "successPages": 4,
    "errorPages": 1,
    "processingTime": "15.2s",
    "startedAt": "2024-01-15T10:30:00Z",
    "completedAt": "2024-01-15T10:30:15Z"
  }
}
```

### GET `/download/markdown`
Download scraped website content as markdown file with table of contents.

**Request:** `GET /download/markdown?url=https://example.com&depth=3&delay=1000&external=false`

**Parameters:**
- `url` (required): Website URL to scrape
- `depth` (optional): Maximum scraping depth (1-10, default: 3)
- `delay` (optional): Delay between requests in ms (default: 1000, min: 100)
- `external` (optional): Follow external links (default: false)

**⚡ Speed Note:** Minimum 100ms delay enforced for respectful scraping. Use 100-500ms for fast but responsible scraping.

**Response:** Markdown file with table of contents, statistics, and all scraped content. Includes anchor links for easy navigation.

### GET `/status`
Get server status and available endpoints.

### GET `/health`
Health check endpoint.

## ⚙️ Configuration

### Environment Variables
```bash
export PORT=8080                    # API server port
export DEFAULT_MAX_DEPTH=3          # Default scraping depth
export DEFAULT_DELAY=1000           # Default delay between requests
export DEFAULT_USER_AGENT="Website-Markdown-Converter/1.0"
```

### Respectful Scraping
- **⏱️ Configurable delays** (100ms-3000ms) between requests
- **🤖 Proper User-Agent** identification
- **🚫 Smart filtering** of non-HTML content, files, and minimal pages
- **📝 robots.txt respect** (planned feature)
- **🔒 Built-in rate limiting** to prevent abuse

## 🛠️ Development

### Backend (Go)
```bash
cd backend
go mod tidy                          # Install dependencies
go run main.go --help               # Show CLI help
go run main.go --server             # Start API server
go build -o website-markdown main.go # Build binary

# Live reload (optional)
go install github.com/cosmtrek/air@latest
air  # Auto-restart on changes
```

### Frontend (Svelte + Bun)
```bash
cd frontend
bun install     # Install dependencies
bun run dev     # Start dev server
```

### Project Structure
```
website-markdown/
├── backend/              # Go API & CLI
│   ├── cmd/cli.go       # CLI commands & file operations
│   ├── internal/
│   │   ├── api/server.go    # HTTP endpoints
│   │   └── scraper/scraper.go # Core scraping logic
│   └── main.go          # Entry point
├── frontend/             # Svelte UI
│   ├── src/routes/+page.svelte # Main UI
│   └── package.json     # Dependencies
├── README.md             # Documentation
└── LICENSE               # MIT License
```

## 📝 Real-World Examples

### 🖥️ **CLI Examples**

**📚 Documentation Sites:**

```bash
# Scrape entire documentation section
./website-markdown https://docs.svelte.dev --depth 3 --format single
# Output: docs-svelte-dev_2024-01-15_14-30-25.md

# Multiple formats for Go docs
./website-markdown https://pkg.go.dev/net/http --depth 2 --format json
# Output: pkg-go-dev_2024-01-15_14-30-25.json
```

**📰 News & Blogs:**

```bash
# Respectful blog scraping
./website-markdown https://blog.golang.org --depth 2 --delay 2000 --external
# Gets main posts + linked articles

# News aggregator (single level)
./website-markdown https://news.ycombinator.com --depth 1 --delay 3000 --format single
```

**🏢 Company Sites:**

```bash
# Marketing site audit
./website-markdown https://stripe.com --depth 2 --output ./stripe-audit/
# Creates folder: stripe-audit/ with organized pages
```

### 🌐 **API Examples**

```bash
# Download markdown file directly
curl "http://localhost:8080/download/markdown?url=https://docs.github.com&depth=2&delay=100" \
  --output docs.md

# JSON API response
curl -X POST http://localhost:8080/scrape \
  -H "Content-Type: application/json" \
  -d '{"url": "https://github.com/features", "maxDepth": 2, "delay": 1000}' \
  | jq '.stats.successPages'
```

## 🤝 Contributing

We welcome contributions from everyone! Whether you're fixing bugs, adding features, improving documentation, or sharing ideas.

### 🚀 **Quick Contribute**

1. 🍴 Fork the repository
2. 🌿 Create a feature branch (`git checkout -b feature/amazing-feature`)
3. ✨ Make your changes
4. 🧪 Test your changes
5. 📝 Commit with clear messages (`git commit -m '✨ Add amazing feature'`)
6. 🚀 Push to your branch (`git push origin feature/amazing-feature`)
7. 🎯 Open a Pull Request

### 📋 **Areas We Need Help With:**

- 🐛 **Bug fixes** - Duplicate detection, error handling
- ✨ **Features** - Robots.txt support, sitemap parsing, PDF output
- 🎨 **UI/UX** - Mobile improvements, dark mode, better progress
- 📚 **Documentation** - More examples, translations
- 🧪 **Testing** - Edge cases, performance testing

See our [Contributing Guide](./CONTRIBUTING.md) for detailed information!

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](./LICENSE) file for details.

**TL;DR**: ✅ Free to use, modify, distribute, and even sell. Just keep the license notice!

## 🙏 Acknowledgments

- [goquery](https://github.com/PuerkitoBio/goquery) for HTML parsing
- [html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown) for conversion
- [Gin](https://github.com/gin-gonic/gin) for the web framework
- [Cobra](https://github.com/spf13/cobra) for CLI interface
- [Svelte](https://svelte.dev/) and [Tailwind CSS](https://tailwindcss.com/) for the frontend

---
