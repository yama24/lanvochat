#!/bin/bash
export PKG_CONFIG_PATH=/usr/lib/x86_64-linux-gnu/pkgconfig:/usr/share/pkgconfig
export CGO_ENABLED=1

# Auto-detect WebKit version
if pkg-config --exists webkit2gtk-4.1; then
    WEBKIT_VERSION="4.1"
    BUILD_TAGS="-tags webkit2gtk_4_1"
    echo "Detected: WebKit2GTK 4.1"
    
    # Create wrapper for systems with 4.1 only
    mkdir -p .build
    cat > .build/pkg-config << 'WRAPPER'
#!/bin/bash
args=("$@")
for i in "${!args[@]}"; do
    if [[ "${args[i]}" == "webkit2gtk-4.0" ]]; then
        args[i]="webkit2gtk-4.1"
    fi
done
/usr/bin/pkg-config "${args[@]}"
WRAPPER
    chmod +x .build/pkg-config
    export PATH="$PWD/.build:$PATH"
elif pkg-config --exists webkit2gtk-4.0; then
    WEBKIT_VERSION="4.0"
    BUILD_TAGS=""
    echo "Detected: WebKit2GTK 4.0"
else
    echo "ERROR: WebKit2GTK not found!"
    echo "Install: sudo apt install libwebkit2gtk-4.1-dev"
    exit 1
fi

export WEBKIT_VERSION

# Build production
$HOME/go/bin/wails build $BUILD_TAGS

echo ""
echo "Build selesai! Executable ada di: ./build/bin/lanvochat"
