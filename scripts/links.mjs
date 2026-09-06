import { readFileSync, existsSync } from 'node:fs';
import { dirname, resolve } from 'node:path';
for (const path of ['AGENTS.md', 'ARCHITECTURE.md', 'PRODUCT.md', 'PERFORMANCE_BUDGET.md', 'README.md', 'docs/ERP_REFERENCE.md', 'docs/DEPENDENCIES.md']) {
  for (const match of readFileSync(path, 'utf8').matchAll(/\]\(([^)]+)\)/g)) {
    if (/^(https?:|#)/.test(match[1])) continue;
    if (!existsSync(resolve(dirname(path), match[1].split('#')[0]))) throw Error(`${path}: broken link ${match[1]}`);
  }
}
console.log('Foundation links valid');
