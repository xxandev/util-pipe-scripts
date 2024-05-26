package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"util-pipe/internal/dbg"
	"util-pipe/internal/utils"
	"util-pipe/internal/xj"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

const (
	APP_VERSION       float64 = 1.1
	APP_LOG_SEPARATOR string  = "\n====================================================="
)

//go:embed repository
var rep embed.FS

func init() {
	log.SetPrefix("[PIPE] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmsgprefix)

	flag.Usage = func() { fmt.Println(APP_LOG_SEPARATOR); flag.PrintDefaults() }

	flag.StringVar(&config.Server.Host, "host", ":8090", "server host, example: 127.0.0.1:8080, :8080...")
	flag.StringVar(&config.ScriptsPath, "scripts-path", "", "scripts path, example: /root/.scripts or c:\\.scripts")
	flag.StringVar(&config.WikiPath, "wiki-path", "", "wiki path, example: /root/.wiki or c:\\.wiki")
	flag.StringVar(&config.BasicAuth.Login, "user", "", "basic auth user login, example: user")
	flag.StringVar(&config.BasicAuth.Pass, "pass", "", "basic auth user password, example: qwerty")
	flag.StringVar(&config.LDAP.URL, "ldap-url", "", "ldap url, example: ldap://ldap.example.com:389")
	flag.StringVar(&config.LDAP.DN, "ldap-dn", "", "ldap dn, example: ou=users,dc=ldap,dc=example,dc=com")
	flag.StringVar(&config.Server.SSL.CRT, "ssl-crt", "", "ssl crt, example: /root/.certs/server.crt")
	flag.StringVar(&config.Server.SSL.Key, "ssl-key", "", "ssl key, example: /root/.certs/server.key")

	flag.BoolFunc("mute", "deactivate log", func(string) error { log.SetOutput(io.Discard); return nil })
	flag.BoolFunc("debug", "activate debug", func(string) error { dbg.Log.SetOutput(os.Stdout); return nil })
	flag.BoolFunc("v", "app version", func(string) error { fmt.Println(APP_VERSION); os.Exit(0); return nil })
	flag.BoolFunc("version", "app version", func(string) error { fmt.Println(APP_VERSION); os.Exit(0); return nil })

	flag.Func("c", "path to configuration file", func(p string) error { return config.Init(p) })
	flag.Func("config", "path to configuration file", func(p string) error { return config.Init(p) })
	flag.Func("example", "bla bla bla bla blas", func(f string) error { fmt.Println(config.Example(f)); os.Exit(0); return nil })

	flag.Parse()

	if err := config.Check(); err != nil {
		log.Fatalf("error check config: %v\n", err)
	}
	if err := authorization.Init(); err != nil {
		log.Fatalf("error init auth: %v\n", err)
	}
}

func main() {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
		Debug:            false,
	}))
	router.Use(mwWhereAmI)

	if err := InitWiki(router, config.WikiPath); err != nil {
		log.Printf("init wiki: %v\n", err)
	}

	pageIndex := InitEmbedPage(rep, "repository/front/index.html")
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		pageIndex.Execute(w, config.wiki.ListLinks)
	})

	pageNotFound := InitEmbedPage(rep, "repository/front/not-found.html")
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		pageNotFound.Execute(w, nil)
	})

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("pong")) })
	router.Get("/str2md5", hGetStr2MD5)
	router.Get("/statistic", hGetStatistic)
	router.With(mwAuth).Get("/check-auth", hGetCheckAuth)
	router.With(mwAuth).Get("/script-exec", hGetScriptExec)
	router.With(mwAuth).Post("/script-exec", hPostScriptExec)
	router.With(mwAuth).Get("/json-read", hGetJsonRead)
	router.With(mwAuth).Post("/json-create", hPostJsonCreate)
	router.With(mwAuth).Get("/gotty-run", hGetGottyRun)
	router.With(mwAuth).Get("/gotty-kill", hGetGottyKill)

	server := &http.Server{Addr: config.Server.Host, Handler: router}
	utils.OnTermination(func() { server.Close() })
	if len(config.Server.SSL.CRT)+len(config.Server.SSL.Key) > 0 {
		log.Println("server https run, host:", config.Server.Host)
		log.Fatalf("server: %v", server.ListenAndServeTLS(config.Server.SSL.CRT, config.Server.SSL.Key))
	}
	log.Println("server http run, host:", config.Server.Host)
	log.Fatalf("server: %v", server.ListenAndServe())
}

func InitEmbedPage(fs embed.FS, name string) *template.Template {
	cont, err := rep.ReadFile(name)
	if err != nil {
		return errTmpPage(name, err)
	}
	tmp, err := template.New(name).Parse(string(cont))
	if err != nil {
		return errTmpPage(name, err)
	}
	return tmp
}

func errTmpPage(name string, err error) *template.Template {
	tmp, _ := template.New(name).Parse(xj.Errf("error page - %v", err).Str())
	return tmp
}
