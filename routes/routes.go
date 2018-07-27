package routes

import (
	"aplicacoes/projeto-zumbie/config"
	"aplicacoes/projeto-zumbie/controller"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var portaAplicacao string

// HandleFunc ...
func HandleFunc() {
	rotas := mux.NewRouter()
	config.TryConn()
	portaAplicacao = ":3000"
	fmt.Println("Aplicação ON: porta => ", portaAplicacao)

	rotas.HandleFunc("/api/", controller.HomeAPI).Methods("GET")
	rotas.HandleFunc("/api/sobreviventes", controller.BuscarTodosSobrevivente).Methods("GET")
	rotas.HandleFunc("/api/adicionar/sobrevivente", controller.AdicionarNovoSobrevivente).Methods("POST")
	rotas.HandleFunc("/api/sobrevivente/{sobrevivente1}/{sobrevivente2}", controller.BuscarSobreviventes).Methods("GET")
	rotas.HandleFunc("/api/trocar", controller.RealizarTroca).Methods("POST")

	log.Fatal(http.ListenAndServe(portaAplicacao, rotas))
}
