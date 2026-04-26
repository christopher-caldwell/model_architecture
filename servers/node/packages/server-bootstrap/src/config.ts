export interface ServerConfig {
  databaseRoUrl: string;
  databaseRwUrl: string;
  jwtSecret: string;
  serverPort: number;
}

export function loadServerConfig(env: NodeJS.ProcessEnv = process.env): ServerConfig {
  return {
    databaseRoUrl: requiredEnv(env, "DATABASE_RO_URL"),
    databaseRwUrl: requiredEnv(env, "DATABASE_RW_URL"),
    jwtSecret: requiredEnv(env, "JWT_SECRET"),
    serverPort: parsePort(env.SERVER_PORT ?? "3000")
  };
}

function requiredEnv(env: NodeJS.ProcessEnv, name: string): string {
  const value = env[name];
  if (value === undefined || value.length === 0) {
    throw new Error(`${name} must be set`);
  }

  return value;
}

function parsePort(value: string): number {
  const port = Number(value);
  if (!Number.isInteger(port) || port < 0 || port > 65535) {
    throw new Error("SERVER_PORT must be a valid port number (0-65535)");
  }

  return port;
}
