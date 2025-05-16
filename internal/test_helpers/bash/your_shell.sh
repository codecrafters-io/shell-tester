#!/bin/sh
export HISTFILE=/dev/null
BASH_SILENCE_DEPRECATION_WARNING=1 PS1='$ ' exec bash --norc -i
