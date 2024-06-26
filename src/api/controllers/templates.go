package controllers

var DockerSourceJson = `{
	"Job": {
		"ID": "{{ .User }}-{{ .Name }}-prospector",
		"Name": "{{ .User }}-{{ .Name }}-prospector",
		"Type": "service",
		"Datacenters": [
			"dc1"
		],
		"TaskGroups": [
			{{ range $i, $_ := .Components }}{
				{{ $component := . }}
				"Name": "{{ .Name }}",
				"Count": 1,
				"Tasks": [
					{
						"Name": "{{ .Name }}",
						"Driver": "docker",
						"Config": {
							"image": "{{ .Image }}",
							"ports": [
								"{{ .Name }}"
							],
							"volumes": [
								"/data/prospector/{{ .UserConfig.User }}:/mnt/user-storage"{{ if .Volumes }},{{ end }}
								{{ if .Volumes }}{{ range $i, $v := .Volumes }}
								"/data/prospector/{{ $component.UserConfig.User }}/{{ $component.Name }}/{{ $v }}:/mnt/component-storage/{{ $v }}"{{ if not (last $i $component.Volumes) }},{{ end }}
								{{ end }}{{ end }}
							]
						},
						"Services": [
							{
								"Name": "{{ .Name }}",
								{{ if .Network.Expose }}"Tags": [
									"traefik.enable=true",
									"traefik.http.routers.{{ .UserConfig.User }}-{{ $.Name }}-{{ .Name }}-prospector.rule=Host(` + "`" + `{{ .Name }}-{{ $.Name }}-{{ .UserConfig.User }}.prospector.ie` + "`" + `)",
									"traefik.http.routers.{{ .UserConfig.User }}-{{ $.Name }}-{{ .Name }}-prospector.entrypoints=websecure",
									"traefik.http.routers.{{ .UserConfig.User }}-{{ $.Name }}-{{ .Name }}-prospector.tls=true",
									"traefik.http.routers.{{ .UserConfig.User }}-{{ $.Name }}-{{ .Name }}-prospector.tls.certresolver=lets-encrypt",
									"prometheus.io/scrape=false"
								],{{ else }}"Tags": [
									"prometheus.io/scrape=false"
								],{{ end }}
								"PortLabel": "{{ .Name }}"
							}
						],
						"Resources": {
							"CPU": {{ .Resources.Cpu }},
							"MemoryMB": {{ .Resources.Memory }}
						}
					}
				],
				"Networks": [
					{
						"DynamicPorts": [
							{
								"Label": "{{ .Name }}",
								"Value": 0,
								"To": {{ .Network.Port }}
							}
						]
					}
				]
			}{{ if not (last $i $.Components) }},{{ end }}{{ end }}	
		],
		"Meta": {
			"job-type": "docker",
			"job-definition": "{{ json | escapeQuotes }}"
		}
	}
}`

var VMSourceJson = `{
	"Job": {
		"ID": "{{ .User }}-{{ .Name }}-prospector",
		"Name": "{{ .User }}-{{ .Name }}-prospector",
		"Type": "service",
		"Datacenters": [
			"dc1"
		],
		"TaskGroups": [
			{{ range $i, $_ := .Components }}{
				"Name": "{{ .Name }}",
				"Count": 1,
				"Tasks": [
					{
						"Name": "{{ .Name }}",
						"Driver": "qemu",
						"Config": {
                            "accelerator": "kvm",
                            "args": [
                                "-netdev",
                                "bridge,id=hn0",
                                "-device",
                                "virtio-net-pci,netdev=hn0,id=nic1,mac={{ .Mac }}",
                                "-smbios",
                                "type=1,serial=ds=nocloud-net;s=https://prospector.ie/api/vm-config/{{ $.User }}-{{ .Name }}-vm/"
                            ],
                            "drive_interface": "virtio",
                            "image_path": "local/{{ .Name }}-vm.qcow2"
                        },
						"Constraints": [
                            {
                                "LTarget": "${attr.unique.hostname}",
                                "RTarget": "hermes",
                                "Operand": "="
                            }
                        ],
						"Resources": {
							"CPU": {{ .Resources.Cpu }},
							"MemoryMB": {{ .Resources.Memory }}
						},
						"Artifacts": [
                            {
                                "GetterSource": "https://cloud.debian.org/images/cloud/bookworm/latest/debian-12-genericcloud-amd64.qcow2",
                                "GetterMode": "file",
                                "RelativeDest": "local/{{ .Name }}-vm.qcow2"
                            }
                        ]
					}
				],
				"Services": [
                    {
                        "Name": "{{ .Name }}-vm"
                    }
                ]
			}{{ if not (last $i $.Components) }},{{ end }}{{ end }}
		],
		"Meta": {
			"job-type": "vm",
			"job-definition": "{{ json | escapeQuotes }}"

		}
	}
}`
