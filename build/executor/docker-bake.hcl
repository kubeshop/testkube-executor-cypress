// docker-bake.hcl
target "docker-metadata-action" {}

target "build" {
  inherits = ["docker-metadata-action"]
  context = "./"
  dockerfile = "build/executor/Dockerfile"
  platforms = [
    "linux/amd64",
  ]
}