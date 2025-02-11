# Xendit Webhook Dumper

## Project Description
A lightweight Go-based microservice designed to receive Xendit webhooks and persistently store them as JSON files, organized by date.

## Features
- Receives Xendit webhook payloads via HTTP POST
- Saves each webhook payload as a unique JSON file
- Organizes files by date in `webhooks/data/YYYY-MM-DD/` directory
- Configurable port via environment variable
- Lightweight and easy to deploy

## Prerequisites
- Go 1.21+
- PM2 (Process Manager)
- Basic Linux/Unix environment

## Development Setup

### Local Development
1. Clone the repository
```bash
git clone <your-repo-url>
cd xendit-webhook-dumper
```

2. Initialize a new Go module
```bash
go mod init xendit-webhook-dumper
go mod tidy
```

3. Run the application locally
```bash
# Default port (8080)
go run main.go

# Custom port
PORT=3000 go run main.go
```

### Building for Production
```bash
# Check Go environment
go env GOOS GOARCH

# Build for server
GOOS=linux GOARCH=amd64 go build -o /var/www/xendit-webhook-dumper main.go
```

### Deploy with PM2
#### Initial Deployment
1. Transfer the binary to your server
2. Configure the ecosystem.config.js
3. Start the application
```bash
pm2 start ecosystem.config.js
```

#### Service Management
```bash
# Start Service
pm2 start ecosystem.config.js

# Stop Service
pm2 stop xendit-webhook-dumper

# Restart Service
pm2 restart xendit-webhook-dumper

# View Logs
pm2 logs xendit-webhook-dumper
```

#### Updating the Service
When code changes:
1. Rebuild the binary
2. Replace the existing binary
3. Restart the service
```bash
# Rebuild using the same command as `### Building for Production`
# Restart service
pm2 restart xendit-webhook-dumper
```

#### Webhook Storage
- Webhooks are stored in: webhooks/data/YYYY-MM-DD/
- Filename format: <uuid>_<webhook-id>.json

#### Security Considerations
- Implement authentication in production
- Use environment variables for sensitive configurations
- Regularly rotate and manage stored webhook data

#### Troubleshooting
- Check PM2 logs: pm2 logs xendit-webhook-dumper
- Verify binary permissions: chmod +x xendit-webhook-dumper
- Ensure correct path in ecosystem config

## License

MIT License

Copyright (c) 2025 Taufiq Ridwan Soleh

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.