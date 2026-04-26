import type { Claims } from "./claims.js";

export interface AuthVerifier {
  verifyToken(token: string): Claims;
}
