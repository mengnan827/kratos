package log

import "context"

type Valuer func(ctx context.Context) any

func bindValues(ctx context.Context, keyVals []any) {
	for i := 1; i < len(keyVals); i += 2 {
		if v, ok := keyVals[i].(Valuer); ok {
			keyVals[i] = v(ctx)
		}
	}
}

func containsValuer(keyVals []any) bool {
	for i := 1; i < len(keyVals); i += 2 {
		if _, ok := keyVals[i].(Valuer); ok {
			return true
		}
	}
	return false
}
