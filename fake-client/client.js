import dgram from "node:dgram";
import { getData } from './data.js'

const client = dgram.createSocket("udp4");

const [,,n, maxResponseTime] = process.argv;

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
    const buff = Buffer.from(JSON.stringify(data))
    client.send(Buffer.from(buff), (err) => {
      if (err) {
        throw new Error(err);
      }

      if (i === n - 1) {
        client.close();
      }

    });
  }
});
