export class UnexpectedCommandError extends Error {
  override readonly name = "UnexpectedCommandError";

  constructor(message: string, options?: ErrorOptions) {
    super(message, options);
  }
}

export function isUnexpectedCommandError(error: unknown): error is UnexpectedCommandError {
  return error instanceof UnexpectedCommandError;
}

export function toUnexpectedCommandError(message: string, error: unknown): UnexpectedCommandError {
  if (error instanceof UnexpectedCommandError) return error;
  if (error instanceof Error) return new UnexpectedCommandError(message, { cause: error });
  return new UnexpectedCommandError(message, { cause: new Error(String(error)) });
}
