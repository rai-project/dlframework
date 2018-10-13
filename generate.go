//go:generate go get -v github.com/mailru/easyjson/...
//go:generate easyjson -snake_case -all -pkg framework/agent
//go:generate easyjson -snake_case -all -disallow_unknown_fields -pkg httpapi/models
//go:generate easyjson -snake_case -disallow_unknown_fields -pkg .

package dlframework
