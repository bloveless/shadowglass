package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"

	"shadowglass/internal/gen/tradersv1"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

type TradersEngine struct {
	runtime wazero.Runtime
	logic   tradersv1.Logic
}

func NewTradersEngine(ctx context.Context, r wazero.Runtime, l tradersv1.Logic) (*TradersEngine, error) {
	engine := TradersEngine{
		runtime: r,
		logic:   l,
	}

	err := tradersv1.AttachHostFunctions(ctx, engine.runtime, engine.logic)

	if err != nil {
		return new(TradersEngine), err
	}

	return &engine, nil
}

func (engine *TradersEngine) CompileModule(ctx context.Context, module []byte) (wazero.CompiledModule, string, error) {
	cm, err := engine.runtime.CompileModule(ctx, module)
	if err != nil {
		return nil, "", fmt.Errorf("failed to compile module: %w", err)
	}

	modSha := sha256.Sum256(goWasm)
	modName := fmt.Sprintf("rom-%x-%d", modSha, rand.Int())

	return cm, modName, nil
}

func (engine *TradersEngine) Invoke(ctx context.Context, mod wazero.CompiledModule, modName string) (retErr error) {
	modConf := wazero.NewModuleConfig().
		WithStdin(os.Stdin).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithName(modName)

	im, err := engine.runtime.InstantiateModule(ctx, mod, modConf)
	if err != nil {
		return fmt.Errorf("failed to instantiate module: %w", err)
	}
	defer func(im api.Module, ctx context.Context) {
		err := im.Close(ctx)
		if err != nil {
			retErr = fmt.Errorf("failed to close module: %w", err)
		}
	}(im, ctx)

	return err
}
