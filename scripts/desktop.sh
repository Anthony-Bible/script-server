#!/usr/bin/env bash
#

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
        exit 1
        ;;
esac

# Remove /tmp/img.jpg if it exists
if [[ -f ~/Desktop/img.jpg ]]; then
    rm ~/Desktop/img.jpg
fi
curl -Lo ~/Desktop/img.jpg "https://anthony.bible/image?redirect=true" 

if [[ ! -f ~/Desktop/img.jpg ]]; then
    echo "Hmm that didn't download the file properly"
    exit 1
fi
osascript -e "tell application \"System Events\" to tell every desktop to set picture to \"$default_wallpaper\""
sleep 1
# Reset desktop picture
osascript -e 'tell application "System Events" to tell every desktop to set picture to "~/Desktop/img.jpg"'

sleep 1
osascript -e 'tell application "System Events" to set miniaturized of every window to true'

osascript -e 'tell app "System Events" to display dialog "An unknown error occurred. Have you tried turning it off and on again?" with title "System Error"'

osascript -e 'tell application "System Events"
    set allApps to every application process whose visible is true
    repeat with anApp in allApps
        set allWindows to every window of anApp
        repeat with aWindow in allWindows
            set miniaturized of aWindow to true
        end repeat
    end repeat
end tell'
