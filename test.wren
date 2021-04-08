import "web" for Routes, App

Routes.GET("/") {
    return "hello"
}

App.run(6969)