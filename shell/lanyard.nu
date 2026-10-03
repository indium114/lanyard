# lanyard integration for nushell.
# !! DO NOT EDIT !!

# MARK: instructions
# 1. save this to somewhere like '~/.config/nushell/lanyard.nu'
# 2. add 'source ~/.config/nushell/lanyard.nu' to your config.nu file

lanyard;
bash -c "cat $HOME/.lanyard/activate.nu" | from nuon | load-env;
