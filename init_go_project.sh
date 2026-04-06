#!/bin/bash

read -p "Enter project name (one word, lowercase): " name

if [[ ! "$name" =~ ^[a-z][a-z0-9_]*$ ]] || [[ "$name" =~ [[:space:]] ]]; then
    echo "Error: Project name must be one word, lowercase letters and numbers only."
    exit 1
fi

mkdir -p "$name"
cd "$name" || exit 1

go mod init "github.com/mdalasini/ccgo/$name"

cat > main.go << 'EOF'
package main

import "fmt"

func main() {
	fmt.Println("hello from $NAME")
}
EOF

sed -i "s/\$NAME/$name/" main.go

echo "Project '$name' initialized successfully."
