#!/bin/bash

current_shell_path=$(dirname $0)
current_shell_path=${current_shell_path/\./$(pwd)}

cd $current_shell_path


make build-wasm-debug
mv ./_output/bin/openIM.wasm ../futrtalk-dashboard/public/futrtalk.wasm
