let token =
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM1NDI0ODQsImlkIjoidGVzbGEiLCJvcmlnX2lhdCI6MTcxMzUzODg4NH0.iW2lpo-p3pkzygE6OJ6sAUDZqtq6AYlKja-5fLVkzF4";

let isStopped = false;
let type = "stdout";
let logs = "";
let abortController;

let stream;

// make a http request that will be upgraded to a websocket
// the url is at localhost:3434/api/v1/tesla-nginx-prospector/logs

const getLogs = async () => {
    const url = "http://localhost:3434/api/v1/jobs/tesla-nginx-prospector/logs";
    // add query params to the url
    // ?type=stdout and task=nginx
    const params = new URLSearchParams({ type, task: "nginx" });
    const fullUrl = `${url}?${params.toString()}`;
    const headers = {
        Authorization: `Bearer ${token}`,
        Connection: "Upgrade",
        Upgrade: "websocket",
    };

    const response = await fetch(fullUrl, {
        headers,
    });

    if (response.status === 200) {
        stream = await response.body.getReader();
        readStream(type);
    }

    return response;
};

const readStream = async () => {
    const { value, done } = await stream.read();
    if (done) {
        console.log("Stream is done");
        return;
    }

    console.log("Reading stdout");
    const text = new TextDecoder().decode(value);
    console.log(text);
    logs += text;
    document.getElementById("output").innerText = logs;
    scrollToBottom(document.getElementById("output"));

    readStream();
};

const stopStream = async () => {
    isStopped = true;
    stream.cancel();
};

// get logs for both stdout and stderr

function scrollToBottom(element) {
    element.scroll({ top: element.scrollHeight, behavior: "smooth" });
}

getLogs();
