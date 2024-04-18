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
							]
						},
						"Services": [
							{
								"Name": "{{ .Name }}",
								"Checks": [
									{
										"Name": "{{ .Name }}-health",
										"Type": "http",
										"Path": "/",
										"Interval": 10000000000,
										"Timeout": 2000000000
									}
								],
								{{ if .Network.Expose }}"Tags": [
									"traefik.enable=true",
									"traefik.http.routers.{{ .Name }}-{{ $.Name }}-{{ .UserConfig.User }}-prospector.rule=Host(` + "`" + `{{ .Name }}-{{ $.Name }}-{{ .UserConfig.User }}.prospector.ie` + "`" + `)",
									"traefik.http.routers.{{ .Name }}-{{ $.Name }}-{{ .UserConfig.User }}-prospector.entrypoints=websecure",
									"traefik.http.routers.{{ .Name }}-{{ $.Name }}-{{ .UserConfig.User }}-prospector.tls=true",
									"traefik.http.routers.{{ .Name }}-{{ $.Name }}-{{ .UserConfig.User }}-prospector.tls.certresolver=lets-encrypt"
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
			"job-type": "docker"
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
                                "type=1,serial=ds=nocloud-net;s=https://prospector.ie/api/vm-config/{{ .Name }}-vm/"
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
                                "RelativeDest": "local/name-vm.qcow2"
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
			"job-type": "vm"
		}
	}
}`
