package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
	mux := pat.New()

	dynamicMiddleware := func(next http.Handler) http.Handler {
		return app.session.Enable(noSurf(app.authenticate(next)))
	}

	mux.Get("/", dynamicMiddleware(http.HandlerFunc(app.home)))
	mux.Get("/snippet/create", dynamicMiddleware(app.requireAuthentication(http.HandlerFunc(app.createSnippetForm))))
	mux.Post("/snippet/create", dynamicMiddleware(app.requireAuthentication(http.HandlerFunc(app.createSnippet))))
	mux.Get("/snippet/:id", dynamicMiddleware(http.HandlerFunc(app.showSnippet)))

	mux.Get("/user/signup", dynamicMiddleware(http.HandlerFunc(app.signupUserForm)))
	mux.Post("/user/signup", dynamicMiddleware(http.HandlerFunc(app.signupUser)))
	mux.Get("/user/login", dynamicMiddleware(http.HandlerFunc(app.loginUserForm)))
	mux.Post("/user/login", dynamicMiddleware(http.HandlerFunc(app.loginUser)))
	mux.Post("/user/logout", dynamicMiddleware(app.requireAuthentication(http.HandlerFunc(app.logoutUser))))

	mux.Get("/ping", http.HandlerFunc(ping))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
