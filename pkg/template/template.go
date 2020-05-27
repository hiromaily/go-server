package template

import (
	"errors"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	fl "github.com/hiromaily/golibs/files"
	lg "github.com/hiromaily/golibs/log"
)

//type key int
//const templateKey = key(51)

var tempFiles *template.Template

// template FuncMap
func getTempFunc() template.FuncMap {
	//type FuncMap map[string]interface{}

	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"strAry": func(ary []string, i int) string {
			return ary[i]
		},
		"dateFmt": func(t time.Time) string {
			//fmt := "August 17, 2016 9:51 pm"
			//fmt := "2006-01-02 15:04:05"
			//fmt := "Mon Jan _2 15:04:05 2006"
			fmt := "Mon Jan _2 15:04:05"
			return t.Format(fmt)
		},
	}
	return funcMap
}

// LoadTemplatesFiles is to load template files for html response
func LoadTemplatesFiles() {
	ext := []string{"html"}
	//
	pwd, _ := os.Getwd()
	//fmt.Println("pwd", pwd)
	if strings.HasSuffix(pwd, "/cmd") {
		pwd += "/.."
	}
	files := fl.GetFileList(pwd+"/web/templates/", ext)
	if len(files) == 0{
		panic("template file can not be found")
	}
	//files2 := fl.GetFileList(pwd+"/submodules/global", ext)
	//files := append(append(files1, files2...), files3...)
	//files := append(files1, files2...)

	//*Template
	tempFiles = template.Must(template.New("").Funcs(getTempFunc()).ParseFiles(files...))
	//if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
	//	log.Println(err.Error())
	//	http.Error(w, http.StatusText(500), 500)
	//}
}

//func GetTemplate(ctx context.Context) (*template.Template, error) {
//	temp, ok := ctx.Value(templateKey).(*template.Template)
//	if !ok {
//		return nil, fmt.Errorf("%s", "couldn't find template in context")
//	}
//	return temp, nil
//}

//func SetTemplate(ctx context.Context, tmp *template.Template) context.Context {
//	ctx = context.WithValue(ctx, templateKey, tmp)
//	return ctx
//}

// Execute is to execute template for html response
func Execute(res http.ResponseWriter, key string, data interface{}) {
	if tempFiles != nil {
		err := tempFiles.ExecuteTemplate(res, key, data)
		if err != nil {
			lg.Error(err.Error())
			http.Error(res, http.StatusText(500), 500)
		}
	}

}
