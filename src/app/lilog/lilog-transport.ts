import { once } from "events";
import fs from "fs";

type Options = {
  filename: string;
}

export default async (options: Options) => {
  const ws = fs.createWriteStream(`${process.cwd()}/logs/${options.filename}`);
  await once(ws, 'open');
  return ws
}
