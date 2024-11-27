#! /bin/sh

sudo launchctl list | grep gotradecypto
sudo launchctl bootout system/localhost.gotradecypto.main
sudo rm /Library/LaunchDaemons/localhost.gotradecypto.main.plist
sudo launchctl list | grep gotradecypto