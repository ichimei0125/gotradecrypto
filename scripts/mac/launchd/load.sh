#! /bin/sh

go build -o gotradecypto cmd/main.go

sudo launchctl list | grep gotradecypto
sudo launchctl bootout system/localhost.gotradecypto.main
sudo launchctl list | grep gotradecypto
sudo chown root:wheel localhost.gotradecypto.main.plist 
sudo cp scripts/mac/launchd/localhost.gotradecypto.main.plist /Library/LaunchDaemons/
sudo launchctl bootstrap system /Library/LaunchDaemons/localhost.gotradecypto.main.plist
sudo launchctl list | grep gotradecypto
