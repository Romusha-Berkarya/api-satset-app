package main

import (
	"encoding/json"
	"gateway-api-satset/entity"
	"gateway-api-satset/helper"
	"gateway-api-satset/middleware"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var ConnectMySql *gorm.DB

func main() {
	log.SetOutput(os.Stdout)

	viper.AddConfigPath("./config") //Viper looks here for the files.
	viper.SetConfigType("yaml")     //Sets the format of the config file.
	viper.SetConfigName("config")   // So that Viper loads default.yml.
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Warning could not load configuration: %v", err)
	}
	viper.AutomaticEnv() // Merges any overrides set through env vars.

	gatewayConfig := &entity.GatewayConfig{}

	err = viper.UnmarshalKey("gateway", gatewayConfig)
	if err != nil {
		panic(err)
	}

	helper.ConnectMysql(gatewayConfig.MySql)

	log.Println("Initializing routes...")

	r := mux.NewRouter()

	for _, route := range gatewayConfig.Routes {
		// Returns a proxy for the target url.
		proxy, err := NewProxy(route.Target)
		if err != nil {
			panic(err)
		}
		// Just logging the mapping.
		log.Printf("Mapping '%v' | %v ---> %v", route.Name, route.Context, route.Target)
		// Maps the HandlerFunc fn returned by NewHandler() fn
		// that delegates the requests to the proxy.
		r.HandleFunc(route.Context+"/{targetPath:.*}", NewHandler(proxy))
	}

	// Handle Cors
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}).Handler(r)

	log.Printf("Started server on %v", gatewayConfig.ListenAddr)
	log.Fatal(http.ListenAndServe(gatewayConfig.ListenAddr, corsHandler))
}

func NewProxy(targetUrl string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetUrl)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ModifyResponse = func(response *http.Response) error {
		dumpResponse, err := httputil.DumpResponse(response, false)
		if err != nil {
			return err
		}
		log.Println("Response: \r\n", string(dumpResponse))
		return nil
	}
	return proxy, nil
}

func NewHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		_, token := middleware.IsAuthorized(w, r)

		if token != nil {
			tokenEncode, err := json.Marshal(token)

			if err != nil {
				helper.ErrorWithMessage(err, "error token encode")
			}

			r.Header.Set("client-authorization", string(tokenEncode))
		}

		r.URL.Path = mux.Vars(r)["targetPath"]

		log.Println("Request URL: ", r.URL.String())
		log.Println("Method: ", r.Method)

		p.ServeHTTP(w, r)
	}
}
