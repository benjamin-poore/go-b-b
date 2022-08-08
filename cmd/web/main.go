package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/benjamin-poore/project/pkg/config"
	"github.com/benjamin-poore/project/pkg/handlers"
	"github.com/benjamin-poore/project/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application
func main() {
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create tempalte cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	handlers.Repo = &handlers.Repository{App: &app}

	render.NewTemplates(&app)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	fmt.Println(fmt.Sprintf("Starting on port %s", portNumber))

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
