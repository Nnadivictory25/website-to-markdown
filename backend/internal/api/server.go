package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"website-markdown/internal/scraper"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	port   string
}

type ScrapeRequest struct {
	URL            string `json:"url" binding:"required"`
	MaxDepth       int    `json:"maxDepth"`
	Delay          int    `json:"delay"`
	FollowExternal bool   `json:"followExternal"`
}

type ScrapeResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Pages   []*scraper.ScrapedPage `json:"pages,omitempty"`
	Error   string                 `json:"error,omitempty"`
	Stats   *ScrapeStats           `json:"stats,omitempty"`
}

type ScrapeStats struct {
	TotalPages     int       `json:"totalPages"`
	SuccessPages   int       `json:"successPages"`
	ErrorPages     int       `json:"errorPages"`
	ProcessingTime string    `json:"processingTime"`
	StartedAt      time.Time `json:"startedAt"`
	CompletedAt    time.Time `json:"completedAt"`
}

func NewServer(port string) *Server {
	// Set gin mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:4173"} // Vite dev and preview ports
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	server := &Server{
		router: router,
		port:   port,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", s.healthCheck)

	// API routes
	api := s.router.Group("/api/v1")
	{
		api.POST("/scrape", s.scrapeWebsite)
		api.GET("/status", s.getStatus)
	}

	// Serve static files for docs (optional)
	s.router.Static("/docs", "./docs")
}

func (s *Server) Start() error {
	fmt.Printf("🚀 Starting API server on port %s\n", s.port)
	fmt.Printf("📝 API Endpoints:\n")
	fmt.Printf("   POST /api/v1/scrape - Scrape a website\n")
	fmt.Printf("   GET  /api/v1/status - Get server status\n")
	fmt.Printf("   GET  /health - Health check\n")
	fmt.Printf("🌐 CORS enabled for frontend on localhost:5173 and localhost:4173\n")

	return s.router.Run(":" + s.port)
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"time":    time.Now().UTC(),
		"service": "website-markdown-converter",
	})
}

func (s *Server) getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "running",
		"version": "1.0.0",
		"endpoints": []string{
			"POST /api/v1/scrape",
			"GET /api/v1/status",
			"GET /health",
		},
	})
}

func (s *Server) scrapeWebsite(c *gin.Context) {
	startTime := time.Now()

	var req ScrapeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ScrapeResponse{
			Success: false,
			Error:   fmt.Sprintf("❌ Invalid request: %v", err),
		})
		return
	}

	// Validate URL
	if req.URL == "" {
		c.JSON(http.StatusBadRequest, ScrapeResponse{
			Success: false,
			Error:   "❌ URL is required",
		})
		return
	}

	// Set defaults
	if req.MaxDepth <= 0 {
		req.MaxDepth = 3
	}
	if req.MaxDepth > 10 {
		req.MaxDepth = 10 // Prevent abuse
	}
	if req.Delay <= 0 {
		req.Delay = 1000
	}
	if req.Delay < 500 {
		req.Delay = 500 // Minimum delay to be respectful
	}

	fmt.Printf("🔄 API scrape request: %s (depth: %d, delay: %dms, external: %t)\n",
		req.URL, req.MaxDepth, req.Delay, req.FollowExternal)

	// Create scraper config
	config := &scraper.ScrapingConfig{
		MaxDepth:       req.MaxDepth,
		Delay:          time.Duration(req.Delay) * time.Millisecond,
		FollowExternal: req.FollowExternal,
		UserAgent:      "Website-Markdown-API/1.0",
	}

	// Perform scraping
	scrapeInstance := scraper.NewScraper(config)
	pages, err := scrapeInstance.ScrapeWebsite(req.URL)

	endTime := time.Now()
	processingTime := endTime.Sub(startTime)

	if err != nil {
		fmt.Printf("❌ Scraping failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, ScrapeResponse{
			Success: false,
			Error:   fmt.Sprintf("Scraping failed: %v", err),
			Stats: &ScrapeStats{
				ProcessingTime: processingTime.String(),
				StartedAt:      startTime,
				CompletedAt:    endTime,
			},
		})
		return
	}

	// Calculate stats
	stats := calculateStats(pages, startTime, endTime)

	fmt.Printf("✅ Scraping completed: %d pages (%d successful, %d errors) in %v\n",
		stats.TotalPages, stats.SuccessPages, stats.ErrorPages, processingTime)

	c.JSON(http.StatusOK, ScrapeResponse{
		Success: true,
		Message: fmt.Sprintf("🎉 Successfully scraped %d pages", stats.SuccessPages),
		Pages:   pages,
		Stats:   stats,
	})
}

func calculateStats(pages []*scraper.ScrapedPage, startTime, endTime time.Time) *ScrapeStats {
	stats := &ScrapeStats{
		TotalPages:     len(pages),
		ProcessingTime: endTime.Sub(startTime).String(),
		StartedAt:      startTime,
		CompletedAt:    endTime,
	}

	for _, page := range pages {
		if page.Error != "" {
			stats.ErrorPages++
		} else {
			stats.SuccessPages++
		}
	}

	return stats
}

// Helper function to start server from main
func StartAPIServer(port string) error {
	if port == "" {
		port = "8080"
	}

	// Validate port
	if _, err := strconv.Atoi(port); err != nil {
		return fmt.Errorf("❌ Invalid port: %s", port)
	}

	server := NewServer(port)
	return server.Start()
}
