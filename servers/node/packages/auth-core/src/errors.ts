export class InvalidTokenError extends Error {
  override readonly name = "InvalidTokenError";

  constructor(options?: ErrorOptions) {
    super("invalid token", options);
  }
}
