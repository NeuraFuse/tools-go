#!/bin/bash

printf "\nContainer | NeuraFuse\n\n"
printf "Preparing environment..\n"

printf "Fetching dependencies..\n\n"
export PATH=$PATH:/usr/local/go/bin

go get -d ./...

export GO111MODULE=on
go get sigs.k8s.io/kustomize/pkg/ifc > /dev/null 2>&1
go get sigs.k8s.io/kustomize/pkg/types > /dev/null 2>&1
go get sigs.k8s.io/kustomize/pkg/gvk > /dev/null 2>&1
go get sigs.k8s.io/kustomize/pkg/resmap > /dev/null 2>&1
go get sigs.k8s.io/kustomize/pkg/transformers > /dev/null 2>&1
go get sigs.k8s.io/kustomize/pkg/resource > /dev/null 2>&1
go get sigs.k8s.io/kustomize/pkg/factory > /dev/null 2>&1
go get sigs.k8s.io/kustomize/pkg/commands/build > /dev/null 2>&1
go get sigs.k8s.io/kustomize/pkg/fs > /dev/null 2>&1

# kube portfwd
go get sigs.k8s.io/kustomize/pkg/commands/build

export GO111MODULE=off

printf "\n\nContainer prepared.\n\n"