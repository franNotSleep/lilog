import fs from "node:fs";
import tls, { TLSSocket } from "node:tls";
import { Writable } from "node:stream";
import { once } from "node:events";

const tlsOptions = {
  host: "localhost",
  ca: [fs.readFileSync("serverCert.pem")],
};

export default async () => {
  const t = new Transport();

  return t;
};

class Transport extends Writable {
  /**
   * @type TLSSocket
   * @private
   */
  tlsSocket;

  constructor() {
    super();
  }

  async _construct(cb) {
    try {
      this.tlsSocket = tls.connect(4119, tlsOptions);
      await once(this.tlsSocket, "secure");
      cb();
    } catch (error) {
      cb(error);
    }
  }

  _write(chunk, encoding, callback) {
    const logs = chunk.toString().split("\n");

    for (let i = 0; i < logs.length; i++) {
      const log = logs[i];

      if (!log) continue;

      const payload = this.prepareData(log);
      this.tlsSocket.write(payload);
    }

    callback();
  }

  _final(callback) {
    const data = chunk.toString();
    const payload = this.prepareData(data);

    this.tlsSocket.write(payload);

    callback();
  }

  _destroy(error, callback) {
    this.tlsSocket.destroy();
    callback(error);
  }

  prepareData(data) {
    const server = "web api";
    const invoice = data;
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
}
