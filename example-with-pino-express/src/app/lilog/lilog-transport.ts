import { once } from "events";
import { Writable } from "stream";
import dgram from "node:dgram";

type Options = {
  filename: string;
};

class DGramStream extends Writable {
  private server: dgram.Socket;
  constructor(
    type: dgram.SocketType,
    callback?: (msg: Buffer, rinfo: dgram.RemoteInfo) => void,
  ) {
    super();
    this.server = dgram.createSocket(type, callback);
  }

  async waitToConnect() {
    await once(this.server, "connect");
  }

  _write(
    chunk: any,
    _: BufferEncoding,
    callback: (error?: Error | null) => void,
  ): void {
    process.stdout.write(chunk);
    callback
  }
}

export default async (_: Options) => {
  const server = new DGramStream("udp4");
  await server.waitToConnect();
  return server;
};
