alias logbook=${LOGBOOK_PATH}/logbook/logbook
PROMPT_COMMAND='logbook add -command "$(history 1)"'
bind -m emacs-standard -x '"\C-r": xdotool type "$(logbook raw_query --query "SELECT * FROM command GROUP BY command_name ORDER BY id DESC" | fzf)"'
bind -m vi-command -x '"\C-r": xdotool type "$(logbook raw_query --query "SELECT * FROM command GROUP BY command_name ORDER BY id DESC" | fzf)"'
bind -m vi-insert -x '"\C-r": xdotool type "$(logbook raw_query --query "SELECT * FROM command GROUP BY command_name ORDER BY id DESC" | fzf)"'
