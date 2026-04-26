import SchemaBuilder from "@pothos/core";

import type { ServerDeps } from "@library/server-bootstrap";

import { mapGraphqlError } from "./errors.js";

export interface GraphqlContext {
  deps: ServerDeps;
}

export const builder = new SchemaBuilder<{
  Context: GraphqlContext;
}>({});

export function safeResolver<TParent, TArgs, TResult>(
  resolver: (parent: TParent, args: TArgs, ctx: GraphqlContext) => Promise<TResult>
): (parent: TParent, args: TArgs, ctx: GraphqlContext) => Promise<TResult> {
  return async (parent, args, ctx) => {
    try {
      return await resolver(parent, args, ctx);
    } catch (error) {
      throw mapGraphqlError(error);
    }
  };
}
