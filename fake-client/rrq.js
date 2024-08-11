import dgram from "node:dgram";

const server = "web api";
const buffer = new ArrayBuffer(1 + server.length + 1 + 8 + 8 + 1);
const view = new DataView(buffer);
const client = dgram.createSocket("udp4");

view.setUint8(0, 1);
let i;
for (i = 0; i < server.length; i++) {
  view.setUint8(1 + i, server[i].charCodeAt());
}

view.setUint8(1 + server.length, "\x00");

let from = 255 + 255;
let to = 255 + 100;
view.setBigUint64(1 + server.length + 1, BigInt(from), false);
view.setBigUint64(1 + server.length + 1 + 8, BigInt(to), false);

view.setUint8(1 + server.length + 1 + 8 + 8, "\x00");

console.log(buffer.slice(8));
console.log({ to }, { from });
console.log(view.getUint8(0));
console.log("to", view.getBigUint64(1 + server.length + 1 + 8 + 1));
console.log("from", view.getBigUint64(1 + server.length + 1));

client.connect(4119, "127.0.0.1", (err) => {
  if (err) {
    throw new Error(err);
  }
  client.send(view, (err) => {
    if (err) {
      throw new Error(err);
    }

    client.on("message", (msg, rinfo) => {
      console.log(`server got: ${msg} from ${rinfo.address}:${rinfo.port}`);

      client.close();
    });
  });
});
