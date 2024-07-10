import express from 'express'
import bodyParser from 'body-parser';
import  { pinoHttp } from 'pino-http';
import pino from 'pino';

const app = express();

app.use(bodyParser.json())
app.use(
  pinoHttp({
    logger: createLogger()
  })
)

app.get('/ping', async (_, res) => {
  return res.status(200).json('pong')
})

app.post('/ping', async (req, res) => {
  req.log.debug({ body: req.body })
  return res.status(200).json('pong')
})

function createLogger() {
  return pino();
}

export { app }
