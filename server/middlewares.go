package server

import (
	"context"
	"log"

	"github.com/fulldump/box"

	persistence "spikepersistence"
)

func ContextPersistence(ctx context.Context) persistence.Persistencer {
	p, ok := ctx.Value("persistence").(persistence.Persistencer)
	if !ok {
		panic("Persistence is not in context, developer fault!")
	}

	return p
}

func MiddlewareSetPersistence(p persistence.Persistencer) box.I {
	return func(next box.H) box.H {
		return func(ctx context.Context) {
			ctx = context.WithValue(ctx, "persistence", p)
			next(ctx)
		}
	}
}

func MiddlewareAccessLog() box.I {
	return func(next box.H) box.H {
		return func(ctx context.Context) {
			r := box.GetRequest(ctx)

			action := box.GetBoxContext(ctx).Action
			actionName := ""
			if action != nil {
				actionName = action.Name
			}

			log.Println(r.Method, r.URL.String(), actionName)
			next(ctx)
		}
	}
}
