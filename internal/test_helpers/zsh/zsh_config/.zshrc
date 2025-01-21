# https://forums.freebsd.org/threads/bracketed-paste.81314/#post-522692

autoload -U compinit
compinit -u
# zstyle ':completion:*' completer _complete _match _approximate
zstyle ':completion:*' ignored-patterns 'echotc' 'echoti'

typeset -A zle_bracketed_paste # Disable bracketed paste mode
unsetopt prompt_cr # Remove PROMPT_EOL_MARK + cursor movement      