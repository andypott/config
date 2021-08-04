export EDITOR=nvim
export MOZ_ENABLE_WAYLAND=1

[[ ":$PATH:" != *"$HOME/config/bin"* ]] && PATH="$HOME/config/bin:${PATH}"
[[ ":$PATH:" != *"$HOME/bin"* ]] && PATH="$HOME/bin:${PATH}"
[[ ":$PATH:" != *"$HOME/go/bin"* ]] && PATH="$HOME/go/bin:${PATH}"
[[ ":$PATH:" != *"$HOME/.config/composer/vendor/bin"* ]] && PATH="$HOME/.config/composer/vendor/bin:${PATH}"
[[ ":$PATH:" != *"$HOME/.local/bin"* ]] && PATH="$HOME/.local/bin:${PATH}"
export PATH

export XKB_DEFAULT_LAYOUT='gb'


if [ -n "$BASH_VERSION" ]; then
    # include .bashrc if it exists
    if [ -f "$HOME/.bashrc" ]; then
	. "$HOME/.bashrc"
    fi
fi
