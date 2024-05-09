FROM ghcr.io/scop/bash-completion/test:centos7 AS bash_completion
FROM mcr.microsoft.com/devcontainers/go:latest

RUN apt-get update && apt-get install -y \
    bash-completion \
    && rm -rf /var/lib/apt/lists/*

# Set a default value for ZSH_CUSTOM if it's not already set
ENV ZSH_CUSTOM=/home/vscode/.oh-my-zsh/custom

# Clone the powerlevel10k theme for zsh and change owner
RUN git clone --depth=1 https://github.com/romkatv/powerlevel10k.git $ZSH_CUSTOM/themes/powerlevel10k \
    && chown -R vscode:vscode $ZSH_CUSTOM/themes/powerlevel10k

# Clone the zsh-autosuggestions repository and change owner
RUN git clone https://github.com/zsh-users/zsh-autosuggestions $ZSH_CUSTOM/plugins/zsh-autosuggestions \
    && chown -R vscode:vscode $ZSH_CUSTOM/plugins/zsh-autosuggestions

COPY --chown=vscode:vscode configurations/.p10k.zsh /home/vscode/.p10k.zsh

COPY --chown=vscode:vscode configurations/.zshrc /home/vscode/.zshrc

COPY --chown=vscode:vscode configurations/p10k-instant-prompt-vscode.zsh /home/vscode/.cache/p10k-instant-prompt-vscode.zsh

# Set the default shell to Zsh
SHELL ["/bin/zsh", "-c"]

# Start a simple loop to keep the container running
CMD ["zsh", "-c", "while true; do sleep 3600; done"]