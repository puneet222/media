package main

func (app *App) InitializeRoutes() {
	app.Router.HandleFunc("/ws", app.createWebsocket)
	app.Router.HandleFunc("/data", app.GetData)
}
