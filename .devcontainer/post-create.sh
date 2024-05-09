#!/bin/sh

# Link WSL bash_history to bash_history
ln -sf /WSL_USER/.zsh_history ~/.zsh_history

# Install commitizen globally
npm install -g commitizen

# Run go mod tidy
go mod tidy

# Add configuration to .czrc for commitizen
echo '{"path": "cz-conventional-changelog"}' >> ~/.czrc
