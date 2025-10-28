package main

import (
	"context"
)

var (
	ctx, ctxCancel = context.WithCancel(context.Background())
)
