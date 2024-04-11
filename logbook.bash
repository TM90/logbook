
alias logbook=${LOGBOOK_PATH}/logbook/logbook
export LOGBOOK_UUID=$(uuidgen -r)
PROMPT_COMMAND='logbook add --exit-code=$? --command="$(history 1)" --uuid=$LOGBOOK_UUID'
bind -m emacs-standard -x '"\C-r": xdotool type "$(logbook raw_query --query "SELECT * FROM command GROUP BY command_name ORDER BY id DESC" | fzf)"'
bind -m vi-command -x '"\C-r": xdotool type "$(logbook raw_query --query "SELECT * FROM command GROUP BY command_name ORDER BY id DESC" | fzf)"'
bind -m vi-insert -x '"\C-r": xdotool type "$(logbook raw_query --query "SELECT * FROM command GROUP BY command_name ORDER BY id DESC" | fzf)"'
