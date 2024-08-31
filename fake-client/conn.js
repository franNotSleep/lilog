import tls from "node:tls";
import fs from "fs";

const options = {
  host: "localhost",
  ca: [fs.readFileSync("serverCert.pem")],
};

const socket = tls.connect(4119, options, () => {
  socket.write(reqView());
});

socket.setEncoding("utf8");

socket.on("data", (data) => {
  console.log(JSON.parse(data));
});

socket.on("end", () => {
  console.log("server ends connection");
});

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
