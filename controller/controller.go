package controller

import (
	"aplicacoes/projeto-zumbie/config"
	"encoding/json"
	"fmt"
	"net/http"
)

// APP ...
type APP struct {
	Versao    int64  `json:"versao"`
	Descricao string `json:"descrisao"`
	Data      string `json:"data"`
	Linguagem string `json:"linguagem"`
}

// Sobreviventes ...
type Sobreviventes struct {
	Codigo     uint        `json:"codigosobrevivente"`
	Nome       string      `json:"nome"`
	Idade      int         `json:""idade`
	Genero     string      `json:"genero"`
	Infectado  bool        `json:"infectado"`
	Inventario Inventarios `json:"inventario"`
}

// Inventarios ...
type Inventarios struct {
	Agua        int `json:"agua"`
	Comida      int `json:"comida"`
	Medicamento int `json:"medicamento"`
	Municao     int `json:"municao"`
}

var app []APP
var sobreviventes []Sobreviventes
var inventarios []Inventarios
var db = config.DB // var do banco
var versao int64
var descricao, data, linguagem string
var query string

// AtualizaAPP ...
func AtualizaAPP(v int64, d string, data string, l string) {
	app = append(app, APP{
		Versao:    v,
		Descricao: d,
		Data:      data,
		Linguagem: l,
	})
}

// HomeAPI ...
func HomeAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	versao = 1
	descricao = "Aplicação Sobrevivência Zumbi"
	data = "25/07/2018"
	linguagem = "Go"

	AtualizaAPP(versao, descricao, data, linguagem)

	json.NewEncoder(w).Encode(app)
}

// BuscarTodosSobrevivente ...
func BuscarTodosSobrevivente(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query = "SELECT sobrevivente.codigo, nome, idade, genero, infectado, " +
		"agua, comida, medicamento, municao " +
		"FROM sobrevivente, inventario" +
		" WHERE " +
		"sobrevivente.codigo = inventario.codigosobrevivente"

	rows, err := db.Query(query)
	CheckError(err)

	sobreviventes = sobreviventes[:0]

	for rows.Next() {
		sobrevivente := Sobreviventes{}
		rows.Scan(
			&sobrevivente.Codigo, &sobrevivente.Nome, &sobrevivente.Idade,
			&sobrevivente.Genero, &sobrevivente.Infectado, &sobrevivente.Inventario.Agua,
			&sobrevivente.Inventario.Comida, &sobrevivente.Inventario.Medicamento, &sobrevivente.Inventario.Municao)
		sobreviventes = append(sobreviventes, sobrevivente)
	}

	json.NewEncoder(w).Encode(sobreviventes)

}

// AdicionarNovoSobrevivente ...
func AdicionarNovoSobrevivente(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sobreviventes = sobreviventes[:0]
	inventarios = inventarios[:0]

	sobrevivente := Sobreviventes{}
	_ = json.NewDecoder(r.Body).Decode(&sobrevivente)
	sobreviventes = append(sobreviventes, sobrevivente)

	nome := sobrevivente.Nome
	idade := sobrevivente.Idade
	genero := sobrevivente.Genero
	infectado := sobrevivente.Infectado
	agua := sobrevivente.Inventario.Agua
	comida := sobrevivente.Inventario.Comida
	medicamento := sobrevivente.Inventario.Medicamento
	municao := sobrevivente.Inventario.Municao

	sobreviventes = sobreviventes[:0]
	inventarios = inventarios[:0]

	querySobrevivente := "INSERT INTO sobrevivente (nome,idade,genero,infectado) VALUES(?,?,?,?)"
	queryInventario := "INSERT INTO inventario (codigosobrevivente, agua, comida, medicamento, municao) VALUES (LAST_INSERT_ID(), ?, ?, ?, ?);"

	fmt.Println(query)

	stmt, err := db.Prepare(querySobrevivente)
	CheckError(err)
	_, err = stmt.Exec(nome, idade, genero, infectado)
	CheckError(err)

	stmt, err = db.Prepare(queryInventario)
	CheckError(err)
	_, err = stmt.Exec(agua, comida, medicamento, municao)

}

// CheckError ...
func CheckError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
