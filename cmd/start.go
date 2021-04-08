package cmd

import (
	_ "embed"
	"github.com/gin-gonic/gin"

	wren "github.com/crazyinfin8/WrenGo"
	"github.com/nailuj29gaming/wren-web/utils"
	"github.com/nailuj29gaming/wren-web/web"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a server to run your app",
	Long: `Starts a development server on port 3333 by default. 
This can be changed with --port=3000, for example`,
	Run: func(cmd *cobra.Command, args []string) {
		debugMode, err := cmd.Flags().GetBool("debug")
		if err != nil {
			log.Fatal("Could not get debug flag. This shouldn't ever happen")
		}
		if debugMode {
			log.SetLevel(log.DebugLevel)
		}
		log.Debug("Initializing VM")
		cfg := wren.NewConfig()
		cfg.LoadModuleFn = func(vm *wren.VM, name string) (source string, ok bool) {
			if name == "web" {
				return utils.WebModuleSource, true
			}
			return "", false
		}
		vm := cfg.NewVM()
		defer vm.Free()
		app := web.App{
			Router: gin.Default(),
		}
		log.Debug("Creating foreign classes")
		web.CreateForeignClasses(vm, &app)
		err = vm.InterpretFile(args[0])
		if err != nil {
			log.Fatalf("An error occurred while running %s: %s", args[0], err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
