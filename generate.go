//go:generate go get -v github.com/mailru/easyjson/...
//go:generate easyjson -all -pkg framework/agent
//go:generate easyjson -all -disallow_unknown_fields -pkg httpapi/models
//go:generate easyjson -disallow_unknown_fields -pkg .

package dlframework
