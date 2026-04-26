import "dotenv/config";

import { createServerDeps, loadServerConfig } from "@library/server-bootstrap";

import { newRouter } from "./router.js";

const config = loadServerConfig();
const deps = await createServerDeps(config);
const app = newRouter(deps);

app.listen(config.serverPort, "0.0.0.0", () => {
  console.log(`HTTP server listening on :${config.serverPort}`);
});
