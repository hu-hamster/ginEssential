package main

import (
	"hu-hamster/ginEssential/common"
	"hu-hamster/ginEssential/router"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {
	dir, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(dir + "/config")
	if err = viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}
	common.InitDB()
}

func main() {
	r := gin.Default()
	r = router.CollectRouter(r)
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	} else {
		panic(r.Run())
	}
}
