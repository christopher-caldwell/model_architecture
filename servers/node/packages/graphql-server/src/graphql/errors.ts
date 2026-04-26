import { GraphQLError } from "graphql";

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

export function mapGraphqlError(error: unknown): GraphQLError {
  if (error instanceof GraphQLError) return error;

  if (error instanceof Error && notFoundErrors.has(error.name)) {
    return new GraphQLError(error.message, {
      extensions: { code: "NOT_FOUND" }
    });
  }

  if (error instanceof Error && conflictErrors.has(error.name)) {
    return new GraphQLError(error.message, {
      extensions: { code: "CONFLICT" }
    });
  }

  console.error("Unhandled GraphQL request error:", error);
  return new GraphQLError("Something went wrong", {
    extensions: { code: "INTERNAL_SERVER_ERROR" }
  });
}
