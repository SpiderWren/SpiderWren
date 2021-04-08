package web

import (
	wren "github.com/crazyinfin8/WrenGo"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strconv"
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

					callHandle, err := handle.Func("call()")
					if err != nil {
						log.Fatal("Must pass a handle with 0 parameters to the second argument of Routes.GET")
					}
					result, err := callHandle.Call()
					if err != nil {
						context.Header("Content-Type", "text/html")
						context.String(500,"An error occurred: {{.error}}", gin.H{
							"error": err.Error(),
						})
						return
					}
					out, ok := result.(string)
					if !ok {
						log.Fatal("Must return a string")
					}
					context.Header("Content-Type","text/html")
					context.String(200, out, nil)

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
				err := app.Router.Run("localhost:"  + strconv.Itoa(port))
				return nil, err
			},
		}),
	}))

}
