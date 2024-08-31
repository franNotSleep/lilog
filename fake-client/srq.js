import tls from "node:tls";
import fs from "fs";

import { getData } from "./data.js";

const options = {
  host: "localhost",
  ca: [fs.readFileSync("serverCert.pem")],
};

const socket = tls.connect(4119, options, () => {
  setInterval(() => {
    const maxResponseTime = Math.floor(Math.random() * 500) + 200;
    const data = getData(maxResponseTime);
    const buff = getBuffer(data);
    socket.write(buff);
  }, 2 * 1000);
});

socket.on("data", (data) => {
  console.log(data);
});

socket.on("end", () => {
  console.log("server ends connection");
});

function getBuffer(data) {
  const server = "web api";
  const invoice = JSON.stringify(data);
  const buffer = new ArrayBuffer(1 + server.length + 1 + invoice.length + 1);
  const view = new DataView(buffer);

  view.setUint8(0, 3);

  for (let i = 0; i < server.length; i++) {
    view.setUint8(1 + i, server[i].charCodeAt());
  }

  view.setUint8(1 + server.length, "\x00");

  for (let i = 0; i < invoice.length; i++) {
    view.setUint8(1 + 1 + server.length + i, invoice[i].charCodeAt());
  }

  view.setUint8(1 + 1 + server.length + invoice.length, "\x00");

  return view;
}
