package endpoints

import "strings"

type Endpoint = string

const (
	Secrets Endpoint = "/secrets"
)

func Join(server string, endpoint Endpoint) string {
	return strings.Join([]string{server, endpoint}, "")
}
