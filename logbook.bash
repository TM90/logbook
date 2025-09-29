function __logbook_hook__() {
  if [ "$XDG_SESSION_TYPE" = "wayland" ]; then
    wl-copy "$(logbook raw_query --query "SELECT id, command_name, history_id, exit_code, uuid, MAX(exec_time) FROM command GROUP BY command_name ORDER BY exec_time DESC" | fzf)"
  else
    xdotool type "$(logbook raw_query --query "SELECT id, command_name, history_id, exit_code, uuid, MAX(exec_time) FROM command GROUP BY command_name ORDER BY exec_time DESC" | fzf)"
  fi
}

alias logbook=${LOGBOOK_PATH}/logbook/logbook
export LOGBOOK_UUID=$(uuidgen -r)
PROMPT_COMMAND=('logbook add --exit-code=$? --command="$(history 1)" --uuid=$LOGBOOK_UUID' "${PROMPT_COMMAND[@]}")
bind -m emacs-standard -x '"\C-r": __logbook_hook__'
bind -m vi-command -x '"\C-r": __logbook_hook__'
bind -m vi-insert -x '"\C-r": __logbook_hook__'
