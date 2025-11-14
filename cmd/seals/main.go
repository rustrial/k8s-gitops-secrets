package main

import (
	"context"

	"github.com/rustrial/k8s-gitops-secrets/cmd/seals/cmd"
)

func main() {
	cmd.Execute(context.Background())
}
