package logic

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/a-h/templ"
)

func BuildFilterURL(
	param, value, q, status, domain string,
	difficulty int,
	sort string,
) templ.SafeURL {
	u, _ := url.Parse("/exercises")
	query := u.Query()

	// ✅ Param modifié (reset si value vide)
	if value != "" {
		query.Set(param, value)
	}

	// ✅ Conserver les autres filtres actifs
	if q != "" {
		query.Set("q", q)
	}
	if status != "" && param != "status" {
		query.Set("status", status)
	}
	if domain != "" && param != "domain" {
		query.Set("domain", domain)
	}
	if difficulty > 0 && param != "difficulty" {
		query.Set("difficulty", strconv.Itoa(difficulty))
	}
	if sort != "" && param != "sort" {
		query.Set("sort", sort)
	}

	u.RawQuery = query.Encode()
	return templ.SafeURL(u.String())
}

// buildHxVals (inchangé)
func BuildHxVals(param, value, q, status, domain string, difficulty int, sort string) string {
	vals := make(map[string]string)
	vals[param] = value

	if q != "" {
		vals["q"] = q
	}
	if status != "" && param != "status" {
		vals["status"] = status
	}
	if domain != "" && param != "domain" {
		vals["domain"] = domain
	}
	if difficulty > 0 && param != "difficulty" {
		vals["difficulty"] = strconv.Itoa(difficulty)
	}
	if sort != "" && param != "sort" {
		vals["sort"] = sort
	}

	json := "{"
	first := true
	for k, v := range vals {
		if !first {
			json += ","
		}
		json += fmt.Sprintf(`"%s":"%s"`, k, v)
		first = false
	}
	json += "}"

	return json
}
