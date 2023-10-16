#!/usr/bin/env bash
#   
# Clone git@github.com:Anthony-Bible/dotfiles.git to ~/dotfiles and run setup.sh in that folder
#
# make sure git is installed
if ! command -v git  > /dev/null; then
  echo "git isn't installed"
  OSTYPE = $(uname -s)
  #install git for mac or linux
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        sudo apt-get install git
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        brew install git
    fi
fi

# clone dotfiles repo
# if dotfiles folder exists, cd into it and pull latest
# else clone repo
if [ -d ~/dotfiles ]; then
  cd ~/dotfiles
  git pull
else
  git clone
fi
# run setup.sh
#
~/dotfiles/setup.sh

