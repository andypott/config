# .bashrc

# Source global definitions
# Fedora
if [ -f /etc/bashrc ]; then
	. /etc/bashrc
fi
# Debian
if [ -f /etc/bash_completion ]; then
	. /etc/bash_completion
fi

# User specific environment
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
if [ -f /usr/share/fzf/shell/key-bindings.bash ]; then
	. /usr/share/fzf/shell/key-bindings.bash
fi
# Debian
if [ -f /usr/share/doc/fzf/examples/key-bindings.bash ]; then
	. /usr/share/doc/fzf/examples/key-bindings.bash
fi
# Arch
if [ -f /usr/share/fzf/key-bindings.bash ]; then
	. /usr/share/fzf/key-bindings.bash
fi

# Use ripgrep for fzf
export FZF_DEFAULT_COMMAND='rg --files --hidden'
export FZF_CTRL_T_COMMAND="$FZF_DEFAULT_COMMAND"

export XKB_DEFAULT_LAYOUT='gb'

# Aliases
alias ls="ls --color"
alias ll="ls -lA"

# Disable command not found searching as it is slow and pointless
unset command_not_found_handle

