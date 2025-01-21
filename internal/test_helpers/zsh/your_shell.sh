#!/bin/sh
rm -rf /etc/zsh
rm -rf /etc/zshenv
rm -rf ~/.zshenv
rm -rf /etc/zprofile
rm -rf ~/.zprofile
rm -rf /etc/zshrc
rm -rf ~/.zshrc
rm -rf /etc/zlogin
rm -rf ~/.zlogin 
ZDOTDIR='/workspaces/shell-tester/internal/test_helpers/zsh/zsh_config' zsh
