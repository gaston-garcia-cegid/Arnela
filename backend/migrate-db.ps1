# Script para migrar base de datos existente a Docker
# Uso: .\migrate-db.ps1

Write-Host "=== Migración de Base de Datos a Docker ===" -ForegroundColor Cyan
Write-Host ""

# Configuración de la DB origen (tu PGAdmin actual)
$SourceHost = Read-Host "Host de la DB origen (ej: localhost)"
$SourcePort = Read-Host "Puerto de la DB origen (ej: 5432)"
$SourceDB = Read-Host "Nombre de la DB origen"
$SourceUser = Read-Host "Usuario de la DB origen"
$SourcePassword = Read-Host "Password de la DB origen" -AsSecureString
$SourcePasswordPlain = [Runtime.InteropServices.Marshal]::PtrToStringAuto(
    [Runtime.InteropServices.Marshal]::SecureStringToBSTR($SourcePassword)
)

# Configuración de la DB destino (Docker)
$DestHost = "localhost"
$DestPort = "5432"
$DestDB = "arnela_db"
$DestUser = "arnela_user"
$DestPassword = "arnela_password"

Write-Host ""
Write-Host "Configuración:" -ForegroundColor Yellow
Write-Host "  Origen: $SourceUser@$SourceHost`:$SourcePort/$SourceDB"
Write-Host "  Destino: $DestUser@$DestHost`:$DestPort/$DestDB"
Write-Host ""

$confirm = Read-Host "¿Continuar con la migración? (s/n)"
if ($confirm -ne "s") {
    Write-Host "Migración cancelada" -ForegroundColor Red
    exit
}

# Paso 1: Crear backup de la DB origen
Write-Host ""
Write-Host "[1/4] Creando backup de la DB origen..." -ForegroundColor Green
$backupFile = ".\backup_$(Get-Date -Format 'yyyyMMdd_HHmmss').sql"

$env:PGPASSWORD = $SourcePasswordPlain
pg_dump -h $SourceHost -p $SourcePort -U $SourceUser -d $SourceDB -F p -f $backupFile

if ($LASTEXITCODE -ne 0) {
    Write-Host "Error al crear backup" -ForegroundColor Red
    exit 1
}
Write-Host "✓ Backup creado: $backupFile" -ForegroundColor Green

# Paso 2: Verificar que Docker esté corriendo
Write-Host ""
Write-Host "[2/4] Verificando contenedor Docker..." -ForegroundColor Green
$container = docker ps --filter "name=arnela-postgres" --format "{{.Names}}"
if ($container -ne "arnela-postgres") {
    Write-Host "Contenedor arnela-postgres no está corriendo" -ForegroundColor Red
    Write-Host "Ejecuta: docker-compose up -d" -ForegroundColor Yellow
    exit 1
}
Write-Host "✓ Contenedor arnela-postgres activo" -ForegroundColor Green

# Paso 3: Limpiar DB destino (opcional)
Write-Host ""
$cleanDB = Read-Host "[3/4] ¿Limpiar DB destino antes de importar? (s/n)"
if ($cleanDB -eq "s") {
    Write-Host "Limpiando DB destino..." -ForegroundColor Yellow
    
    # Eliminar todas las tablas existentes
    $env:PGPASSWORD = $DestPassword
    $dropScript = @"
DO `$`$ DECLARE
    r RECORD;
BEGIN
    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
        EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
    END LOOP;
END `$`$;
"@
    
    $dropScript | docker exec -i arnela-postgres psql -U $DestUser -d $DestDB
    Write-Host "✓ DB destino limpiada" -ForegroundColor Green
}

# Paso 4: Restaurar backup en DB destino
Write-Host ""
Write-Host "[4/4] Restaurando backup en DB destino..." -ForegroundColor Green

# Copiar backup al contenedor
docker cp $backupFile arnela-postgres:/tmp/backup.sql

# Restaurar desde dentro del contenedor
$env:PGPASSWORD = $DestPassword
docker exec -i arnela-postgres psql -U $DestUser -d $DestDB -f /tmp/backup.sql

if ($LASTEXITCODE -ne 0) {
    Write-Host "Advertencia: Algunos errores durante la restauración (pueden ser normales)" -ForegroundColor Yellow
}

# Limpiar backup temporal del contenedor
docker exec arnela-postgres rm /tmp/backup.sql

Write-Host ""
Write-Host "✓ Migración completada!" -ForegroundColor Green
Write-Host ""
Write-Host "Pasos siguientes:" -ForegroundColor Cyan
Write-Host "  1. Verificar datos: docker exec -it arnela-postgres psql -U arnela_user -d arnela_db"
Write-Host "  2. Listar tablas: \dt"
Write-Host "  3. Si todo está OK, puedes eliminar el backup: $backupFile"
Write-Host ""
