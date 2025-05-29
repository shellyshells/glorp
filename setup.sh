#!/bin/bash

# GoForum Setup Script
echo "🚀 Setting up GoForum..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.19 or higher."
    exit 1
fi

echo "✅ Go is installed: $(go version)"

# Initialize Go module if not exists
if [ ! -f "go.mod" ]; then
    echo "📦 Initializing Go module..."
    go mod init goforum
fi

# Install dependencies
echo "📥 Installing dependencies..."
go get github.com/golang-jwt/jwt/v4
go get github.com/gorilla/mux
go get github.com/mattn/go-sqlite3
go mod tidy

# Create directories if they don't exist
echo "📁 Creating directory structure..."
mkdir -p static/css static/js static/images
mkdir -p views/layouts views/auth views/threads views/admin
mkdir -p database docs

# Build the application
echo "🔨 Building application..."
go build -o goforum

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo ""
    echo "🎉 GoForum is ready to run!"
    echo ""
    echo "To start the server:"
    echo "  ./goforum"
    echo "  or"
    echo "  go run main.go"
    echo ""
    echo "Then open your browser and go to: http://localhost:8080"
    echo ""
    echo "Default admin credentials:"
    echo "  Username: admin"
    echo "  Password: AdminPassword123!"
else
    echo "❌ Build failed. Please check the error messages above."
    exit 1
fi