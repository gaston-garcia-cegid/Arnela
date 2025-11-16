# ============================================
# Arnela Database Setup Script (PowerShell)
# ============================================

Write-Host "🚀 Starting Arnela Database Setup..." -ForegroundColor Cyan
Write-Host ""

# Configuration
$DB_USER = "arnela_user"
$DB_PASSWORD = "arnela_secure_pass_2024"
$DB_NAME = "arnela_db"
$POSTGRES_SUPERUSER = "postgres"

# Check if psql is available
try {
    $psqlVersion = psql --version
    Write-Host "✓ PostgreSQL client found: $psqlVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Error: psql command not found!" -ForegroundColor Red
    Write-Host "Please ensure PostgreSQL is installed and psql is in your PATH" -ForegroundColor Yellow
    Write-Host "You can add it manually, for example:" -ForegroundColor Yellow
    Write-Host '  $env:Path += ";C:\Program Files\PostgreSQL\16\bin"' -ForegroundColor Gray
    exit 1
}

Write-Host ""
Write-Host "📋 Configuration:" -ForegroundColor Cyan
Write-Host "  Database: $DB_NAME"
Write-Host "  User: $DB_USER"
Write-Host "  Password: $DB_PASSWORD"
Write-Host ""

# Prompt for postgres password
Write-Host "Enter password for PostgreSQL superuser ($POSTGRES_SUPERUSER):" -ForegroundColor Yellow
$PGPASSWORD = Read-Host -AsSecureString
$env:PGPASSWORD = [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR($PGPASSWORD))

Write-Host ""
Write-Host "🔧 Creating database and user..." -ForegroundColor Cyan

# Run the SQL script
try {
    psql -U $POSTGRES_SUPERUSER -f "setup_database.sql"
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "✅ Database setup completed successfully!" -ForegroundColor Green
        Write-Host ""
        Write-Host "📝 Connection details:" -ForegroundColor Cyan
        Write-Host "  Host: localhost"
        Write-Host "  Port: 5432"
        Write-Host "  Database: $DB_NAME"
        Write-Host "  User: $DB_USER"
        Write-Host "  Password: $DB_PASSWORD"
        Write-Host ""
        Write-Host "🚀 You can now start the backend with:" -ForegroundColor Green
        Write-Host "  go run cmd/api/main.go" -ForegroundColor Gray
        Write-Host ""
    } else {
        Write-Host ""
        Write-Host "❌ Error: Database setup failed!" -ForegroundColor Red
        Write-Host "Please check the error messages above." -ForegroundColor Yellow
        exit 1
    }
} catch {
    Write-Host ""
    Write-Host "❌ Error executing SQL script: $_" -ForegroundColor Red
    exit 1
} finally {
    # Clear password from environment
    Remove-Item Env:\PGPASSWORD -ErrorAction SilentlyContinue
}
