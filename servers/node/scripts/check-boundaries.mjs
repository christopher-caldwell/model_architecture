import { readFileSync, readdirSync, statSync } from "node:fs";
import { join, relative } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const packageDir = join(root, "packages");

const allowedImports = new Map([
  ["domain", new Set()],
  ["application", new Set(["domain"])],
  ["auth-core", new Set()],
  ["persistence", new Set(["domain", "application"])],
  ["server-bootstrap", new Set(["application", "auth-core", "persistence"])],
  ["http-server", new Set(["server-bootstrap", "auth-core"])],
  ["graphql-server", new Set(["server-bootstrap", "auth-core"])]
]);

function walk(dir) {
  const files = [];
  for (const entry of readdirSync(dir)) {
    const path = join(dir, entry);
    const stat = statSync(path);
    if (stat.isDirectory()) {
      if (entry !== "dist" && entry !== "node_modules") files.push(...walk(path));
    } else if (path.endsWith(".ts") && !path.endsWith(".d.ts")) {
      files.push(path);
    }
  }
  return files;
}

const importPattern = /(?:import|export)\s+(?:type\s+)?(?:[^'"]*?\s+from\s+)?["'](@library\/([^/'"]+)(?:\/[^'"]*)?)["']/g;
const violations = [];

for (const packageName of allowedImports.keys()) {
  const srcDir = join(packageDir, packageName, "src");
  const allowed = allowedImports.get(packageName);

  for (const file of walk(srcDir)) {
    const source = readFileSync(file, "utf8");
    for (const match of source.matchAll(importPattern)) {
      const importedPackage = match[2];
      if (importedPackage === packageName) continue;
      if (!allowed.has(importedPackage)) {
        violations.push(
          `${relative(root, file)} imports ${match[1]}, but @library/${packageName} may only import ${formatAllowed(allowed)}`
        );
      }
    }
  }
}

if (violations.length > 0) {
  console.error("Import boundary violations:");
  for (const violation of violations) console.error(`- ${violation}`);
  process.exit(1);
}

console.log("Import boundaries ok");

function formatAllowed(allowed) {
  return allowed.size === 0 ? "no workspace packages" : [...allowed].map((name) => `@library/${name}`).join(", ");
}
