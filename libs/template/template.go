package template

import (
	"context"
	"fmt"
	lg "github.com/hiromaily/golibs/log"
	"html/template"
	"net/http"
)

type key int

const templateKey = key(51)

func GetTemplate(ctx context.Context) (*template.Template, error) {
	temp, ok := ctx.Value(templateKey).(*template.Template)
	if !ok {
		return nil, fmt.Errorf("%s", "couldn't find template in context")
	}
	return temp, nil
}

func SetTemplate(ctx context.Context, tmp *template.Template) context.Context {
	ctx = context.WithValue(ctx, templateKey, tmp)
	return ctx
}

func Execute(res http.ResponseWriter, ctx context.Context, key string) {
	tmpl, err := GetTemplate(ctx)
	if err == nil {
		err = tmpl.ExecuteTemplate(res, key, nil)
	}

	if err != nil {
		lg.Error(err.Error())
		http.Error(res, http.StatusText(500), 500)
	}
}
