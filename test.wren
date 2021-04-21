import "web" for Routes, App

Routes.GET("/") {
    return "hello"
}

Routes.GET("/add/:num1/:num2") { | params |
    var num1 = Num.fromString(params["num1"])
    var num2 = Num.fromString(params["num2"])


    return "%(num1) + %(num2) = %(num1 + num2)"
}

Routes.GET("/param/:param") { | params |
    return params["param"]
}

App.run(3000)