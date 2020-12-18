# .bashrc

# Source global definitions
# Fedora
[ -f /etc/bashrc ] && . /etc/bashrc
# Debian
[ -f /etc/bash_completion ] && . /etc/bash_completion
# Arch
[ -f /usr/share/bash-completion/bash_complation ] && . /usr/share/bash-completion/bash_complation

# User specific environment
[[ ":$PATH:" != *"$HOME/config/bin"* ]] && PATH="$HOME/config/bin:${PATH}"
[[ ":$PATH:" != *"$HOME/bin"* ]] && PATH="$HOME/bin:${PATH}"
[[ ":$PATH:" != *"$HOME/go/bin"* ]] && PATH="$HOME/go/bin:${PATH}"
[[ ":$PATH:" != *"$HOME/.config/composer/vendor/bin"* ]] && PATH="$HOME/.config/composer/vendor/bin:${PATH}"
[[ ":$PATH:" != *"$HOME/.local/bin"* ]] && PATH="$HOME/.local/bin:${PATH}"
export PATH

# Uncomment the following line if you don't like systemctl's auto-paging feature:
# export SYSTEMD_PAGER=

# User specific aliases and functions

# Add FZF keybindings
# Fedora
[ -f /usr/share/fzf/shell/key-bindings.bash ] && . /usr/share/fzf/shell/key-bindings.bash
# Debian
[ -f /usr/share/doc/fzf/examples/key-bindings.bash ] && . /usr/share/doc/fzf/examples/key-bindings.bash
# Arch
[ -f /usr/share/fzf/key-bindings.bash ] && . /usr/share/fzf/key-bindings.bash

# Use ripgrep for fzf
export FZF_DEFAULT_COMMAND='rg --files --hidden'
export FZF_CTRL_T_COMMAND="$FZF_DEFAULT_COMMAND"

export XKB_DEFAULT_LAYOUT='gb'

# Aliases
alias ls="ls --color"
alias ll="ls -lA"

# Disable command not found searching as it is slow and pointless
unset command_not_found_handle

