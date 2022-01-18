package api

import (
	"context"
	"fmt"
	"testing"
)

func TestNewApiName(t *testing.T) {
	// os.Setenv("PORT", "8080")
	t.Setenv(LISTEN_ADDR, ":8080")
	Create(context.Background(), false)
	fmt.Printf("%v/n", Name())

}
