package cmd

const (
	bash_completion_func = `_go-pass_complete_entries () {
    local IFS=$'\n'
    COMPREPLY=( $(compgen -W "$(go-pass ls | tail +2 )" -- "${COMP_WORDS[${#COMP_WORDS[@]}-1]}") )
}

__custom_func() {
    case ${last_command} in
        go-pass_show | go-pass_cp | go-pass_edit | go-pass_generate | go-pass_ls | go-pass_mv | go-pass_rm | go-pass | go-pass_otp)
            _go-pass_complete_entries
            return
            ;;
        *)
            ;;
    esac
}
`
	zsh_completion_func = `
_go-pass_complete_entries () {
	local IFS=$'\n'
	_values -C 'passwords' $(go-pass ls | tail +2)
}
`
)
