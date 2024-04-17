package controllers

var DockerSource = `job "{{ .User }}-{{ .Name }}-prospector" {
	datacenters = ["dc1"]
	type = "service"
	
	meta {
		job-type = "docker"
	}
	{{ range .Components }}
	group "{{ .Name }}" {
		count = 1

		network {
			port "{{ .Name }}" {
				to = {{ .Network.Port }}
			}
		}
		
		task "{{ .Name }}" {
			driver = "docker"
			
			config {
				image = "{{ .Image }}"
				ports = ["{{ .Name }}"]
			}

			resources {
				cpu    = {{ .Resources.Cpu }}
				memory = {{ .Resources.Memory }}
			}

			service {
				name = "{{ .Name }}"
				port = "{{ .Name }}"

				check {
					name = "{{ .Name }}-health"
					type = "http"
					path = "/"
					interval = "10s"
					timeout = "2s"
				}
				{{ if .Network.Expose }}
				tags = [
					"traefik.enable=true",
					"traefik.http.routers.{{ .Name }}-{{ $.Name }}-{{ .UserConfig.User }}-prospector.rule=Host(` + "`" + `{{ .Name }}-{{ $.Name }}-{{ .UserConfig.User }}.prospector.ie` + "`" + `)",
					"traefik.http.routers.{{ .Name }}-{{ $.Name }}-{{ .UserConfig.User }}-prospector.entrypoints=websecure",
					"traefik.http.routers.{{ .Name }}-{{ $.Name }}-{{ .UserConfig.User }}-prospector.tls=true",
					"traefik.http.routers.{{ .Name }}-{{ $.Name }}-{{ .UserConfig.User }}-prospector.tls.certresolver=lets-encrypt"
				]
				{{ end }}
			}
		}
	}
	{{ end }}
}
`

var VMSource = `job "{{ .User }}-{{ .Name }}-prospector" {
  datacenters = ["dc1"]

  meta {
	job-type = "vm"
  }

  {{ range .Components }}
  group "{{ .User }}-{{ .Name }}" {

    network {
      mode = "host"
    }

    service {
      name = "{{ .Name }}-vm"
    }

    task "{{ .Name }}" {
      constraint {
        attribute = "${attr.unique.hostname}"
        value     = "hermes"
      }

      resources {
        cpu    = {{ .Resources.Cpu }}
        memory = {{ .Resources.Memory }}
      }

      artifact {
        source      = "https://cloud.debian.org/images/cloud/bookworm/latest/debian-12-genericcloud-amd64.qcow2"
		// source   = "{{ .Image }}"
        destination = "local/{{ .Name }}-vm.qcow2"
        mode        = "file"
      }

      driver = "qemu"

      config {
        image_path = "local/{{ .Name }}-vm.qcow2"
        accelerator = "kvm"
        drive_interface = "virtio"

        args = [
          "-netdev",
          "bridge,id=hn0",
          "-device",
          "virtio-net-pci,netdev=hn0,id=nic1,mac={{ .Mac }}",
          "-smbios",
          "type=1,serial=ds=nocloud-net;s=https://prospector.ie/api/vm-config/{{ .Name }}-vm/",
        ]
      }
    }
  }
  {{ end }}
}
`
