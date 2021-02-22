# [Choice] Go version: 1, 1.16, 1.15
ARG GO_VARIANT=1
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${GO_VARIANT}

ARG INSTALL_NODE="true"
ARG NODE_VARIANT="lts/*"
RUN if [ "${INSTALL_NODE}" = "true" ]; then su vscode -c "umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VARIANT} 2>&1"; fi

# [Optional] Uncomment this section to install additional OS packages.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends nginx

# [Optional] Uncomment if you want to install an additional version of node using nvm
# ARG EXTRA_NODE_VERSION=10
# RUN su node -c "source /usr/local/share/nvm/nvm.sh && nvm install ${EXTRA_NODE_VERSION}"

# [Optional] Uncomment if you want to install more global node modules
ARG NODE_MODULES=""
RUN if [ "${NODE_MODULES}" != "" ]; then su vscode -c "npm install -g ${NODE_MODULES}"; fi
