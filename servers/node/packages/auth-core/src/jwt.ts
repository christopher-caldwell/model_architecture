import jwt from "jsonwebtoken";

import type { Claims } from "./claims.js";
import { InvalidTokenError } from "./errors.js";
import type { AuthVerifier } from "./verifier.js";

export const jwtAudience = "ops.craftcode.solutions";

export class JwtAuthAdapter implements AuthVerifier {
  constructor(private readonly jwtSecret: string) {}

  verifyToken(token: string): Claims {
    try {
      const decoded = jwt.verify(token, this.jwtSecret, {
        algorithms: ["HS256"],
        audience: jwtAudience
      });

      if (typeof decoded === "string") {
        throw new InvalidTokenError();
      }

      return decoded as Claims;
    } catch (error) {
      if (error instanceof InvalidTokenError) throw error;
      throw new InvalidTokenError({ cause: error });
    }
  }
}
