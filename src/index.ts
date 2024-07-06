import { app } from "./app/server";

const PORT = 3032

app.listen(PORT, () => {
  console.log(`:${PORT}`)
})
