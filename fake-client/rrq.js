import dgram from "node:dgram";

const server = dgram.createSocket("udp4");

server.on("listening", () => {
  server.send(reqView(), 4119, "127.0.0.1", (err) => {
    if (err) {
      throw new Error(err);
    }
  });
});

let i = 1;

server.on("message", (msg, rinfo) => {
  let data;
  try {
    data = JSON.parse(msg);
  } catch (error) {
    data = msg;
  }
  console.log(rinfo, data);

  if (i <= 0) {
    setTimeout(() => {
      server.send(ackView(), rinfo.port, rinfo.address, (err) => {
        if (err) {
          throw new Error(err);
        }
      });
    }, 1000);
  }

  i--;
});

function ackView() {
  const buff = new ArrayBuffer(2);
  const view = new DataView(buff);
  view.setUint16(0, 3);
  return view;
}

function reqView() {
  const server = "web api";
  const buffer = new ArrayBuffer(1 + server.length + 1 + 8 + 8 + 1);
  const view = new DataView(buffer);

  view.setUint8(0, 1);
  for (let i = 0; i < server.length; i++) {
    view.setUint8(1 + i, server[i].charCodeAt());
  }

  view.setUint8(1 + server.length, "\x00");

  let from = 255 + 255;
  let to = 255 + 100;
  view.setBigUint64(1 + server.length + 1, BigInt(from), false);
  view.setBigUint64(1 + server.length + 1 + 8, BigInt(to), false);
  view.setUint8(1 + server.length + 1 + 8 + 8, 1);

  return view;
}

process.on("SIGINT", () => {
  console.log("\nclosing client... 👀");
  server.close();
  console.log("closed. bye bye :) 👊🏾");
});

server.bind(55231);
