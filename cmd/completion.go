package cmd

const (
	bash_completion_func = `_go-pass_complete_entries () {
	local prefix="${PASSWORD_STORE_DIR:-$HOME/.password-store/}"
	prefix="${prefix%/}/"
	local suffix=".gpg"
	local autoexpand=${1:-0}

	local IFS=$'\n'
	local items=($(compgen -f $prefix$cur))
	local firstitem=""
	local i=0
        local item

	for item in ${items[@]}; do
		[[ $item =~ /\.[^/]*$ ]] && continue
		if [[ ${#items[@]} -eq 1 && $autoexpand -eq 1 ]]; then
			while [[ -d $item ]]; do
				local subitems=($(compgen -f "$item/"))
				local filtereditems=( ) item2
				for item2 in "${subitems[@]}"; do
					[[ $item2 =~ /\.[^/]*$ ]] && continue
					filtereditems+=( "$item2" )
				done
				if [[ ${#filtereditems[@]} -eq 1 ]]; then
					item="${filtereditems[0]}"
				else
					break
				fi
			done
		fi
		[[ -d $item ]] && item="$item/"

		item="${item%$suffix}"
		COMPREPLY+=("${item#$prefix}")
		if [[ $i -eq 0 ]]; then
			firstitem=$item
		fi
		let i+=1
	done
	if [[ $i -gt 1 || ( $i -eq 1 && -d $firstitem ) ]]; then
		compopt -o nospace
	fi
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
