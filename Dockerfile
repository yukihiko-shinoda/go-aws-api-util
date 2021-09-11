FROM golang:1.16.5-buster as base
# For compatibility with Visual Studio Code
WORKDIR /workspace
# 2021-06-10 Visual Studio Code uses
# Ctrl + Shift + P -> Go: Install/Update Tools
RUN go get github.com/uudashr/gopkgs/v2/cmd/gopkgs
RUN go get github.com/ramya-rao-a/go-outline
RUN go get github.com/cweill/gotests/...
RUN go get github.com/fatih/gomodifytags
RUN go get github.com/josharian/impl
RUN go get github.com/haya14busa/goplay/cmd/goplay
RUN go get github.com/go-delve/delve/cmd/dlv
# see: https://github.com/golang/vscode-go/blob/master/docs/dlv-dap.md
RUN GOBIN=/tmp/ go get github.com/go-delve/delve/cmd/dlv@master \
 && mv /tmp/dlv $GOPATH/bin/dlv-dap
RUN go get github.com/golangci/golangci-lint/cmd/golangci-lint
RUN go get golang.org/x/tools/gopls
# Install dependencies
# COPY . /workspace/
# RUN go install
