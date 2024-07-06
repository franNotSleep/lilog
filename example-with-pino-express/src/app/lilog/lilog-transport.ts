import { once } from "events";
//import fs from "fs";
import net from 'net';

type Options = {
  filename: string;
}


export default async (_: Options) => {
  const connOptions: net.NetConnectOpts = {
    host: '127.0.0.1',
    port: 41169,
  }

  const conn = net.connect(connOptions)

  await once(conn, 'connect')

  conn.on('close', () => {
    console.log("closed connection");
    conn.connect(connOptions);
  })

  conn.on('data', () => {
    console.log("new data received");
  })

//  const ws = fs.createWriteStream(`${process.cwd()}/logs/${options.filename}`);
//  await once(ws, 'open');
  return conn
}

