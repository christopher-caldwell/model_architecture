import type { NextFunction, Request, Response } from "express";

import { InvalidTokenError, type Claims } from "@library/auth-core";
import type { ServerDeps } from "@library/server-bootstrap";

export interface AuthenticatedRequest extends Request {
  claims?: Claims;
}

export function authMiddleware(deps: ServerDeps) {
  return (req: AuthenticatedRequest, res: Response, next: NextFunction): void => {
    const authHeader = req.header("authorization");
    if (authHeader === undefined) {
      res.sendStatus(401);
      return;
    }

    const token = authHeader.startsWith("Bearer ") ? authHeader.slice("Bearer ".length) : null;
    if (token === null || token.length === 0) {
      res.sendStatus(401);
      return;
    }

    try {
      req.claims = deps.auth.verifier.verifyToken(token);
      next();
    } catch (error) {
      if (error instanceof InvalidTokenError) {
        res.sendStatus(403);
        return;
      }

      next(error);
    }
  };
}
