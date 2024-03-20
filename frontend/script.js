let started = false;
let intervalId = null;

const startButton = document.getElementById('btn-start');
const stopButton = document.getElementById('btn-stop');
const textArea = document.getElementById('text-area');

async function doQuery() {
  const backendURL = document.getElementById('backend-url').value;

  const res = await fetch(backendURL, {
    cache: 'no-cache',
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
  });

  const data = await res.json();

  if (!res.ok) {
    return 'query error: ' + data.message;
  }

  return data.info;
}

function start() {
  if (!started) {
    started = true;
    startButton.disabled = true;
    stopButton.disabled = false;

    textArea.value = textArea.value + '***\nStarted querying...\n';
    textArea.scrollTop = textArea.scrollHeight;

    intervalId = setInterval(async () => {
      info = await doQuery();
      textArea.value = `${textArea.value} + Response from host: '${info.hostname}' - ${info.db_string}.\n`;
      textArea.scrollTop = textArea.scrollHeight;
    }, 250);
  }
}

function stop() {
  if (started) {
    started = false;
    startButton.disabled = false;
    stopButton.disabled = true;

    clearInterval(intervalId);

    textArea.value = textArea.value + '***\nStopped querying.\n';
    textArea.scrollTop = textArea.scrollHeight;
  }
}
