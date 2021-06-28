package server

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type api struct {
	router http.Handler
}

type PCA struct {
	Num_components int
	// contains filtered or unexported fields
}
type Server interface {
	Router() http.Handler
}

var Centroids []Node

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			next.ServeHTTP(w, req)
		})
}
func enableCORS(router *mux.Router) {
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}).Methods(http.MethodOptions)
	router.Use(middlewareCors)
}

var DataSetNodes []Node

func leerCSVdesdeURL(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func New() Server {

	a := &api{}
	r := mux.NewRouter()
	//Habilitamos los CORS
	enableCORS(r)
	url := "https://raw.githubusercontent.com/omarmendozaaa/IA-TF-BackEnd/master/server/CovidConPAC2.csv"

	fmt.Println("Se cargó satisfactoriamente el DataSet")

	DataSetLines, err := leerCSVdesdeURL(url)
	if err != nil {
		panic(err)
	}
	if err != nil {
		fmt.Println(err)
	}
	for _, line := range DataSetLines {

		param1, _ := strconv.ParseFloat(line[1], 64)
		param2, _ := strconv.ParseFloat(line[2], 64)

		var datita Node = Node{
			float64(param1),
			float64(param2),
		}
		DataSetNodes = append(DataSetNodes, datita)
	}

	//Train( data, clusters, iteraciones para definir centroide)
	_, Centroids = Train(DataSetNodes, 5, 50)

	r.HandleFunc("/gokmeans/predict", PredictKmeans).Methods("GET", "OPTIONS")
	r.HandleFunc("/gokmeans/centroids", GetCentroids).Methods("GET", "OPTIONS")

	//Iniciar Servidor
	a.router = r
	return a
}

func (a *api) Router() http.Handler {
	return a.router
}
