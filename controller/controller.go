package controller

import (
	"aplicacoes/projeto-zumbie/config"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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
	Idade      int         `json:"idade"`
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

// Trocas ...
type Trocas struct {
	Sobrevivente1 Sobreviventes
	Sobrevivente2 Sobreviventes
}

// ErroTroca ...
type ErroTroca struct {
	NomeSobrevivente string `json:"nomesobrevivente"`
	Mensagem         string `json:"mensagem"`
}

var app []APP
var sobreviventes []Sobreviventes
var inventarios []Inventarios
var errotroca []ErroTroca
var inventarioS1 int
var inventarioS2 int
var trocas []Trocas
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

	stmt, err := db.Prepare(querySobrevivente)
	CheckError(err)
	_, err = stmt.Exec(nome, idade, genero, infectado)
	CheckError(err)

	stmt, err = db.Prepare(queryInventario)
	CheckError(err)
	_, err = stmt.Exec(agua, comida, medicamento, municao)

	json.NewEncoder(w).Encode("Sobrevivente adicionado !!")

}

// BuscarSobreviventes ...
func BuscarSobreviventes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sobrevivente1 := mux.Vars(r)["sobrevivente1"]
	sobrevivente2 := mux.Vars(r)["sobrevivente2"]

	query = "(SELECT sobrevivente.codigo, nome, idade, genero, infectado, " +
		"agua, comida, medicamento, municao " +
		"FROM sobrevivente, inventario" +
		" WHERE " +
		"sobrevivente.codigo = ? AND inventario.codigosobrevivente = ?) " +
		"UNION" +
		" (SELECT sobrevivente.codigo, nome, idade, genero, infectado, " +
		"agua, comida, medicamento, municao " +
		"FROM sobrevivente, inventario" +
		" WHERE " +
		"sobrevivente.codigo = ? AND inventario.codigosobrevivente = ?) "

	rows, err := db.Query(query, sobrevivente1, sobrevivente1, sobrevivente2, sobrevivente2)
	CheckError(err)

	sobreviventes = sobreviventes[:0]
	sobrevivente := Sobreviventes{}

	for rows.Next() {
		rows.Scan(&sobrevivente.Codigo, &sobrevivente.Nome, &sobrevivente.Idade,
			&sobrevivente.Genero, &sobrevivente.Infectado, &sobrevivente.Inventario.Agua,
			&sobrevivente.Inventario.Comida, &sobrevivente.Inventario.Medicamento,
			&sobrevivente.Inventario.Municao)
		sobreviventes = append(sobreviventes, sobrevivente)
	}

	json.NewEncoder(w).Encode(sobreviventes)

}

