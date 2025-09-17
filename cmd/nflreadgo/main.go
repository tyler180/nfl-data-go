package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tyler180/nfl-data-go/internal/download"
)

func main() {
	var (
		repo      string
		file      string
		season    int
		formatStr string
		outputDir string
		force     bool
		showHelp  bool
	)

	flag.StringVar(&repo, "repo", "", "Repository name (e.g., nflverse-data, nfldata)")
	flag.StringVar(&file, "file", "", "File name to download (e.g., plays.parquet)")
	flag.IntVar(&season, "season", 0, "Season year (optional, if applicable)")
	flag.StringVar(&formatStr, "format", "parquet", "File format (parquet or csv)")
	flag.StringVar(&outputDir, "out", ".", "Output directory")
	flag.BoolVar(&force, "force", false, "Force re-download even if file exists")
	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.Parse()

	if showHelp || repo == "" || file == "" {
		flag.Usage()
		os.Exit(1)
	}

	format, err := download.ParseFormat(formatStr)
	if err != nil {
		log.Fatalf("Error parsing format: %v", err)
	}

	dl := download.GetDownloader()
	b, usedURL, err := dl.Download(repo, file, &format, nil)
	if err != nil {
		log.Fatalf("Error downloading file: %v", err)
	}

	outputPath := fmt.Sprintf("%s/%s", outputDir, file)
	err = os.WriteFile(outputPath, b, 0644)
	if err != nil {
		log.Fatalf("Error writing file: %v", err)
	}

	fmt.Printf("File downloaded successfully to %s (source URL: %s)\n", outputPath, usedURL)

	// Alternative approach using DownloadFile:
	// if outputDir == "" {
	// 	outputDir = "."
	// }
	// if _, err := os.Stat(outputDir); os.IsNotExist(err) {
	// 	err = os.MkdirAll(outputDir, 0755)
	// 	if err != nil {
	// 		log.Fatalf("Error creating output directory: %v", err)
	// 	}
	// }

	// outputPath := filepath.Join(outputDir, file)
	// err = download.DownloadFile(repo, file, season, format, outputPath, force)
	// if err != nil {
	// 	log.Fatalf("Error downloading file: %v", err)
	// }

	// fmt.Printf("File downloaded successfully to %s\n", outputPath)
}

// 		Timeout:      30 * time.Second,
// 		CacheDir:     cacheDir,
// 		CacheEnabled: true,
// 	}
// }
// }
