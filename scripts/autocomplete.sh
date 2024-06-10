# Define the array of options
options=("init" "commit" "add")

# Function to handle autocomplete
_sasy_autocomplete() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    COMPREPLY=( $(compgen -W "${options[*]}" -- "$cur") )
}

# Register the autocomplete function for the 'sasy' command
complete -F _sasy_autocomplete sasy
