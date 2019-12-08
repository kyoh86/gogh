function gogh() {
  exec 5>&1
  case $1 in
  "cd" )
    shift
    cd "$(command gogh find "$@" | tee /dev/tty | tail -n1)" || return
    ;;

  "get" )
    local CD=0
    for arg in "$@"; do
      if [ "${arg}" = '--cd' ]; then
        CD=1
      fi
    done

    if [ $CD -eq 1 ]; then
      loc="$(command gogh "$@" | tee /dev/tty | tail -n1)"
      cd "$loc" || return
    else
      command gogh "$@"
    fi
    ;;

  * )
    command gogh "$@"
    ;;
  esac
}
eval "$(command gogh --completion-script-zsh)"
