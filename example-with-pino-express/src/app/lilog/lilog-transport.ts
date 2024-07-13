import { Writable } from "stream";
import dgram from "node:dgram";
//import { once } from "events";

type DGramStreamOptions = {
  port: number;
  highWaterMark: number;
};

class DGramStream extends Writable {
  private server: dgram.Socket;
  private chunks: any[];

  constructor({ highWaterMark }: DGramStreamOptions) {
    super({ highWaterMark });

    this.chunks = [];
    this.server = dgram.createSocket("udp4");
  }

  _write(
    chunk: any,
    _: BufferEncoding,
    callback: (error?: Error | null) => void,
  ): void {
    let start = 0;
    for (let i = 0; i < chunk.length; i++) {
      const byte = chunk[i];

      if (byte === 10) {
        this.server.send(chunk.slice(start, i), 4119, "localhost");
        start = i;
      }
    }

    callback();
  }

  _final(callback: (error?: Error | null) => void): void {
    this.server.send(Buffer.concat(this.chunks), 4119, "localhost", callback);
  }

  _destroy(
    error: Error | null,
    callback: (error?: Error | null) => void,
  ): void {
    callback(error);
  }
}

export default () => {
  const server = new DGramStream({ port: 4111, highWaterMark: 1024 });
  return server;
};
