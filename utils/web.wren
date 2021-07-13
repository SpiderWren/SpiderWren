foreign class Routes {
    foreign static GET(path, callback)
}

foreign class App {
    foreign static run(port)
}

foreign class TemplatesHelper {
    foreign static render(path, data, keys)
}

class Templates {
    static render(path, data) {
        return TemplatesHelper.render(path, data, data.keys.toList)
    }
}