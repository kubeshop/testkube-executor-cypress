module github.com/kubeshop/kubtest-executor-cypress

go 1.16

replace github.com/kubeshop/kubtest v0.0.0-20210901115324-4505d1df0c19 => ../kubtest

require (
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kubeshop/kubtest v0.0.0-20210901115324-4505d1df0c19
	github.com/stretchr/testify v1.7.0
)
