package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/seesawhq/seesaw/models"
	"github.com/seesawhq/seesaw/views/home"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("seesaw.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.User{}, &models.Team{}, &models.Member{}, &models.FeatureFlag{})
	ctx := context.Background()
	_, err = gorm.G[models.User](db).Where("email = ?", "sagar.kaurav@gmail.com").First(ctx)
	if err != nil {
		_ = gorm.G[models.User](db).Create(ctx, &models.User{FirstName: "sagar", LastName: "kaurav", Email: "sagar.kaurav@gmail.com"})
	}
	mux := http.NewServeMux()
	home.New(mux)
	fmt.Println("started server at 0.0.0.0:3000")
	http.ListenAndServe("0.0.0.0:3000", mux)
}
