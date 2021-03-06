package app

import (
	"crypto/tls"
	"errors"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"

	"github.com/moocss/chi-webserver/src/config"
	"github.com/moocss/chi-webserver/src/dao"
	"github.com/moocss/chi-webserver/src/pkg/log"
	"github.com/moocss/chi-webserver/src/router"
	"github.com/moocss/chi-webserver/src/service"
	"github.com/moocss/chi-webserver/src/util"
)

// App 项目
type App struct {
	config  *config.Config
	dao     *dao.Dao
	service service.Service
}

// New 实例化App
func New(cfg *config.Config, dao *dao.Dao, svc service.Service) *App {
	return &App{
		config:  cfg,
		dao:     dao,
		service: svc,
	}
}

// InitLog 初始化日志配置
func (app *App) InitLog() {
	if app.config.Log.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if app.config.Log.Trace {
		logrus.SetLevel(logrus.TraceLevel)
	}
	if app.config.Log.Text {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:   app.config.Log.Color,
			DisableColors: !app.config.Log.Color,
		})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: app.config.Log.Pretty,
		})
	}
}

// RunHTTPServer provide run http or https protocol.
func (app *App) RunHTTPServer() (err error) {
	if !app.config.Core.Enabled {
		log.Debug("httpd app is disabled.")
		return nil
	}

	if app.config.Core.AutoTLS.Enabled {
		return app.autoTLSServer()
	} else if app.config.Core.TLS.CertPath != "" && app.config.Core.TLS.KeyPath != "" {
		return app.defaultTLSServer()
	} else {
		return app.defaultServer()
	}
}

func (app *App) autoTLSServer() error {
	var g errgroup.Group

	dir := util.CacheDir()
	_ = os.MkdirAll(dir, 0700)

	manager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(app.config.Core.AutoTLS.Host),
		// Cache:      autocert.DirCache(app.config.Core.AutoTLS.Folder),
		Cache: autocert.DirCache(dir),
	}

	g.Go(func() error {
		return http.ListenAndServe(":http", manager.HTTPHandler(http.HandlerFunc(app.redirect)))
	})

	g.Go(func() error {
		serve := &http.Server{
			Addr: ":https",
			TLSConfig: &tls.Config{
				GetCertificate: manager.GetCertificate,
				NextProtos:     []string{"http/1.1"}, // disable h2 because Safari :(
			},
			Handler: serve(app),
		}
		app.handleSignal(serve)
		log.Info("Start to listening the incoming requests on https address")
		return serve.ListenAndServeTLS("", "")
	})

	return g.Wait()
}

func (app *App) defaultTLSServer() error {
	var g errgroup.Group
	g.Go(func() error {
		return http.ListenAndServe(":http", http.HandlerFunc(app.redirect))
	})
	g.Go(func() error {
		serve := &http.Server{
			Addr:    ":https",
			Handler: serve(app),
			TLSConfig: &tls.Config{
				NextProtos: []string{"http/1.1"}, // disable h2 because Safari :(
			},
		}
		app.handleSignal(serve)
		log.Infof("Start to listening the incoming requests on https address: %s", app.config.Core.TLS.Port)
		return serve.ListenAndServeTLS(
			app.config.Core.TLS.CertPath,
			app.config.Core.TLS.KeyPath,
		)
	})
	return g.Wait()
}

func (app *App) defaultServer() error {
	serve := &http.Server{
		Addr:    "0.0.0.0:" + app.config.Core.Port,
		Handler: serve(app),
	}
	app.handleSignal(serve)
	log.Infof("Start to listening the incoming requests on http address: %s", app.config.Core.Port)
	return serve.ListenAndServe()
}

// redirect ...
func (app *App) redirect(w http.ResponseWriter, req *http.Request) {
	var serverHost = app.config.Core.Host
	serverHost = strings.TrimPrefix(serverHost, "http://")
	serverHost = strings.TrimPrefix(serverHost, "https://")
	req.URL.Scheme = "https"
	req.URL.Host = serverHost

	w.Header().Set("Strict-Transport-Security", "max-age=31536000")

	http.Redirect(w, req, req.URL.String(), http.StatusMovedPermanently)
}

// serve returns a app instance
func serve(app *App) *chi.Mux {
	// Set gin mode.
	setRuntimeMode(app.config.Core.Mode)

	// Setup the app
	handler := router.Load(
		// Services
		app.service,

		// Middlwares

	)

	return handler
}

// setRuntimeMode 设置开发模式
func setRuntimeMode(mode string) {
	switch mode {
	case "dev":
		gin.SetMode(gin.DebugMode)
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		panic("unknown mode")
	}
}

// handleSignal handles system signal for graceful shutdown.
func (app *App) handleSignal(server *http.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		s := <-c
		log.Infof("got signal [%s], exiting apiserver now", s)
		if err := server.Close(); nil != err {
			log.Errorf("server close failed: %s ", err)
		}

		// 退出服务时，关闭数据库
		app.dao.DB.Close()

		log.Info("WebServer exited")
		os.Exit(0)
	}()
}

// PingServer 服务心跳检查
func (app *App) PingServer() (err error) {
	maxPingConf := app.config.Core.MaxPingCount
	for i := 0; i < maxPingConf; i++ {
		// Ping the app by sending a GET request to `/health`.
		resp, err := http.Get("http://localhost:" + app.config.Core.Port + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	err = errors.New("Cannot connect to the router.")
	return
}
