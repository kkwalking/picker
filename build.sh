# brew install mingw-w64

CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -ldflags="-H=windowsgui" -o rollcall.exe

# build in macos
/Users/zhouzekun/go/bin/fyne package -os darwin -icon icon.png --name picker
sudo spctl --master-disable
xattr -cr /Applications/picker.app