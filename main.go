package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/seesawhq/seesaw/models"
	"github.com/seesawhq/seesaw/views/auth"
	"github.com/seesawhq/seesaw/views/home"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed static
var StaticFiles embed.FS

//go:embed templates
var TemplateFS embed.FS

func main() {
	db, err := gorm.Open(sqlite.Open("seesaw.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.User{}, &models.Team{}, &models.Member{}, &models.FeatureFlag{})
	mux := http.NewServeMux()
	fs := http.FileServer(http.FS(StaticFiles))
	mux.Handle("/static/", fs)
	home.New(TemplateFS, mux)
	auth.New(TemplateFS, mux)
	fmt.Println("started server at 0.0.0.0:3000")
	http.ListenAndServe("0.0.0.0:3000", mux)
}
