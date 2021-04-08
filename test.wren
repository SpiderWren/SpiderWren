import "web" for Routes, App

Routes.GET("/") {
    return "hello"
}

Routes.GET("/twoplustwo") {
    return "%(2 + 2)"
}

App.run(6969)