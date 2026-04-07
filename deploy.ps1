$SERVER = "192.168.1.100"
$USER = "root"
$REMOTE_DIR = "/DATA/AppData/Arnela"

Write-Host "=== Arnela Deploy Script ===" -ForegroundColor Cyan
Write-Host ""

# Files and directories to send
$items = @(
    "docker-compose.prod.yml",
    ".env.prod.example",
    "backend",
    "frontend",
    "nginx"
)

Write-Host "[1/4] Creating remote directory..." -ForegroundColor Yellow
ssh "${USER}@${SERVER}" "mkdir -p ${REMOTE_DIR}"

Write-Host "[2/4] Sending files to server..." -ForegroundColor Yellow
foreach ($item in $items) {
    $localPath = Join-Path $PSScriptRoot $item
    if (Test-Path $localPath) {
        Write-Host "  -> $item" -ForegroundColor Gray
        scp -r $localPath "${USER}@${SERVER}:${REMOTE_DIR}/"
    } else {
        Write-Host "  [SKIP] $item not found" -ForegroundColor Red
    }
}

Write-Host "[3/4] Creating .env.prod on server (if not exists)..." -ForegroundColor Yellow
ssh "${USER}@${SERVER}" @"
cd ${REMOTE_DIR}
if [ ! -f .env.prod ]; then
    cp .env.prod.example .env.prod
    # Generate secure passwords
    JWT=\$(openssl rand -base64 32)
    DB_PASS=\$(openssl rand -base64 16)
    REDIS_PASS=\$(openssl rand -base64 16)
    sed -i "s|CHANGE_ME_jwt_secret_here|\$JWT|" .env.prod
    sed -i "s|CHANGE_ME_strong_password_here|\$DB_PASS|" .env.prod
    sed -i "s|CHANGE_ME_redis_password_here|\$REDIS_PASS|" .env.prod
    echo '  -> .env.prod created with auto-generated secrets'
else
    echo '  -> .env.prod already exists, skipping'
fi
"@

Write-Host "[4/4] Building and starting containers..." -ForegroundColor Yellow
ssh "${USER}@${SERVER}" @"
cd ${REMOTE_DIR}
docker compose -f docker-compose.prod.yml --env-file .env.prod up -d --build
"@

Write-Host ""
Write-Host "=== Deploy complete! ===" -ForegroundColor Green
Write-Host "App:     http://${SERVER}" -ForegroundColor Cyan
Write-Host "Health:  http://${SERVER}/health" -ForegroundColor Cyan
Write-Host "Swagger: http://${SERVER}/swagger/index.html" -ForegroundColor Cyan
Write-Host "CasaOS:  http://${SERVER}:90" -ForegroundColor Cyan
