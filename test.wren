import "web" for Routes, App

Routes.GET("/") {
    return "hello"
}

Routes.GET("/twoplustwo") {
    return "%(2 + 2)"
}

Routes.GET("/param/:param") { | params |
    return params["param"]
}

App.run(3000)