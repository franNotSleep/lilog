import dgram from "node:dgram";
import { getData } from "./data.js";

const client = dgram.createSocket("udp4");

const [, , n, maxResponseTime] = process.argv;

if (!n || !maxResponseTime) {
  console.log(`Usage: node ${process.argv[1]} <n request> <max response time>`);
  process.exit(1);
}

client.connect(4119, "127.0.0.1", (err) => {
  if (err) {
    throw new Error(err);
  }
  for (let i = 0; i < +n; i++) {
    const data = getData(+maxResponseTime);
    const buff = getBuffer(data);
    client.send(buff, (err) => {
      if (err) {
        throw new Error(err);
      }

      if (i === n - 1) {
        client.close();
      }
    });
  }
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
