module github.com/kubeshop/kubtest-executor-cypress

go 1.16

replace github.com/kubeshop/kubtest v0.0.0-20210823141506-ac90beb1ff74 => ../kubtest

require (
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kubeshop/kubtest v0.0.0-20210823141506-ac90beb1ff74
	github.com/stretchr/testify v1.7.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
