import cors from "cors";
import express from "express";
import { createYoga } from "graphql-yoga";

import type { ServerDeps } from "@library/server-bootstrap";

import { authMiddleware } from "./auth.js";
import { schema, type GraphqlContext } from "./graphql/schema.js";

export async function newRouter(deps: ServerDeps): Promise<express.Express> {
  const app = express();
  const yoga = createYoga<GraphqlContext>({
    graphqlEndpoint: "/graphql",
    schema,
    context: () => ({ deps })
  });

  app.get("/graphql", (_req, res) => {
    res.type("text/plain").send("GraphQL endpoint. Send authenticated POST requests to /graphql.");
  });

  app.use(
    "/graphql",
    cors(),
    authMiddleware(deps),
    yoga as express.RequestHandler
  );

  return app;
}
