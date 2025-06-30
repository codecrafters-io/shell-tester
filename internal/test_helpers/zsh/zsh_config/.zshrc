# https://forums.freebsd.org/threads/bracketed-paste.81314/#post-522692

# echo "Loading zsh config"
autoload -U compinit
compinit -u
# zstyle ':completion:*' completer _complete _match _approximate
zstyle ':completion:*' ignored-patterns 'echotc' 'echoti'

# Set prompt to match bash for consistent testing
PS1='$ '

typeset -A zle_bracketed_paste # Disable bracketed paste mode
unsetopt prompt_cr # Remove PROMPT_EOL_MARK + cursor movement
# setopt autolist
# unsetopt menucomplete
# setopt noautomenu

# zstyle ':completion:*' list-prompt %S
# zstyle ':completion:*' list-lines 1000
# zstyle ':completion:*' menu no
# zstyle ':completion:*' select 0

# _list-or-complete-newline() {
#   zle expand-or-complete
#   if [[ $? -eq 0 ]]; then
#     return
#   fi

#   zle list-choices
#     if [[ $? -eq 0 ]]; then
#         zle push-line
#         zle redisplay
#     fi
# }
# zle -N _list-or-complete-newline # Make it available to zle


# Bind that widget to tab
# bindkey "^I" _list-or-complete-newline
# zstyle ':completion:*' menu no
# zstyle ':completion:*' select 0