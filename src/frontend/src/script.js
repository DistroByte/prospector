document.getElementById("fetch-jobs").addEventListener("click", () => {
    fetchJobs();
});

function fetchJobs() {
    fetch("/api/jobs")
        .then((response) => {
            if (response.status == 200) {
                return response.json();
            }

            return Promise.resolve({ message: "No jobs found" });
        })
        .then((data) => {
            if (data.message == "No jobs found") {
                var jobData = document.getElementById("job-data");

                return (jobData.innerHTML = `<p>No Jobs Found</p>`);
            }

            var jobData = document.getElementById("job-data");

            jobData.innerHTML = data
                .map((job) => {
                    return `
                <div>
                  <h2>${job.id}</h2>
                  <p>Status: ${job.status}</p>
                  <p>Meta: ${job.Meta || "none"}</p>
                  <p>URL: <a href="https://${
                      job.id.split("-")[0]
                  }.prospector.ie">View app</a></p>
    
                  <button onclick="moreInfo('${job.id}')">More Info</button>
    
                  <button onclick="deleteJob('${job.id}', false)">Stop</button>
                  <button onclick="deleteJob('${job.id}', true)">Delete</button>
    
                  <button onclick="hideInfo()">Hide Info</button>
    
                  <div id="job-info" class="job-info"></div>
    
                </div>
              `;
                })
                .join("");
        });
}

function moreInfo(id) {
    fetch(`/api/jobs/${id}`)
        .then((response) => response.json())
        .then((data) => {
            var moreInfo = document.getElementById("job-info");
            moreInfo.innerHTML = `
              <pre>${JSON.stringify(data, null, 2)}</pre>
              `;
        });
}

function deleteJob(id, purge) {
    hideInfo();

    fetch(`/api/jobs/${id}?purge=${purge}`, {
        method: "DELETE",
    })
        .then((response) => response.json())
        .then((data) => {
            setTimeout(() => {
                fetchJobs();
            }, 3000);
        });
}

function hideInfo() {
    var moreInfo = document.getElementById("job-info");
    moreInfo.innerHTML = "";
}

document
    .getElementById("create-job-form")
    .addEventListener("submit", (event) => {
        event.preventDefault();

        var jobName = document.getElementById("job-name").value;
        var jobImage = document.getElementById("job-image").value;
        var jobPort = document.getElementById("job-port").value;
        var jobCpu = document.getElementById("job-cpu").value;
        var jobMemory = document.getElementById("job-memory").value;

        createJob(jobName, jobImage, jobPort, jobCpu, jobMemory);
    });

function createJob(name, image, port, cpu, memory) {
    fetch("/api/jobs", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            name: name,
            image: image,
            port: parseInt(port),
            cpu: parseInt(cpu),
            memory: parseInt(memory),
        }),
    })
        .then((response) => response.json())
        .then(() => {
            setTimeout(() => {
                fetchJobs();
            }, 3000);
        });
}
