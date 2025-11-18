# Windows Installation Guide

This guide provides step-by-step instructions for setting up ReSell on Windows.

## Prerequisites

### Required Software

1. **Go 1.24+**
   - Download from: https://go.dev/dl/
   - Download the Windows installer (`.msi` file)
   - Run the installer and follow the prompts
   - Verify installation: Open PowerShell and run `go version`

2. **Docker Desktop for Windows**
   - Download from: https://www.docker.com/products/docker-desktop/
   - Install Docker Desktop
   - Enable WSL 2 backend (recommended) or Hyper-V
   - Start Docker Desktop and wait for it to be ready
   - Verify: Open PowerShell and run `docker --version` and `docker-compose --version`

3. **Git for Windows**
   - Download from: https://git-scm.com/download/win
   - Install with default options
   - Verify: Open PowerShell and run `git --version`

### Optional Tools

4. **Make for Windows** (recommended for easier commands)
   
   **Option A: Using Chocolatey**
   ```powershell
   # Install Chocolatey first (run PowerShell as Administrator)
   Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
   
   # Then install make
   choco install make
   ```
   
   **Option B: Using Scoop**
   ```powershell
   # Install Scoop first
   Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
   irm get.scoop.sh | iex
   
   # Then install make
   scoop install make
   ```
   
   **Option C: Manual Download**
   - Download from: http://gnuwin32.sourceforge.net/packages/make.htm
   - Add to PATH environment variable

5. **PostgreSQL Client Tools** (optional, for direct database access)
   - Download from: https://www.postgresql.org/download/windows/
   - Or use pgAdmin: https://www.pgadmin.org/download/

## Installation Steps

### Step 1: Clone the Repository

Open PowerShell or Git Bash:

```powershell
git clone https://github.com/Daumas-hugo/ReSell.git
cd ReSell
```

### Step 2: Copy Environment Configuration

```powershell
# PowerShell
Copy-Item .env.example .env

# Or Git Bash
cp .env.example .env
```

### Step 3: Install Go Dependencies

```powershell
go mod download
go mod tidy
```

### Step 4: Install Development Tools

#### Install golang-migrate

**Option A: Using Go**
```powershell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

**Option B: Download Binary**
1. Download from: https://github.com/golang-migrate/migrate/releases
2. Download the Windows `.exe` file
3. Rename to `migrate.exe`
4. Add to PATH or place in a directory that's already in PATH

#### Install SQLC

**Option A: Using Go**
```powershell
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

**Option B: Download Binary**
1. Download from: https://github.com/sqlc-dev/sqlc/releases
2. Download the Windows `.zip` file
3. Extract and add to PATH

Verify installations:
```powershell
migrate -version
sqlc version
```

### Step 5: Start Infrastructure Services

Ensure Docker Desktop is running, then:

```powershell
# Using make (if installed)
make docker-up

# Or using docker-compose directly
docker-compose up -d
```

Wait for services to be healthy (about 60 seconds):
```powershell
docker-compose ps
```

You should see:
- `resell_postgres` - healthy
- `resell_redis` - healthy
- `resell_keycloak` - healthy

### Step 6: Run Database Migrations

```powershell
# Using make (if installed)
make migrate-up

# Or using migrate directly
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/resell?sslmode=disable" up
```

### Step 7: Build the Application

```powershell
# Using make
make build

# Or using go directly
go build -o bin/api.exe cmd/api/main.go
```

### Step 8: Start the API Server

```powershell
# Using make
make run

# Or using go directly
go run cmd/api/main.go

# Or run the built binary
.\bin\api.exe
```

The API will be available at: http://localhost:8080

## Verify Installation

Test the API is running:

```powershell
# Using PowerShell
Invoke-RestMethod -Uri http://localhost:8080/health

# Or using curl (if installed)
curl http://localhost:8080/health
```

Expected response: `OK`

## Common Windows-Specific Issues

### Issue 1: Port Already in Use

If ports 5432, 6379, or 8180 are already in use:

1. Edit `docker-compose.yml`
2. Change the port mappings (left side only):
   ```yaml
   ports:
     - "5433:5432"  # Changed from 5432:5432
   ```
3. Update your `.env` file with the new port
4. Restart: `docker-compose down && docker-compose up -d`

### Issue 2: Docker Desktop Not Running

Error: `Cannot connect to the Docker daemon`

Solution: Start Docker Desktop from the Start menu and wait for it to be ready.

### Issue 3: WSL 2 Backend Issues

