job "prospector-[[.environment_slug]]" {
  type = "service"

  datacenters = ["dc1"]

  constraint {
    attribute = "${attr.cpu.arch}"
    value     = "amd64"
  }

  meta {
    git_sha = "[[.git_sha]]"
  }

  group "prospector-[[.environment_slug]]" {
    count = 1

    network {
      port "api" {
        to = 8080
      }
      port "http" {
        to = 80
      }
    }

    service {
      name = "prospector-review-[[.environment_slug]]"
      port = "http"

      check {
        name     = "global_check"
        type     = "http"
        interval = "10s"
        timeout  = "2s"
        path     = "/"
      }
    }

    task "prospector-api" {
      driver = "docker"

      config {
        image = "git.dbyte.xyz/distro/prospector/api:[[.git_sha]]"
        ports = ["api"]
      }

      service {
        name = "prospector-api-[[.environment_slug]]"
        port = "api"

        check {
          name     = "api_check"
          type     = "http"
          interval = "10s"
          timeout  = "2s"
          path     = "/api/health"
        }

        tags = [
          "traefik.enable=true",
          "traefik.http.routers.prospector-api-[[.environment_slug]].rule=Host(`[[.deploy_url]]`) && PathPrefix(`/api`)",
          "traefik.http.routers.prospector-api-[[.environment_slug]].entrypoints=websecure",
          "traefik.http.routers.prospector-api-[[.environment_slug]].tls.certresolver=lets-encrypt"
        ]
      }

      resources {
        cpu    = 128
        memory = 128
      }
    }

    task "prospector-frontend" {
      driver = "docker"

      config {
        image = "git.dbyte.xyz/distro/prospector/frontend:[[.git_sha]]"
        ports = ["http"]
      }

      service {
        name = "prospector-frontend-[[.environment_slug]]"
        port = "http"

        tags = [
          "traefik.enable=true",
          "traefik.http.routers.prospector-frontend-[[.environment_slug]].rule=Host(`[[.deploy_url]]`)",
          "traefik.http.routers.prospector-frontend-[[.environment_slug]].entrypoints=websecure",
          "traefik.http.routers.prospector-frontend-[[.environment_slug]].tls.certresolver=lets-encrypt"
        ]
      }

      resources {
        cpu    = 128
        memory = 128
      }
    }
  }
}
