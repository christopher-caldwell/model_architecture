import type { GraphQLSchema } from "graphql";

import "./catalog/queries.js";
import "./catalog/mutations.js";
import "./lending/queries.js";
import "./lending/mutations.js";
import "./membership/queries.js";
import "./membership/mutations.js";
import { builder, type GraphqlContext } from "./schema-builder.js";

export type { GraphqlContext };

builder.queryType({});
builder.mutationType({});

export const schema: GraphQLSchema = builder.toSchema();
