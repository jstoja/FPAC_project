package main

import (
	"encoding/csv"
	"log"
	"os"
)

const (
	FP_FILE = "./data/FP_WIN.csv"
	AC_FILE = "./data/AC.csv"
)

type CodeLabel struct {
	code  string
	label string
}
type Root struct {
	name    string
	serials map[string]*Serial
}

type ProblemEnv struct {
	articles  map[string]*Article
	causes    map[string]*Cause
	functions map[string]*Function
	failures  map[string]*Failure
}

type Serial struct {
	ProblemEnv
	name string
}

type Article struct {
	CodeLabel
	failures map[string]*Failure
	causes   map[string]*Cause
}

type Cause struct {
	CodeLabel
	articles map[string]*Article
}

type Function struct {
	CodeLabel
	failures map[string]*Failure
}

type Failure struct {
	CodeLabel
	articles  map[string]*Article
	functions map[string]*Function
}

type Env struct {
	roots map[string]*Root

	articles  map[string]*Article
	causes    map[string]*Cause
	functions map[string]*Function
	failures  map[string]*Failure
}

func main() {
	env := &Env{
		make(map[string]*Root),
		make(map[string]*Article),
		make(map[string]*Cause),
		make(map[string]*Function),
		make(map[string]*Failure)}

	file, err := os.Open(FP_FILE)
	if err != nil {
		log.Fatal(err)
	}

	loadFP(file, env)

	file, err = os.Open(AC_FILE)
	if err != nil {
		log.Fatal(err)
	}

	loadAC(file, env)
}

func loadFP(fp_file *os.File, env *Env) {
	csv_reader := csv.NewReader(fp_file)
	csv_reader.Comma = ';'
	for {
		line, err := csv_reader.Read()
		if err != nil {
			log.Print(err)
			return
		}
		saveFPLine(line, env)
	}
}

func (env *Env) GetRoot(name string) *Root {
	root, exists := env.roots[name]
	if exists == false {
		root = &Root{name, make(map[string]*Serial)}
		env.roots[name] = root
	}
	return root
}

func (root *Root) GetSerial(name string) *Serial {
	serial, exists := root.serials[name]
	if exists == false {
		serial = &Serial{ProblemEnv{
			make(map[string]*Article),
			make(map[string]*Cause),
			make(map[string]*Function),
			make(map[string]*Failure)}, name}
		root.serials[name] = serial
	}
	return serial
}

func (env *Env) GetFunction(code string, label string) *Function {
	function, exists := env.functions[code]
	if exists == false {
		function = &Function{CodeLabel{code, label}, make(map[string]*Failure)}
		env.functions[code] = function
	}
	return function
}

func (env *Env) GetFailure(code string, label string) *Failure {
	failure, exists := env.failures[code]
	if exists == false {
		failure = &Failure{CodeLabel{code, label}, make(map[string]*Article), make(map[string]*Function)}
		env.failures[code] = failure
	}
	return failure
}

func (serial *Serial) AddFunction(function *Function) {
	_, exists := serial.functions[function.code]
	if exists == false {
		serial.functions[function.code] = function
	}
}

func (serial *Serial) AddFailure(failure *Failure) {
	_, exists := serial.failures[failure.code]
	if exists == false {
		serial.failures[failure.code] = failure
	}
}

func (function *Function) LinkFailure(failure *Failure) {
	_, exists := function.failures[failure.code]
	if exists == false {
		function.failures[failure.code] = failure
	}
}

func (failure *Failure) LinkFunction(function *Function) {
	_, exists := failure.functions[function.code]
	if exists == false {
		failure.functions[function.code] = function
	}
}

func saveFPLine(line []string, env *Env) {
	root := env.GetRoot(line[0])
	serial := root.GetSerial(line[1])
	function := env.GetFunction(line[2], line[3])
	failure := env.GetFailure(line[4], line[5])

	serial.AddFunction(function)
	serial.AddFailure(failure)

	function.LinkFailure(failure)
	failure.LinkFunction(function)
}

func loadAC(fp_file *os.File, env *Env) {
	csv_reader := csv.NewReader(fp_file)
	csv_reader.Comma = ';'
	for {
		line, err := csv_reader.Read()
		if err != nil {
			log.Print(err)
			return
		}
		saveACLine(line, env)
	}
}

func (env *Env) GetArticle(code string, label string) *Article {
	article, exists := env.articles[code]
	if exists == false {
		article = &Article{CodeLabel{code, label}, make(map[string]*Failure), make(map[string]*Cause)}
		env.articles[code] = article
	}
	return article
}

func (env *Env) GetCause(code string, label string) *Cause {
	cause, exists := env.causes[code]
	if exists == false {
		cause = &Cause{CodeLabel{code, label}, make(map[string]*Article)}
		env.causes[code] = cause
	}
	return cause
}

func (serial *Serial) AddArticle(article *Article) {
	_, exists := serial.articles[article.code]
	if exists == false {
		serial.articles[article.code] = article
	}
}

func (serial *Serial) AddCause(cause *Cause) {
	_, exists := serial.causes[cause.code]
	if exists == false {
		serial.causes[cause.code] = cause
	}
}

func (cause *Cause) LinkArticle(article *Article) {
	_, exists := cause.articles[article.code]
	if exists == false {
		cause.articles[article.code] = article
	}
}

func (article *Article) LinkCause(cause *Cause) {
	_, exists := article.causes[cause.code]
	if exists == false {
		article.causes[cause.code] = cause
	}
}

func saveACLine(line []string, env *Env) {
	root := env.GetRoot(line[0])
	serial := root.GetSerial(line[1])
	article := env.GetArticle(line[2], line[3])
	cause := env.GetCause(line[4], line[5])

	serial.AddArticle(article)
	serial.AddCause(cause)

	cause.LinkArticle(article)
	article.LinkCause(cause)
}
