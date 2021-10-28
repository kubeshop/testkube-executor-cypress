module github.com/kubeshop/testkube-executor-cypress

go 1.16

replace github.com/kubeshop/testkube => /root/go/src/github.com/kubeshop/testkube

require (
	github.com/joshdk/go-junit v0.0.0-20210226021600-6145f504ca0d
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kubeshop/testkube v0.6.4-0.20211028114358-8fb98d0667a3
)
