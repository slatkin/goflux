package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/slatkin/goflux/internal/ui"
	"github.com/slatkin/goflux/pkg/config"
)

func main() {
	initFlag := flag.Bool("init", false, "Initialize default configuration file")
	flag.Parse()

	if *initFlag {
		path, err := config.Init()
		if err != nil {
			fmt.Printf("Error initializing config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Wrote default configuration file to %s\n", path)
		os.Exit(0)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(ui.NewModel(cfg))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