// RealizarTroca ...
func RealizarTroca(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	trocas = trocas[:0]
	errotroca = errotroca[:0]

	troca := Trocas{}
	erroTroca := ErroTroca{}
	// sobrevivente := Sobreviventes{}

	_ = json.NewDecoder(r.Body).Decode(&troca)
	trocas = append(trocas, troca)

	codigoSobrevivente1 := troca.Sobrevivente1.Codigo
	codigoSobrevivente2 := troca.Sobrevivente2.Codigo
	s1 := troca.Sobrevivente1
	s2 := troca.Sobrevivente2

	for _, s := range sobreviventes {
		if s.Codigo == codigoSobrevivente1 {
			// Aqui to tentando recuperar os valores dos alimentos...
			if s1.Inventario.Agua > 0 {
				fmt.Println("Agua será trocada !!")
				inventarioS1 = s1.Inventario.Agua
				if inventarioS1 > s.Inventario.Comida {
					erroTroca.NomeSobrevivente = s.Nome
					erroTroca.Mensagem = "Não possui água suficiente !! ;("
					json.NewEncoder(w).Encode(erroTroca)
				}
			}
			if s1.Inventario.Comida > 0 {
				fmt.Println("comida será trocada !!")
				inventarioS1 = s1.Inventario.Comida
				if inventarioS1 > s.Inventario.Comida {
					erroTroca.NomeSobrevivente = s.Nome
					erroTroca.Mensagem = "Não possui comida suficiente !! ;("
					json.NewEncoder(w).Encode(erroTroca)
				}
			}
			if s1.Inventario.Medicamento > 0 {
				fmt.Println("Medicamento será trocada !!")
				inventarioS1 = s1.Inventario.Medicamento
				if inventarioS1 > s.Inventario.Medicamento {
					erroTroca.NomeSobrevivente = s.Nome
					erroTroca.Mensagem = "Não possui medicamento suficiente !! ;("
					json.NewEncoder(w).Encode(erroTroca)
				}
			}
			if s1.Inventario.Municao > 0 {
				fmt.Println("Munição será trocada !!")
				inventarioS1 = s1.Inventario.Municao
				if inventarioS1 > s.Inventario.Municao {
					erroTroca.NomeSobrevivente = s.Nome
					erroTroca.Mensagem = "Não possui munição suficiente !! ;("
					json.NewEncoder(w).Encode(erroTroca)
				}
			}

		}
	}

	for _, s := range sobreviventes {
		if s.Codigo == codigoSobrevivente2 {
			// Aqui to tentando recuperar os valores dos alimentos...
			if s2.Inventario.Agua > 0 {
				fmt.Println("Agua será trocada !!")
				inventarioS2 = s2.Inventario.Agua
				if inventarioS2 > s.Inventario.Comida {
					erroTroca.NomeSobrevivente = s2.Nome
					erroTroca.Mensagem = "Não possui água suficiente !! ;("
					json.NewEncoder(w).Encode(erroTroca)
				}
			}
			if s2.Inventario.Comida > 0 {
				fmt.Println("comida será trocada !!")
				inventarioS2 = s2.Inventario.Comida
				if inventarioS2 > s.Inventario.Comida {
					erroTroca.NomeSobrevivente = s2.Nome
					erroTroca.Mensagem = "Não possui comida suficiente !! ;("
					json.NewEncoder(w).Encode(erroTroca)
				}
			}
			if s2.Inventario.Medicamento > 0 {
				fmt.Println("Medicamento será trocada !!")
				inventarioS2 = s2.Inventario.Medicamento
				if inventarioS2 > s.Inventario.Medicamento {
					erroTroca.NomeSobrevivente = s2.Nome
					erroTroca.Mensagem = "Não possui medicamento suficiente !! ;("
					json.NewEncoder(w).Encode(erroTroca)
				}
			}
			if s2.Inventario.Municao > 0 {
				fmt.Println("Munição será trocada !!")
				inventarioS2 = s2.Inventario.Municao
				if inventarioS2 > s.Inventario.Municao {
					erroTroca.NomeSobrevivente = s2.Nome
					erroTroca.Mensagem = "Não possui munição suficiente !! ;("
					json.NewEncoder(w).Encode(erroTroca)
				}
			}
		}
	}

	if ComparaTroca(trocas) {
		fmt.Println("ok vamos trocar...")
	} else {
		fmt.Println("não vamos.")
	}
}

// ComparaTroca ...
func ComparaTroca(trocas []Trocas) bool {
	// agua := 4
	// comida := 3
	// medicamento := 2
	// municao := 1
	var s1Agua int
	var s1Comida int
	var s1Med int
	var s1Mun int

	var s2Agua int
	var s2Comida int
	var s2Med int
	var s2Mun int

	var somaS1 int
	var somaS2 int

	for _, s := range trocas {
		s1Agua = s.Sobrevivente1.Inventario.Agua * 2
		s1Comida = s.Sobrevivente1.Inventario.Comida * 2
		s1Med = s.Sobrevivente1.Inventario.Medicamento * 2
		s1Mun = s.Sobrevivente1.Inventario.Municao * 2

		somaS1 = s1Agua + s1Comida + s1Med + s1Mun

		s2Agua = s.Sobrevivente2.Inventario.Agua * 2
		s2Comida = s.Sobrevivente2.Inventario.Comida * 2
		s2Med = s.Sobrevivente2.Inventario.Medicamento * 2
		s2Mun = s.Sobrevivente2.Inventario.Municao * 2

		somaS2 = s2Agua + s2Comida + s2Med + s2Mun

	}

	if somaS1 != somaS2 {
		return false
	}
	return true
}

// CheckError ...
func CheckError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
