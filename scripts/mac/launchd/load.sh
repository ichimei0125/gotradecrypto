#! /bin/sh

go build -o gotradecrypto cmd/main.go

sudo launchctl list | grep gotradecrypto
sudo launchctl bootout system/localhost.gotradecrypto.main
sudo launchctl list | grep gotradecrypto
sudo chown root:wheel localhost.gotradecrypto.main.plist 
sudo cp scripts/mac/launchd/localhost.gotradecrypto.main.plist /Library/LaunchDaemons/
sudo launchctl bootstrap system /Library/LaunchDaemons/localhost.gotradecrypto.main.plist
sudo launchctl list | grep gotradecrypto
