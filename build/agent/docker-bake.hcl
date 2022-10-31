// docker-bake.hcl
target "docker-metadata-action" {}

group "default" {
    targets = ["custom", "cypress8", "cypress9", "cypress10"]
}

target "custom" {
  inherits = ["docker-metadata-action"]
  context = "./"
  dockerfile = "build/agent/Dockerfile.custom"
  platforms = [
    "linux/amd64",
    "linux/arm64"
  ]
}

target "cypress8" {
  inherits = ["docker-metadata-action"]
  context = "./"
  dockerfile = "build/agent/Dockerfile.cypress8"
  platforms = [
    "linux/amd64",
    "linux/arm64"
  ]
}

target "cypress9" {
  inherits = ["docker-metadata-action"]
  context = "./"
  dockerfile = "build/agent/Dockerfile.cypress9"
  platforms = [
    "linux/amd64",
    "linux/arm64"
  ]
}

target "cypress10" {
  inherits = ["docker-metadata-action"]
  context = "./"
  dockerfile = "build/agent/Dockerfile.cypress10"
  platforms = [
    "linux/amd64",
    "linux/arm64"
  ]
}
