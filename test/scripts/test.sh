#!/bin/bash

print_parameters() {
    echo "COUNT PARAMS: $#"
    echo "PARAMS:"
    for param in "$@"; do
        echo "$param"
    done
}

echo "BASH_SH"
print_parameters "$@"