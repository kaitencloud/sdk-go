module github.com/kaitencloud/sdk-go/ai/openai

go 1.25.0

require (
	github.com/kaitencloud/sdk-go v0.2.0
	github.com/openai/openai-go/v3 v3.0.0
)

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/oapi-codegen/runtime v1.3.1 // indirect
	github.com/tidwall/gjson v1.14.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
)

// Until github.com/kaitencloud/sdk-go v0.2.0 is tagged, the root module is the one in
// this repository. Drop this line before tagging ai/openai: a replace never reaches a
// dependent's build.
replace github.com/kaitencloud/sdk-go => ../..
