import {
  CatalogCommands,
  CatalogQueries,
  LendingCommands,
  LendingQueries,
  MembershipCommands,
  MembershipQueries,
  type IdentGenerator
} from "@library/application";
import { JwtAuthAdapter, type AuthVerifier } from "@library/auth-core";
import {
  BookCopyReadRepositoryPostgres,
  BookReadRepositoryPostgres,
  LoanReadRepositoryPostgres,
  MemberReadRepositoryPostgres,
  SqlWriteUnitOfWorkFactory
} from "@library/persistence";
import { nanoid } from "nanoid";
import { Pool } from "pg";

import type { ServerConfig } from "./config.js";

export interface ServerDeps {
  auth: AuthDeps;
  catalog: CatalogDeps;
  lending: LendingDeps;
  membership: MembershipDeps;
}

export interface AuthDeps {
  verifier: AuthVerifier;
}

export interface CatalogDeps {
  commands: CatalogCommands;
  queries: CatalogQueries;
}

export interface LendingDeps {
  commands: LendingCommands;
  queries: LendingQueries;
}

export interface MembershipDeps {
  commands: MembershipCommands;
  queries: MembershipQueries;
}

export async function createServerDeps(config: ServerConfig): Promise<ServerDeps> {
  const roPool = createPool(config.databaseRoUrl);
  const rwPool = createPool(config.databaseRwUrl);

  const writeUnitOfWorkFactory = new SqlWriteUnitOfWorkFactory(rwPool);

  const bookReadRepository = new BookReadRepositoryPostgres(roPool);
  const bookCopyReadRepository = new BookCopyReadRepositoryPostgres(roPool);
  const loanReadRepository = new LoanReadRepositoryPostgres(roPool);
  const memberReadRepository = new MemberReadRepositoryPostgres(roPool);
  const identGenerator = new RandomIdentGenerator();

  return {
    auth: {
      verifier: new JwtAuthAdapter(config.jwtSecret)
    },
    catalog: {
      commands: new CatalogCommands(writeUnitOfWorkFactory),
      queries: new CatalogQueries(bookReadRepository, bookCopyReadRepository)
    },
    lending: {
      commands: new LendingCommands(writeUnitOfWorkFactory),
      queries: new LendingQueries(loanReadRepository)
    },
    membership: {
      commands: new MembershipCommands(writeUnitOfWorkFactory, identGenerator),
      queries: new MembershipQueries(memberReadRepository)
    }
  };
}

function createPool(connectionString: string): Pool {
  return new Pool({
    connectionString,
    max: 5
  });
}

class RandomIdentGenerator implements IdentGenerator {
  generate(): string {
    return nanoid(10);
  }
}
