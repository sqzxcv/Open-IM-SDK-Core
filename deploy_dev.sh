#!/bin/bash

current_shell_path=$(dirname $0)
current_shell_path=${current_shell_path/\./$(pwd)}

cd $current_shell_path


make build-wasm
mv ./_output/bin/openIM.wasm ../futrtalk-dashboard/public/openIM.wasm
