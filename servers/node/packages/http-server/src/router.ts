import cors from "cors";
import express, { type Request, type RequestHandler, type Response, type Router } from "express";

import type { ServerDeps } from "@library/server-bootstrap";

import { authMiddleware } from "./auth.js";
import {
  bookCopyResponse,
  bookResponse,
  loanResponse,
  memberResponse,
  type CreateBookCopyRequestBody,
  type CreateBookRequestBody,
  type CreateLoanRequestBody,
  type CreateMemberRequestBody,
  type HealthCheckResponseBody
} from "./dto.js";
import { errorMiddleware } from "./errors.js";

export function newRouter(deps: ServerDeps): express.Express {
  const app = express();

  app.use(cors());
  app.use(express.json());

  app.get("/health", (_req: Request, res: Response<HealthCheckResponseBody>) => {
    res.json({ message: "ready" });
  });

  const protectedRouter = express.Router();
  protectedRouter.use(authMiddleware(deps));
  mountProtectedRoutes(protectedRouter, deps);

  app.use(protectedRouter);
  app.use(errorMiddleware);

  return app;
}

function mountProtectedRoutes(router: Router, deps: ServerDeps): void {
  router.get("/books", asyncHandler(async (_req, res) => {
    const books = await deps.catalog.queries.getBookCatalog();
    res.json(books.map(bookResponse));
  }));

  router.post("/books", asyncHandler(async (req: Request<unknown, unknown, CreateBookRequestBody>, res) => {
    const book = await deps.catalog.commands.addBook({
      isbn: req.body.isbn,
      title: req.body.title,
      author_name: req.body.author_name
    });
    res.status(201).json(bookResponse(book));
  }));

  router.post("/books/:isbn/copies", asyncHandler(async (req: Request<{ isbn: string }, unknown, CreateBookCopyRequestBody>, res) => {
    const copy = await deps.catalog.commands.addBookCopy({
      isbn: req.params.isbn,
      barcode: req.body.barcode
    });
    res.status(201).json(bookCopyResponse(copy));
  }));

  router.get("/book-copies/:barcode", asyncHandler(async (req: Request<{ barcode: string }>, res) => {
    const copy = await deps.catalog.queries.getBookCopyDetails(req.params.barcode);
    if (copy === null) {
      res.status(404).json({ error: "Book copy not found" });
      return;
    }
    res.json(bookCopyResponse(copy));
  }));

  router.put("/book-copies/:barcode/lost", asyncHandler(async (req: Request<{ barcode: string }>, res) => {
    const copy = await deps.catalog.commands.markBookCopyLost(req.params.barcode);
    res.json(bookCopyResponse(copy));
  }));

  router.delete("/book-copies/:barcode/lost", asyncHandler(async (req: Request<{ barcode: string }>, res) => {
    const copy = await deps.catalog.commands.markBookCopyFound(req.params.barcode);
    res.json(bookCopyResponse(copy));
  }));

  router.put("/book-copies/:barcode/maintenance", asyncHandler(async (req: Request<{ barcode: string }>, res) => {
    const copy = await deps.catalog.commands.sendBookCopyToMaintenance(req.params.barcode);
    res.json(bookCopyResponse(copy));
  }));

  router.delete("/book-copies/:barcode/maintenance", asyncHandler(async (req: Request<{ barcode: string }>, res) => {
    const copy = await deps.catalog.commands.completeBookCopyMaintenance(req.params.barcode);
    res.json(bookCopyResponse(copy));
  }));

  router.post("/book-copies/:barcode/return", asyncHandler(async (req: Request<{ barcode: string }>, res) => {
    const loan = await deps.lending.commands.returnBookCopy(req.params.barcode);
    res.json(loanResponse(loan));
  }));

  router.post("/book-copies/:barcode/report-loss", asyncHandler(async (req: Request<{ barcode: string }>, res) => {
    const copy = await deps.lending.commands.reportLostLoanedBookCopy(req.params.barcode);
    res.json(bookCopyResponse(copy));
  }));

  router.post("/members", asyncHandler(async (req: Request<unknown, unknown, CreateMemberRequestBody>, res) => {
    const member = await deps.membership.commands.registerMember({
      full_name: req.body.full_name,
      max_active_loans: req.body.max_active_loans
    });
    res.status(201).json(memberResponse(member));
  }));

  router.get("/members/:ident", asyncHandler(async (req: Request<{ ident: string }>, res) => {
    const member = await deps.membership.queries.getMemberDetails(req.params.ident);
    if (member === null) {
      res.status(404).json({ error: "Member not found" });
      return;
    }
    res.json(memberResponse(member));
  }));

  router.put("/members/:ident/suspension", asyncHandler(async (req: Request<{ ident: string }>, res) => {
    const member = await deps.membership.commands.suspendMember({ member_ident: req.params.ident });
    res.json(memberResponse(member));
  }));

  router.delete("/members/:ident/suspension", asyncHandler(async (req: Request<{ ident: string }>, res) => {
    const member = await deps.membership.commands.reactivateMember({ member_ident: req.params.ident });
    res.json(memberResponse(member));
  }));

  router.get("/members/:ident/loans", asyncHandler(async (req: Request<{ ident: string }>, res) => {
    const loans = await deps.lending.queries.getMemberLoans(req.params.ident);
    res.json(loans.map(loanResponse));
  }));

  router.post("/loans", asyncHandler(async (req: Request<unknown, unknown, CreateLoanRequestBody>, res) => {
    const loan = await deps.lending.commands.checkOutBookCopy({
      member_ident: req.body.member_ident,
      book_copy_barcode: req.body.book_copy_barcode
    });
    res.status(201).json(loanResponse(loan));
  }));

  router.get("/loans/overdue", asyncHandler(async (_req, res) => {
    const loans = await deps.lending.queries.getOverdueLoans();
    res.json(loans.map(loanResponse));
  }));
}

function asyncHandler(handler: (req: any, res: any) => Promise<void>): RequestHandler {
  return (req, res, next) => {
    handler(req, res).catch(next);
  };
}
