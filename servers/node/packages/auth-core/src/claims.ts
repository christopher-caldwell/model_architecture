export interface Claims {
  sub: string;
  exp: number;
  [claim: string]: unknown;
}
