# https://forums.freebsd.org/threads/bracketed-paste.81314/#post-522692
typeset -A zle_bracketed_paste # Disable bracketed paste mode
unsetopt prompt_cr # Remove PROMPT_EOL_MARK + cursor movement