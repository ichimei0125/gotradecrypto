#! /bin/sh

sudo launchctl list | grep gotradecrypto
sudo launchctl bootout system/localhost.gotradecrypto.main
sudo rm /Library/LaunchDaemons/localhost.gotradecrypto.main.plist
sudo launchctl list | grep gotradecrypto