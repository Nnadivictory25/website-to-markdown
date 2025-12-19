# Contributing to Website to Markdown Converter 🤝

Thank you for your interest in contributing! This project welcomes contributions from everyone.

## 🚀 Quick Start

1. **Fork** the repository
2. **Clone** your fork locally
3. **Create** a feature branch
4. **Make** your changes
5. **Test** your changes
6. **Submit** a pull request

## 💻 Development Setup

### Backend (Go)

```bash
cd backend
go mod tidy
go run main.go --help  # Test CLI
go run main.go --server  # Test API
```

### Frontend (Svelte + Bun)

```bash
cd frontend
bun install
bun run dev  # Start dev server
```

## 🧪 Testing

### CLI Testing

```bash
cd backend
go build -o website-markdown main.go
./website-markdown https://example.com --depth 1
```

### API Testing

```bash
# Terminal 1: Start server
./website-markdown --server

# Terminal 2: Test endpoint
curl -X POST http://localhost:8080/api/v1/scrape \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com","maxDepth":1}'
```

## 📝 Code Style

### Go Code

- Use `gofmt` to format code
- Follow standard Go conventions
- Add emoji to console outputs for consistency
- Keep functions focused and well-named

### Frontend Code

- Use Prettier for formatting
- Follow Svelte 5 best practices
- Use Tailwind CSS for styling
- Maintain TypeScript strict mode

## 🎯 Areas for Contribution

### 🐛 Bug Fixes

- Fix duplicate detection edge cases
- Improve error handling
- Handle special characters in URLs

### ✨ Features

- Add robots.txt respect
- Implement sitemap parsing
- Add more output formats (PDF, EPUB)
- Improve markdown cleaning
- Add content filtering options

### 🎨 UI/UX

- Improve mobile responsiveness
- Add dark mode toggle
- Better progress indicators
- Keyboard shortcuts

### 📚 Documentation

- Add more examples
- Improve API documentation
- Create video tutorials
- Translate to other languages

## 🔍 Pull Request Guidelines

### Before Submitting

- [ ] Code compiles without warnings
- [ ] New features include basic testing
- [ ] Documentation updated if needed
- [ ] Commit messages are clear

### PR Template

```markdown
## Description

Brief description of changes

## Type of Change

- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Performance improvement

## Testing

- [ ] Tested CLI functionality
- [ ] Tested web interface
- [ ] Manual testing completed

## Screenshots (if applicable)

[Add screenshots for UI changes]
```

## 🚨 Issue Reporting

When reporting bugs, please include:

- **URL** you were trying to scrape
- **Command** or settings used
- **Expected** vs **actual** behavior
- **Error messages** or console output
- **Environment** (OS, Go version, etc.)

## 💡 Feature Requests

For new features:

- Explain the **use case**
- Describe the **expected behavior**
- Consider **alternative solutions**
- Think about **implementation complexity**

## 🎉 Recognition

Contributors will be:

- Added to the README acknowledgments
- Mentioned in release notes
- Given credit for their contributions

## ❓ Questions?

- Create an issue for general questions
- Join discussions in existing issues
- Check the README for common solutions

---

**Happy contributing!** 🚀 Every contribution makes this tool better for everyone.