If Docker Desktop fails with WSL 2:

1. Open PowerShell as Administrator
2. Enable WSL 2:
   ```powershell
   wsl --install
   wsl --set-default-version 2
   ```
3. Restart your computer
4. Start Docker Desktop

### Issue 4: Line Endings (CRLF vs LF)

Git might change line endings on Windows. To prevent issues:

```powershell
# Configure git to not change line endings
git config --global core.autocrlf false

# Re-clone the repository
cd ..
Remove-Item -Recurse -Force ReSell
git clone https://github.com/Daumas-hugo/ReSell.git
cd ReSell
```

### Issue 5: Path Length Limit

Windows has a 260 character path limit. To enable long paths:

1. Open PowerShell as Administrator
2. Run:
   ```powershell
   New-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Control\FileSystem" -Name "LongPathsEnabled" -Value 1 -PropertyType DWORD -Force
   ```
3. Restart your computer

### Issue 6: Firewall Blocking Connections

If you can't connect to services:

1. Open Windows Defender Firewall
2. Click "Allow an app through firewall"
3. Add Docker Desktop and allow private/public networks
4. Or temporarily disable firewall for testing (not recommended for production)

## Alternative: Using Git Bash

If you prefer a Unix-like environment on Windows:

1. Install Git for Windows (includes Git Bash)
2. Open Git Bash instead of PowerShell
3. Follow the main Linux/macOS instructions from the README
4. Commands like `make`, `cp`, etc. will work as expected

## Database Access

### Using psql (if PostgreSQL client tools installed)

```powershell
psql -h localhost -p 5432 -U postgres -d resell
# Password: postgres
```

### Using pgAdmin

1. Open pgAdmin
2. Create new server:
   - Name: ReSell Local
   - Host: localhost
   - Port: 5432
   - Database: resell
   - Username: postgres
   - Password: postgres

## IDE Setup

### Visual Studio Code

Recommended extensions:
- Go (by Go Team at Google)
- Docker (by Microsoft)
- PostgreSQL (by Chris Kolkman)
- REST Client (by Huachao Mao) - for testing API endpoints

### GoLand (JetBrains)

1. Open the ReSell folder
2. GoLand will auto-detect the Go project
3. Configure Docker integration in Settings → Build, Execution, Deployment → Docker
4. Use built-in Database Tools for PostgreSQL access

## Running Tests

```powershell
# Run all tests
go test ./...

# Run tests with coverage
go test -v -coverprofile=coverage.txt -covermode=atomic ./...

# View coverage in browser
go tool cover -html=coverage.txt
```

## Development Workflow

### Making Changes

1. Edit code in your favorite IDE
2. Format code:
   ```powershell
   go fmt ./...
   ```
3. Build and test:
   ```powershell
   go build -o bin/api.exe cmd/api/main.go
   go test ./...
   ```
4. Run locally:
   ```powershell
   go run cmd/api/main.go
   ```

### Creating Migrations

```powershell
# Create new migration files
migrate create -ext sql -dir migrations -seq add_new_feature
```

This creates:
- `migrations\NNNNNN_add_new_feature.up.sql`
- `migrations\NNNNNN_add_new_feature.down.sql`

Edit the files, then apply:
```powershell
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/resell?sslmode=disable" up
```

## Stopping Services

```powershell
# Stop but keep data
docker-compose stop

# Stop and remove containers (keeps volumes)
docker-compose down

# Stop and remove everything including data
docker-compose down -v
```

## Cleanup

To completely remove ReSell:

```powershell
# Stop and remove containers, networks, volumes
cd ReSell
docker-compose down -v

# Remove the directory
cd ..
Remove-Item -Recurse -Force ReSell
```

## Additional Resources

- [Main README](../README.md)
- [Quick Start Guide](QUICKSTART.md) - Linux/macOS focused
- [API Documentation](API.md)
- [Architecture Guide](ARCHITECTURE.md)

## Getting Help

If you encounter issues:

1. Check Docker Desktop is running and healthy
2. Verify all prerequisites are installed
3. Check the [Common Issues](#common-windows-specific-issues) section
4. Review Docker logs: `docker-compose logs`
5. Check API logs in the console where you ran `go run`

## Next Steps

Once installed, check out:
- [API Documentation](API.md) - Learn about available endpoints
- [Quick Start Guide](QUICKSTART.md) - Example API calls
- [Architecture Guide](ARCHITECTURE.md) - Understand the system design
