package web

import (
	"errors"
	"strconv"

	wren "github.com/crazyinfin8/WrenGo"
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func CreateForeignClasses(vm *wren.VM, app *App) {
	vm.SetModule("web", wren.NewModule(wren.ClassMap{
		"Routes": wren.NewClass(nil, nil, wren.MethodMap{
			"static GET(_,_)": func(vm *wren.VM, parameters []interface{}) (interface{}, error) {
				log.Debugf("Adding route %s", parameters[1])
				str, ok := parameters[1].(string)
				if !ok {
					log.Fatal("Must pass a string to the first argument of Routes.GET")
				}
				if app.HasRoute(str) {
					log.Errorf("Route %s already registered, ignoring", str)

					return nil, nil
				} else {
					app.Routes = append(app.Routes, str)
				}
				app.Router.GET(str, func(context *gin.Context) {
					handle, ok := parameters[2].(*wren.Handle)
					if !ok {
						log.Fatal("Must pass a handle to the second argument of Routes.GET")
					}

					callHandle, err := handle.Func("call(_)")
					if err != nil {
						log.Fatal("Must pass a handle with 0-1 parameters to the second argument of Routes.GET")
					}
					params, err := vm.NewMap()
					if err != nil {
						log.Fatalf("An error occurred when creating a map: %s", err.Error())
						return // IDE seems to want this
					}
					defer params.Free()
					for _, param := range context.Params {
						params.Set(param.Key, param.Value)
					}
					result, err := callHandle.Call(params)
					if err != nil {
						context.Header("Content-Type", "text/html")
						context.String(500, "An error occurred: %s", err.Error())
						return
					}

					out, ok := result.(string)

					if !ok {
						log.Fatal("Must return a string")
					}

					context.Header("Content-Type", "text/html")
					context.String(200, out)

				})
				return nil, nil
			},
		}),
		"App": wren.NewClass(nil, nil, wren.MethodMap{
			"static run(_)": func(vm *wren.VM, parameters []interface{}) (interface{}, error) {
				if app.IsServing {
					return nil, nil
				} else {
					app.IsServing = true
				}
				portFloat, ok := parameters[1].(float64)
				if !ok {
					log.Fatalf("Invalid port number")
				}
				port := int(portFloat)
				go app.Router.Run("0.0.0.0:" + strconv.Itoa(port))
				return nil, nil
			},
		}),
		"TemplatesHelper": wren.NewClass(nil, nil, wren.MethodMap{
			"static render(_,_,_)": func(vm *wren.VM, parameters []interface{}) (interface{}, error) {
				engine := "jinja"
				path, ok := parameters[1].(string)
				if !ok {
					return nil, errors.New("must pass a string to the first argument of Templates.render")
				}
				wMap, ok := parameters[2].(*wren.MapHandle)
				if !ok {
					return nil, errors.New("must pass a map to the first argument of Templates.render")
				}
				keys, ok := parameters[3].(*wren.ListHandle)
				if !ok {
					return nil, errors.New("must pass a list to the second argument of Templates.render")
				}
				goMap, err := wrenMapToGoMap(wMap, keys)
				if err != nil {
					return nil, err
				}
				if len(parameters) == 5 {
					engine, ok = parameters[4].(string)
					if !ok {
						return nil, errors.New("must pass a string or null to the third argument of Templates.render")
					}
				}
				switch engine {
				case "jinja":
					tpl, err := pongo2.FromFile(path)
					if err != nil {
						return nil, err
					}
					return tpl.Execute(goMap)
				default:
					return nil, errors.New("unknown template engine")
				}
			},
		}),
	}))

}

func wrenMapToGoMap(m *wren.MapHandle, k *wren.ListHandle) (map[string]interface{}, error) {
	var keys []string
	c, err := k.Count()
	if err != nil {
		return nil, err
	}
	for i := 0; i < c; i++ {
		key, err := k.Get(i)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key.(string))
	}
	goMap := make(map[string]interface{})
	for _, key := range keys {
		value, err := m.Get(key)
		if err != nil {
			return nil, err
		}
		goMap[key] = value
	}
	return goMap, nil
}
