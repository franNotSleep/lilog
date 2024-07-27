export function getData(maxResponseTime) {
  const data = {
    level: getRandomLevel(),
    time: 1722088153891,
    pid: 108791,
    hostname: "frannotsleep-on-ubuntu",
    req: {
      id: 1,
      method: "GET",
      url: "/ping/warn/100",
      query: {},
      params: {},
      headers: {
        host: "localhost:3032",
        "user-agent": "curl/7.81.0",
        accept: "*/*",
      },
      remoteAddress: "::ffff:127.0.0.1",
      remotePort: 41760,
    },
    res: {
      statusCode: 200,
      headers: {
        "x-powered-by": "Express",
        "content-type": "application/json; charset=utf-8",
        "content-length": "6",
        etag: 'W/"6-uBnwlsJiQ3kuZAgKKWB4aV5ugdE"',
      },
    },
    responseTime: getRandomResponseTime(maxResponseTime),
    msg: "request completed",
  };

  return data;
}

function getRandomResponseTime(maxResponse) {
  return Math.floor(Math.random() * maxResponse);
}


function getRandomLevel() {
  const levels = [60, 50, 40, 30, 20, 10];
  const indx = Math.floor(Math.random() * levels.length);
  return levels[indx];
}
