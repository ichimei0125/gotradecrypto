#! /bin/sh

sudo launchctl list | grep gotradecypto
sudo launchctl bootout system/localhost.gotradecypto.main
sudo launchctl list | grep gotradecypto
sudo cp scripts/mac/launchd/localhost.gotradecypto.main.plist /Library/LaunchDaemons/
sudo launchctl bootstrap system /Library/LaunchDaemons/localhost.gotradecypto.main.plist
sudo launchctl list | grep gotradecypto
