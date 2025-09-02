package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/flopp/freiburg-fit/internal/app"
	"github.com/flopp/freiburg-fit/internal/utils"
)

type TemplateData struct {
	NowStr string
	Meta   struct {
		Title       string
		Description string
		Canonical   string
	}
	Venue    *app.Venue
	List     *app.Label
	isRemote bool
	basePath string
	CssFiles []string
	JSFiles  []string
	UmamiJS  string
	UmamiId  string
}

func (t TemplateData) IsRemoteTarget() bool {
	return t.isRemote
}

func (t TemplateData) BasePath() string {
	return t.basePath
}

func main() {
	// Parse command line options:
	// -config string
	var configFile string
	flag.StringVar(&configFile, "config", "", "Path to the config file")
	flag.Parse()

	// Load configuration
	config, err := app.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Collect YAML files of venues, extract metadata, build venue objects
	data, err := app.LoadData("data")
	if err != nil {
		log.Fatalf("Error loading data: %v", err)
	}

	// Set the current time
	now := time.Now()
	nowStr := now.Format("2006-01-02 15:04:05")

	// copy static files to output directory
	cssFiles, jsFiles, err := app.CopyAssets(config)
	if err != nil {
		log.Fatalf("Error copying assets: %v", err)
	}

	umamiJS := ""
	otherJS := make([]string, 0)
	// find umami.js file
	for _, jsFile := range jsFiles {
		if strings.Contains(jsFile, "umami") {
			umamiJS = jsFile
		} else {
			otherJS = append(otherJS, jsFile)
		}
	}

	canonical := func(p string) string {
		b := "https://freiburg.fit"
		if !strings.HasPrefix(p, "/") {
			b += "/"
		}
		b += p

		if !strings.Contains(p, ".") && !strings.HasSuffix(p, "/") {
			b += "/"
		}

		return b
	}

	// collect all canonical URLs for creating a sitemap
	sitemapUrls := make([]string, 0)

	// render templates
	pages := []struct {
		Title       string
		Description string
		Canonical   string
		Template    string
		OutFile     string
	}{
		{
			Title:       "freiburg.fit",
			Description: "Eine Übersicht über alle Fitness-Clubs, Studios und Locations in Freiburg.",
			Canonical:   "/",
			Template:    "index.html",
			OutFile:     "index.html",
		},
		{
			Title:       "Impressum - freiburg.fit",
			Description: "Impressum von freiburg.fit.",
			Canonical:   "/impressum.html",
			Template:    "impressum.html",
			OutFile:     "impressum.html",
		},
		{
			Title:       "Datenschutz - freiburg.fitlubs",
			Description: "Datenschutz von freiburg.fit.",
			Canonical:   "datenschutz.html",
			Template:    "datenschutz.html",
			OutFile:     "datenschutz.html",
		},
	}
	for _, page := range pages {
		t := TemplateData{
			NowStr: nowStr,
			Meta: struct {
				Title       string
				Description string
				Canonical   string
			}{
				Title:       page.Title,
				Description: page.Description,
				Canonical:   canonical(page.Canonical),
			},
			Venue:    nil,
			isRemote: config.IsRemoteTarget,
			basePath: config.OutputDir,
			CssFiles: cssFiles,
			JSFiles:  otherJS,
			UmamiJS:  umamiJS,
			UmamiId:  config.UmamiId,
		}
		if err := utils.ExecuteTemplate(page.Template, filepath.Join(config.OutputDir, page.OutFile), t); err != nil {
			log.Fatalf("rendering template %s: %w", page.Template, err)
		}

		sitemapUrls = append(sitemapUrls, t.Meta.Canonical)
	}

	// Render templates with collected data to output files
	for _, venue := range data.Venues {
		t := TemplateData{
			NowStr: nowStr,
			Meta: struct {
				Title       string
				Description string
				Canonical   string
			}{
				Title:     venue.Name,
				Canonical: canonical(venue.Slug()),
			},
			Venue:    venue,
			isRemote: config.IsRemoteTarget,
			basePath: config.OutputDir,
			CssFiles: cssFiles,
			JSFiles:  otherJS,
			UmamiJS:  umamiJS,
			UmamiId:  config.UmamiId,
		}

		target := fmt.Sprintf("%s/%s/index.html", config.OutputDir, venue.Slug())
		err := utils.ExecuteTemplate("venue.html", target, t)
		if err != nil {
			log.Fatalf("Error rendering template for venue %s: %v", venue.Name, err)
		}

		sitemapUrls = append(sitemapUrls, t.Meta.Canonical)
	}

	for _, category := range data.Categories {
		t := TemplateData{
			NowStr: nowStr,
			Meta: struct {
				Title       string
				Description string
				Canonical   string
			}{
				Title:     category.Name,
				Canonical: canonical(category.Slug()),
			},
			List:     category,
			isRemote: config.IsRemoteTarget,
			basePath: config.OutputDir,
			CssFiles: cssFiles,
			JSFiles:  otherJS,
			UmamiJS:  umamiJS,
			UmamiId:  config.UmamiId,
		}

		target := fmt.Sprintf("%s/%s/index.html", config.OutputDir, category.Slug())
		err := utils.ExecuteTemplate("list.html", target, t)
		if err != nil {
			log.Fatalf("Error rendering template for category %s: %v", category.Name, err)
		}

		sitemapUrls = append(sitemapUrls, t.Meta.Canonical)
	}

	/*
		for _, tag := range data.Tags {
			t := TemplateData{
				NowStr: nowStr,
				Meta: struct {
					Title       string
					Description string
					Canonical   string
				}{
					Title:     tag.Name,
					Canonical: canonical(tag.Slug()),
				},
				List:     tag,
				isRemote: config.IsRemoteTarget,
				basePath: config.OutputDir,
				CssFiles: cssFiles,
				JSFiles:  otherJS,
				UmamiJS:  umamiJS,
				UmamiId:  config.UmamiId,
			}

			target := fmt.Sprintf("%s/%s/index.html", config.OutputDir, tag.Slug())
			err := utils.ExecuteTemplate("list.html", target, t)
			if err != nil {
				log.Fatalf("Error rendering template for tag %s: %v", tag.Name, err)
			}
		}
	*/

	// create sitemap.xml
	sitemapFile := filepath.Join(config.OutputDir, "sitemap.xml")
	sitemapData := make([]byte, 0)
	sitemapData = append(sitemapData, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")...)
	sitemapData = append(sitemapData, []byte("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")...)
	for _, url := range sitemapUrls {
		sitemapData = append(sitemapData, []byte("  <url>\n")...)
		sitemapData = append(sitemapData, []byte(fmt.Sprintf("    <loc>%s</loc>\n", url))...)
		sitemapData = append(sitemapData, []byte("  </url>\n")...)
	}
	sitemapData = append(sitemapData, []byte("</urlset>\n")...)
	if err := os.WriteFile(sitemapFile, sitemapData, 0644); err != nil {
		log.Fatalf("writing sitemap file: %v", err)
	}
}
