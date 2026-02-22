variable "url" {
  type    = string
  default = getenv("DB_URL")
}

variable "dev_url" {
  type    = string
  default = getenv("DB_DEV_URL")
}

data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "cmd/migrate/main.go",
  ]
}

env "local" {
  src = data.external_schema.gorm.url
  dev = var.dev_url
  url = var.url
  migration {
    dir = "file://migrations?format=golang-migrate"
  }
}
