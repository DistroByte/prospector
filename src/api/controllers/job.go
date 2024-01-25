package controllers

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
)

type Job struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	Port  struct {
		Label     string `json:"label"`
		Container string `json:"container"`
	} `json:"port"`

	Tags []string `json:"tags"`
}

// @Summary		Create a job
// @Description	Create and submit a job to nomad to deploy
// @Tags			job
// @Accept			json
// @Produce		json
// @Security		None
// @Param			job	body		Job	true	"Job"
// @Success		200	{object}	Message
// @Router			/create [post]
func (c *Controller) CreateJob(ctx *gin.Context) {
	var job Job

	if err := ctx.BindJSON(&job); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := createNomadJob(job)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	println(res)

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// post to nomad server for job creation
func createNomadJob(job Job) (int, error) {
	jobSource := `job "{{ .Name }}-prospector" {
	datacenters = ["dc1"]
	type = "service"

	group "{{ .Name }}" {
		count = 1

		network {
			port "{{ .Port.Label }}" {
				to = {{ .Port.Container }}
			}
		}

		service {
			name = "{{ .Name }}"
			port = "{{ .Port.Label }}"

			check {
				name = "{{ .Name }}-health"
				type = "http"
				path = "/"
				interval = "10s"
				timeout = "2s"
			}

			tags = [
				"{{ .Tags }}"
			]
		}

		task "{{ .Name }}" {
			driver = "docker"
			
			config {
				image = "{{ .Image }}"
			}
		}
	}
}
`

	jobSource = strings.ReplaceAll(jobSource, `"`, `\"`)

	jobTemplate, err := template.New("job").Parse(jobSource)
	if err != nil {
		return 500, err
	}

	var jobSpec bytes.Buffer

	if err := jobTemplate.Execute(&jobSpec, job); err != nil {
		return 500, err
	}

	parseBodySource := `{ "JobHCL": "{{ . }}", "Canonicalize": true }`

	parseBodyTemplate, err := template.New("parseBody").Parse(parseBodySource)
	if err != nil {
		return 500, err
	}

	var parseBody bytes.Buffer

	if err := parseBodyTemplate.Execute(&parseBody, jobSpec.String()); err != nil {
		return 500, err
	}

	// println(parseBody.String())

	parseBodyCleaned := strings.ReplaceAll(parseBody.String(), "\n", `\n`)
	parseBodyCleaned = strings.ReplaceAll(parseBodyCleaned, "\t", ` `)

	println(parseBodyCleaned)

	parseResponse, err := http.Post("http://zeus.internal:4646/v1/jobs/parse", "application/json", strings.NewReader(parseBodyCleaned))
	if err != nil {
		return 500, err
	}

	// parse response
	var parseResponseBuffer bytes.Buffer

	if _, err := io.Copy(&parseResponseBuffer, parseResponse.Body); err != nil {
		return 500, err
	}

	println(parseResponseBuffer.String())

	jobRunSource := `{ "Job": {{ . }} }`

	jobRunTemplate, err := template.New("jobRun").Parse(jobRunSource)
	if err != nil {
		return 500, err
	}

	var jobRun bytes.Buffer

	if err := jobRunTemplate.Execute(&jobRun, parseResponseBuffer.String()); err != nil {
		return 500, err
	}

	// println(jobRun.String())

	jobRunResponse, err := http.Post("http://zeus.internal:4646/v1/jobs", "application/json", &jobRun)
	if err != nil {
		return 500, err
	}

	// parse response
	var jobRunResponseBuffer bytes.Buffer

	if _, err := io.Copy(&jobRunResponseBuffer, jobRunResponse.Body); err != nil {
		return 500, err
	}

	println(jobRunResponseBuffer.String())

	return 200, nil
}
