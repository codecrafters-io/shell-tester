#!/bin/sh
sudo rm -rf /etc/zsh
sudo rm -rf /etc/zshenv
sudo rm -rf ~/.zshenv
sudo rm -rf /etc/zprofile
sudo rm -rf ~/.zprofile
sudo rm -rf /etc/zshrc
sudo rm -rf ~/.zshrc
sudo rm -rf /etc/zlogin
sudo rm -rf ~/.zlogin 
ZDOTDIR='/workspaces/shell-tester/internal/test_helpers/zsh_config' zsh
