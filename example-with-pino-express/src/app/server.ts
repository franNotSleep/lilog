import express from "express";
import bodyParser from "body-parser";
import { pinoHttp } from "pino-http";
import pino from "pino";

const app = express();

const levels = {
  emerg: 80,
  alert: 70,
  crit: 60,
  error: 50,
  warn: 40,
  notice: 30,
  info: 20,
  debug: 10,
};

app.use(bodyParser.json());
app.use(
  pinoHttp({
    logger: createLogger(),
  }),
);

app.get("/ping/:level/:ms", async (req, res) => {
  let level: pino.Level = req.params.level as any;

  req.log[level]("error");

  setTimeout(() => {
    return res.status(200).json("pong");
  }, Number(req.params.ms));
});

app.post("/ping", async (_, res) => {
  return res.status(200).json("pong");
});

function createLogger() {
  const transport = pino.transport({
    targets: [
      {
        target: "../../build/app/lilog/lilog-transport",
      },
    ],
  });

  const logger = pino(transport);
  logger.customLevels = levels;
  logger.useOnlyCustomLevels = true;

  return logger;
}

export { app };
