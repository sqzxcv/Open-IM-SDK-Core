#!/bin/bash

current_shell_path=$(dirname $0)
current_shell_path=${current_shell_path/\./$(pwd)}

cd $current_shell_path


make ios
mv ./build/FutrtalkIMCore.xcframework ../futrtalk-im-sdk-flutter/ios/frameworks/FutrtalkIMCore.xcframework
