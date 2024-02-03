package main

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed examples/rust-rom/rust.wasm
var rustWasm []byte

//go:embed examples/basic-rom/main.wasm
var goWasm []byte

func main() {
	ctx := context.Background()

	ccache := wazero.NewCompilationCache()
	rcfg := wazero.NewRuntimeConfig().WithCompilationCache(ccache)
	r := wazero.NewRuntimeWithConfig(ctx, rcfg)
	defer func(r wazero.Runtime, ctx context.Context) {
		// This closes everything this Runtime created.
		err := r.Close(ctx)
		if err != nil {
			panic(err)
		}
	}(r, ctx)

	// Instantiate WASI, which implements host functions needed for TinyGo to
	// implement `panic`.
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	t := Traders{}
	engine, err := NewTradersEngine(ctx, r, t)
	if err != nil {
		panic(err)
	}

	timeoutContext, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
	defer cancel()

	goCm, goModName, err := engine.CompileModule(ctx, goWasm)
	if err != nil {
		panic(err)
	}

	rustCm, rustModName, err := engine.CompileModule(ctx, rustWasm)
	if err != nil {
		panic(err)
	}

	totalGoTime := int64(0)
	totalRustTime := int64(0)
	for i := 0; i < 100; i++ {
		goStart := time.Now()
		if err = engine.Invoke(timeoutContext, goCm, goModName); err != nil {
			panic(err)
		}
		goElapsed := time.Now().Sub(goStart)
		totalGoTime += goElapsed.Milliseconds()
		fmt.Printf("Go wasm execution %d took %dms\n", i, goElapsed.Milliseconds())

		rustStart := time.Now()
		if err = engine.Invoke(timeoutContext, rustCm, rustModName); err != nil {
			panic(err)
		}
		rustElapsed := time.Now().Sub(rustStart)
		totalRustTime += rustElapsed.Milliseconds()
		fmt.Printf("Rust wasm execution %d took %dms\n", i, rustElapsed.Milliseconds())
	}

	fmt.Printf("Total Go Time %dms; Avg Go Time %fms\n", totalGoTime, float64(totalGoTime)/100)
	fmt.Printf("Total Rust Time %dms; Avg Rust Time %fms\n", totalRustTime, float64(totalRustTime)/100)
}
