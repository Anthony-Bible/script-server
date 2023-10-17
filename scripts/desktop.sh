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
        default_wallpaper=$(curl -s 'https://anthony.bible/image?redirect=true')
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
# Determine the number of desktops (monitors)
desktop_count=$(osascript -e 'tell application "System Events" to count of desktops')

# Loop through each desktop (monitor)
for i in $(seq 1 $desktop_count);
do
    # Set the image path based on monitor number
    image_path="~/Desktop/img${i}.jpg"

    # Remove the image if it exists
    if [[ -f $image_path ]]; then
        rm $image_path
    fi

    # Download a new image to the image path. Assume the server provides a different image based on the monitor number or some other criteria
    curl -Lo $image_path "https://anthony.bible/image?redirect=true"

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

osascript -e 'tell app "System Events" to display dialog "An unknown error occurred. Have you tried turning it off and on again?" with title "System Error"'

