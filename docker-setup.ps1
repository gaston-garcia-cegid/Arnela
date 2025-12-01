# Docker Setup and Management Script for Arnela
# Usage: .\docker-setup.ps1 [command]
# Commands: build, up, down, restart, logs, clean

param(
    [Parameter(Mandatory=$false)]
    [ValidateSet('build', 'up', 'down', 'restart', 'logs', 'clean', 'status')]
    [string]$Command = 'up'
)

$ErrorActionPreference = "Stop"

function Write-ColorOutput($ForegroundColor) {
    $fc = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    if ($args) {
        Write-Output $args
    }
    $host.UI.RawUI.ForegroundColor = $fc
}

Write-ColorOutput Green "=== Arnela Docker Manager ==="
Write-Output ""

switch ($Command) {
    'build' {
        Write-ColorOutput Cyan "Building Docker images..."
        docker-compose build
        Write-ColorOutput Green "✓ Build complete"
    }
    
    'up' {
        Write-ColorOutput Cyan "Starting Docker services..."
        docker-compose up -d
        Start-Sleep -Seconds 3
        Write-Output ""
        Write-ColorOutput Green "✓ Services started"
        Write-Output ""
        Write-ColorOutput Yellow "Service Status:"
        docker-compose ps
        Write-Output ""
        Write-ColorOutput Cyan "Health Check:"
        Start-Sleep -Seconds 2
        try {
            $response = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get
            $response | ConvertTo-Json
        } catch {
            Write-ColorOutput Red "Health check failed - services may still be starting..."
        }
    }
    
    'down' {
        Write-ColorOutput Cyan "Stopping Docker services..."
        docker-compose down
        Write-ColorOutput Green "✓ Services stopped"
    }
    
    'restart' {
        Write-ColorOutput Cyan "Restarting Docker services..."
        docker-compose restart
        Start-Sleep -Seconds 3
        Write-ColorOutput Green "✓ Services restarted"
        docker-compose ps
    }
    
    'logs' {
        Write-ColorOutput Cyan "Showing logs (Ctrl+C to exit)..."
        docker-compose logs -f
    }
    
    'clean' {
        Write-ColorOutput Yellow "Warning: This will remove all containers, volumes, and images"
        $confirmation = Read-Host "Are you sure? (yes/no)"
        if ($confirmation -eq 'yes') {
            Write-ColorOutput Cyan "Cleaning up..."
            docker-compose down -v
            docker system prune -f
            Write-ColorOutput Green "✓ Cleanup complete"
        } else {
            Write-ColorOutput Yellow "Cleanup cancelled"
        }
    }
    
    'status' {
        Write-ColorOutput Cyan "Service Status:"
        docker-compose ps
        Write-Output ""
        Write-ColorOutput Cyan "Health Check:"
        try {
            $response = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get
            $response | ConvertTo-Json -Depth 10
        } catch {
            Write-ColorOutput Red "API not reachable"
        }
    }
}

Write-Output ""
Write-ColorOutput Green "=== Done ==="
