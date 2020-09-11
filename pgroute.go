package pgroute

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func MountFunctionRoute(db *gorm.DB) http.Handler {
	r := chi.NewRouter()

	r.Post("/{funcName}", func(w http.ResponseWriter, r *http.Request) {
		funcName := chi.URLParam(r, "funcName")
		argList, err := getFunctionArgs(db, funcName)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		paramMap, err := getParamsFromReq(r)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		paramList := getParamList(argList, paramMap)

		x := make([]map[string]interface{}, 0)
		db.Raw(buildSqlStmt(funcName, argList), paramList...).Scan(&x)

		res, _ := json.Marshal(x)
		w.Write(res)
	})
	return r
}

func getFunctionArgs(db *gorm.DB, name string) ([]string, error) {
	var res struct {
		A string
	}
	err := db.Raw("SELECT pg_get_function_arguments(oid) AS a FROM pg_proc p WHERE p.proname = ?;", name).Scan(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no such function")
		}
		return nil, err
	}

	if res.A == "" {
		return []string{}, nil
	}

	// get arguments
	argList := strings.Split(res.A, ",")

	// clean argList
	for i, a := range argList {
		a = strings.TrimSpace(a)
		a = strings.Split(a, " ")[0]

		argList[i] = a
	}

	return argList, nil
}

func getParamsFromReq(r *http.Request) (out map[string]interface{}, err error) {
	d := json.NewDecoder(r.Body)
	d.UseNumber()
	d.Decode(&out)

	return out, nil
}

func getParamList(argList []string, paramMap map[string]interface{}) (out []interface{}) {
	for _, a := range argList {
		val, ok := paramMap[a].(json.Number)
		if !ok {
			out = append(out, paramMap[a])
			continue
		}
		if i, err := val.Int64(); err == nil {
			out = append(out, i)
			continue
		}
		if f, err := val.Float64(); err == nil {
			out = append(out, f)
			continue
		}
	}

	return out
}

func buildSqlStmt(funcName string, argList []string) string {
	params := ""

	if len(argList) > 0 {
		for i, a := range argList {
			argList[i] = fmt.Sprintf("%v => ?", a)
		}

		params = strings.Join(argList, ", ")
	}

	return fmt.Sprintf("SELECT * FROM %v(%v)", funcName, params)
}
