package backend

import (
	"github.com/Masterminds/sprig"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	// log "github.com/sirupsen/logrus"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Use(sessions.Sessions("session", cookie.NewStore([]byte("secret"))))

	// FuncMapの設定
	r.SetFuncMap(sprig.FuncMap())

	r.Use(cors.New(cors.Config{
		AllowMethods: []string{
			"POST",
			"GET",
			"OPTIONS",
			"PUT",
			"DELETE",
		},
		AllowHeaders: []string{
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"Access-Control-Allow-Origin",
		},
		AllowOrigins: []string{
			"*",
		},
		MaxAge: 24 * time.Hour,
	}))

	r.Use(secure.New(secure.Config{
		ContentSecurityPolicy: "default-src '*';",
	}))

	// r.LoadHTMLGlob("frontend/build/*.html")
	// r.Use(static.Serve("/", static.LocalFile("frontend/build", false)))

	// templatesの読み込み
	// r.LoadHTMLGlob("backend/templates/*.html")

	// 静的ファイルシステムの設定
	r.StaticFS("/backend/static", http.Dir("backend/static"))

	// GET
	r.GET("/", HandleGetAll)
	r.GET("/page", HandleGetAll)
	r.GET("/user", HandleGetAllByUser)
	r.GET("/list", HandleGetMyShops)
	r.GET("/show", HandleGet)
	r.GET("/edit", HandleEdit)

	// POST
	r.POST("/confirm", HandleConfirm)
	r.POST("/post", HandlePost)
	r.POST("/put", HandlePut)
	r.POST("/merge", HandleMerge)
	r.POST("/delete", HandleDelete)

	return r
}
