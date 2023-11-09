job "prospector-api" {
  type = "service"
  datacenters = ["dc1"]
  
  constraint {
    attribute = "${attr.cpu.arch}"
    value     = "amd64"
  }

  meta {
    git_sha = "[[.git_sha]]"
  }

  group "prospector-api" {
    count = 1

    network {
      port  "http"{
        to = 8080
      }
    }

    update {
      max_parallel = 1
      canary = 1
    }

    service {
      name = "prospector-api"
      port = "http"

      tags = [
        "traefik.enable=true",
        "traefik.http.routers.prospector-api.rule=Host(`prospector.ie`) && PathPrefix(`/api/`)",
        "traefik.http.routers.prospector-api.entrypoints=websecure",
        "traefik.http.routers.prospector-api.tls.certresolver=lets-encrypt"
      ]

      canary_tags = [
        "traefik.enable=false"
      ]

      check {
        type = "http"
        path = "/api/health"
        interval = "5s"
        timeout = "1s"
      }
    }

    service {
      name = "prospector-api-canary"
      port = "http"

      tags = []
      canary_tags = [
        "traefik.enable=true",
        "traefik.http.routers.prospector-api.rule=Host(`canary.prospector.ie`) && PathPrefix(`/api/`)",
        "traefik.http.routers.prospector-api.entrypoints=websecure",
        "traefik.http.routers.prospector-api.tls.certresolver=lets-encrypt"
      ]

      check {
        type = "http"
        path = "/api/health"
        interval = "5s"
        timeout = "1s"
      }
    }

    task "server" {
      driver = "docker"

      config {
        image = "git.dbyte.xyz/distro/prospector/api:[[.git_sha]]"
        ports = ["http"]
      }

      resources {
        cpu = 128
        memory = 128
        
      }      
    }
  }
}
