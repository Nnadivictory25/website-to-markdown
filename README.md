# 🔄 Website to Markdown Converter

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)
[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org/)
[![Svelte](https://img.shields.io/badge/Svelte-5-orange.svg)](https://svelte.dev/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind-4-blue.svg)](https://tailwindcss.com/)
[![Bun](https://img.shields.io/badge/Bun-Latest-black.svg)](https://bun.sh/)

A powerful Go-based tool that recursively converts websites to clean markdown format with **intelligent duplicate prevention**. Features both a **CLI tool** for terminal usage and a **REST API** with a beautiful **Svelte + Tailwind CSS** web interface.

> 🎯 **Perfect for**: Documentation archiving, SEO audits, content migration, AI training data, and research projects!

## ✨ Features

### 🔧 **Core Functionality**

- 🌐 **Recursive Scraping**: Automatically follow links to scrape entire website sections
- 🚫 **Duplicate Prevention**: Smart URL normalization prevents infinite loops and duplicate pages
- 📊 **Configurable Depth**: Control scraping depth (1-10 levels) with intelligent limits
- ⏱️ **Rate Limiting**: Respectful delays between requests (500ms-3000ms)
- 🌍 **External Link Support**: Option to follow external links or stay within domain
- 🎯 **Smart Filtering**: Auto-skips non-HTML content, files, and tracking parameters

### 📄 **Output Options**

- 📁 **Individual Files**: Separate markdown files for each page
- 📄 **Single File**: Combined markdown with all pages
- 📦 **JSON Export**: Structured data format for programmatic use
- 🏷️ **Smart Naming**: Files named with website + timestamp (e.g., `github-com_2024-01-15_14-30-25.md`)

### 🖥️ **Dual Interface**

- 🖥️ **CLI Tool**: Full-featured command-line interface
- 🌐 **Web API**: RESTful API with CORS support
- 📱 **Modern UI**: Responsive Svelte frontend with real-time progress
- 💾 **One-Click Downloads**: Direct download buttons for both formats

## 🚀 Quick Start

### Prerequisites

- **Go 1.19+** for the backend
- **Bun** for the frontend (faster than npm/yarn)
- Modern web browser for the UI

### Installation

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd website-markdown
   ```

2. **Build the CLI tool**

   ```bash
   cd backend
   go build -o website-markdown main.go
   ```

3. **Install frontend dependencies**
   ```bash
   cd ../frontend
   bun install
   ```

### Quick Test

```bash
# Test CLI
cd backend
./website-markdown https://example.com --depth 2

# Start web interface
cd backend && ./website-markdown --server &  # Start API
cd frontend && bun run dev                   # Start UI
# Visit http://localhost:5173
```

## 🖥️ CLI Usage

### Basic Usage

```bash
# Convert a website with default settings (3 levels deep)
./website-markdown https://example.com

# Specify custom depth and output directory
./website-markdown https://example.com --depth 2 --output ./docs

# Follow external links and use custom delay
./website-markdown https://example.com --external --delay 2000

# Output as single markdown file
./website-markdown https://example.com --format single

# Output as JSON
./website-markdown https://example.com --format json
```

### CLI Options

| Flag           | Description                                 | Default                        |
| -------------- | ------------------------------------------- | ------------------------------ |
| `--depth, -d`  | Maximum depth for recursive scraping        | 3                              |
| `--delay`      | Delay between requests in milliseconds      | 1000                           |
| `--external`   | Follow external links                       | false                          |
| `--output, -o` | Output directory                            | current directory              |
| `--format, -f` | Output format: `files`, `single`, or `json` | files                          |
| `--user-agent` | Custom User-Agent string                    | Website-Markdown-Converter/1.0 |

### CLI Examples

```bash
# Scrape a blog with 2 levels depth
./website-markdown https://blog.example.com --depth 2 --output ./blog-content

# Scrape documentation site into single file
./website-markdown https://docs.example.com --format single --output ./docs

# Respectful scraping with longer delays
./website-markdown https://example.com --delay 3000 --depth 1
```

## 🌐 Web Interface Usage

### Start Both Servers

**Option 1: Separate terminals**

```bash
# Terminal 1: API Server
cd backend
./website-markdown --server --port 8080

# Terminal 2: Frontend (Svelte + Vite)
cd frontend
bun run dev
```

**Option 2: Background API**

```bash
cd backend && ./website-markdown --server &
cd frontend && bun run dev
# Visit http://localhost:5173
```

### 🎨 **Modern Svelte UI Features**

- 🌐 **URL Input**: Smart validation with error handling
- ⚙️ **Visual Controls**:
  - Depth slider (1-5 levels)
  - Delay selector (500ms-3000ms)
  - External links toggle
- 📊 **Real-time Stats**: Live progress with emoji feedback
  - Total pages found
  - Success/error counts
  - Processing time
  - Duplicate prevention stats
- 💾 **Smart Downloads**: One-click downloads with intelligent naming
  - `website-name_2024-01-15_14-30-25.md`
  - `website-name_2024-01-15_14-30-25.json`
- 📄 **Live Preview**: Expandable preview of scraped content
- 📱 **Responsive**: Works perfectly on mobile and desktop
- 🎯 **Built with**: Svelte 5 + Tailwind CSS 4 + TypeScript

## 📁 Output Formats & Smart Naming

All files now use **smart naming** with website + timestamp for better organization.

### Individual Files (`--format files`)

Creates separate markdown files for each page:

```
github-com_2024-01-15_14-30-25/
├── page-001-GitHub-Homepage.md
├── page-002-About-GitHub.md
├── page-003-Features.md
...
```

### Single Combined File (`--format single`)

Creates one comprehensive markdown file:

```
github-com_2024-01-15_14-30-25.md
```

Content structure:

```markdown
# Website Content: https://github.com

_Scraped on 2024-01-15 14:30:25_
_Found 15 unique pages (skipped 8 duplicates)_

---

## 📄 Page 1: GitHub Homepage

**URL:** https://github.com
**Depth:** 0

---

[Clean markdown content here]

## 📄 Page 2: About GitHub

**URL:** https://github.com/about
**Depth:** 1

---

[More content...]
```

### JSON Output (`--format json`)

Creates structured data file: `github-com_2024-01-15_14-30-25.json`

```json
[
	{
		"url": "https://github.com",
		"title": "GitHub: Let's build from here",
		"markdown": "[converted content]",
		"depth": 0,
		"error": ""
	},
	{
		"url": "https://github.com/about",
		"title": "About GitHub",
		"markdown": "[about content]",
		"depth": 1,
		"error": ""
	}
]
```

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

### POST `/api/v1/scrape`

Scrape a website and return the results.

**Request Body:**

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

### GET `/api/v1/status`

Get server status and available endpoints.

### GET `/health`

Health check endpoint.

## ⚙️ Configuration

### Environment Variables

You can configure the application using environment variables:

```bash
# API server port
export PORT=8080

# Default scraping settings
export DEFAULT_MAX_DEPTH=3
export DEFAULT_DELAY=1000
export DEFAULT_USER_AGENT="Website-Markdown-Converter/1.0"
```

### Respectful Scraping

The tool is designed to be respectful to target websites:

- ⏱️ **Default 1-second delays** between requests
- 🤖 **Proper User-Agent** identification
- 🚫 **Skips non-HTML content** automatically
- 📝 **Respects robots.txt** (planned feature)
- 🔒 **Rate limiting** in API mode

## 🛠️ Development

### 🔨 **Backend Development (Go)**

```bash
cd backend
go mod tidy                              # Install dependencies
go run main.go --help                   # Show CLI help
go run main.go --server --port 8080     # Start API server
go build -o website-markdown main.go    # Build binary
```

**Live reload during development:**

```bash
# Install air for live reloading (optional)
go install github.com/cosmtrek/air@latest
cd backend && air  # Auto-restart on file changes
```

### 🎨 **Frontend Development (Svelte + Bun)**

```bash
cd frontend
bun i       # Install dependencies
bun run dev         # Start dev server
```

### 📁 **Project Structure**

```
website-markdown/
├── 🗂️ backend/                    # Go API & CLI
│   ├── main.go                   # Entry point + server/CLI routing
│   ├── go.mod                    # Go dependencies
│   ├── website-markdown          # Built binary (after go build)
│   ├── cmd/
│   │   └── cli.go               # CLI commands & file operations
│   └── internal/
│       ├── scraper/
│       │   └── scraper.go       # Core scraping + duplicate prevention
│       └── api/
│           └── server.go        # HTTP API endpoints + CORS
├── 🎨 frontend/                   # Svelte UI
│   ├── src/
│   │   ├── routes/
│   │   │   └── +page.svelte     # Main UI component
│   │   └── app.html             # HTML template
│   ├── package.json             # Bun/Node dependencies
│   ├── tailwind.config.js       # Tailwind CSS config
│   └── vite.config.ts           # Vite build config
├── 📚 README.md                   # This documentation
├── 📄 LICENSE                     # MIT License
├── 🤝 CONTRIBUTING.md             # Contribution guidelines
└── 🏗️ .gitignore                 # Git ignore patterns
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

**Direct curl:**

```bash
# Quick scrape via API
curl -X POST http://localhost:8080/api/v1/scrape \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://github.com/features",
    "maxDepth": 2,
    "delay": 1000,
    "followExternal": false
  }'

# Health check
curl http://localhost:8080/health
# Returns: {"status":"healthy","time":"2024-01-15T14:30:25Z"}
```

**JavaScript fetch:**

```javascript
const response = await fetch('http://localhost:8080/api/v1/scrape', {
	method: 'POST',
	headers: { 'Content-Type': 'application/json' },
	body: JSON.stringify({
		url: 'https://docs.github.com',
		maxDepth: 3,
		delay: 1500,
		followExternal: false,
	}),
});

const result = await response.json();
console.log(
	`✅ Found ${result.stats.successPages} pages in ${result.stats.processingTime}`
);
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

## 🚀 Ready to Start Converting?

### 🏃‍♀️ **Quick Start (30 seconds)**

```bash
# Clone & build
git clone <repo-url> && cd website-markdown/backend
go build -o website-markdown main.go

# Test it immediately
./website-markdown https://example.com --depth 2
# 📄 Creates: example-com_2024-01-15_14-30-25.md
```

### 🌐 **Full Web Experience**

```bash
# Terminal 1: API Server
cd backend && ./website-markdown --server &

# Terminal 2: Svelte Frontend
cd frontend && bun run dev

# Open browser: http://localhost:5173 ✨
```

### 🎯 **Perfect For:**

- 📚 **Documentation archiving** - Offline access to docs
- 🔍 **SEO content audits** - Analyze site content structure
- 📖 **Research projects** - Convert sites to readable format
- 🤖 **AI/ML training data** - Clean markdown from web content
- 📱 **Content migration** - Export existing sites
- 🧑‍💻 **Developer workflows** - Integrate into build processes

### 💡 **Pro Tips:**

- Start with `depth 1-2` for large sites
- Use `--delay 2000+` for respectful scraping
- Try `--format json` for programmatic processing
- Use the web UI for interactive exploration
- Check console logs for duplicate prevention stats

**Ready to scrape smarter, not harder?** 🎯  
Visit `http://localhost:5173` and start converting! 🚀
