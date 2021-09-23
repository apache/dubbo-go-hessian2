module hessian_test

go 1.17

require (
	dup_struct_name v0.0.0
	github.com/apache/dubbo-go-hessian2 v0.0.0
	github.com/stretchr/testify v1.4.0
)

require (
	github.com/davecgh/go-spew v1.1.0
	github.com/dubbogo/gost v1.9.0
	github.com/pkg/errors v0.9.1
	github.com/pmezard/go-difflib v1.0.0
	gopkg.in/yaml.v2 v2.2.2
)

replace (
	dup_struct_name => ./dup_struct_name_demo
	github.com/apache/dubbo-go-hessian2 => ../
)
