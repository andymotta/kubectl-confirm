# kubectl-confirm

`kubectl-confirm` is a utility that helps you prevent accidental dangerous operations on critical Kubernetes contexts. It works by intercepting `kubectl` commands and checking the context and verb against a list of critical contexts and dangerous verbs specified in a configuration file.


## Quick Start

1. `cp sample_config.yaml ~/kube-confirm-config.yaml`
2. Update your `.zshrc` or a sourced file with the following Zsh function:
    ```zsh
    function kubectl() {
        if [ $# -lt 1 ]; then
            command kubectl
        else
            local context=$(command kubectl config current-context | cut -d/ -f2)
            kubectl-confirm "${context}" "$1"
            local exit_code=$?
            if [ $exit_code -eq 0 ]; then
                command kubectl "$@"
            fi
        fi
    }
    ```
3. *Optionally* set KUBECTL_CONFIRM_CONFIG in `.zshrc` to a custom path
    ```zsh
    export KUBECTL_CONFIRM_CONFIG=/path/to/your/config.yaml
    ```
