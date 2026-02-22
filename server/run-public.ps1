param(
    [string]$HTTP_PORT = $(if ($env:HTTP_PORT) { $env:HTTP_PORT } else { '8080' })
)

# Check for ngrok
if (-not (Get-Command ngrok -ErrorAction SilentlyContinue)) {
    Write-Error "ngrok not found in PATH. Install ngrok from https://ngrok.com and add it to PATH."
    exit 1
}

Write-Host "Starting ngrok on port $HTTP_PORT..."
Start-Process ngrok -ArgumentList "http", $HTTP_PORT -NoNewWindow

Write-Host "Starting server (from this folder) with HTTP_PORT=$HTTP_PORT..."
$env:HTTP_PORT = $HTTP_PORT
go run .
