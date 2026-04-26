import type { ErrorRequestHandler, Response } from "express";

import type { ErrorResponseBody } from "./dto.js";

const notFoundErrors = new Set([
  "BookNotFoundError",
  "BookCopyNotFoundError",
  "MemberNotFoundError"
]);

const conflictErrors = new Set([
  "BookCopyCannotBeBorrowedError",
  "BookCopyCannotBeSentToMaintenanceError",
  "BookCopyCannotBeReturnedFromMaintenanceError",
  "BookCopyCannotMarkLostError",
  "BookCopyCannotBeReturnedFromLostError",
  "MemberCannotBeSuspendedError",
  "MemberCannotBeReactivatedError",
  "MemberCannotBorrowWhileSuspendedError",
  "MemberLoanLimitReachedError",
  "NoActiveLoanForBookCopyError",
  "LoanCannotBeReturnedError"
]);

export const errorMiddleware: ErrorRequestHandler = (error, _req, res, _next) => {
  sendError(res, error);
};

export function sendError(res: Response, error: unknown): void {
  if (isNamedError(error) && notFoundErrors.has(error.name)) {
    res.status(404).json(errorBody(error.message));
    return;
  }

  if (isNamedError(error) && conflictErrors.has(error.name)) {
    res.status(409).json(errorBody(error.message));
    return;
  }

  console.error("Unhandled HTTP request error:", error);
  res.status(500).json(errorBody("Something went wrong"));
}

function errorBody(error: string): ErrorResponseBody {
  return { error };
}

function isNamedError(error: unknown): error is Error {
  return error instanceof Error && typeof error.name === "string";
}
