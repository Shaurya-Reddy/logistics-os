import { readFileSync, writeFileSync } from 'node:fs';
import { execFileSync } from 'node:child_process';

const lock = JSON.parse(readFileSync('web/package-lock.json', 'utf8'));
const root = lock.packages[''];
const frontend = Object.entries(lock.packages).filter(([path]) => path).map(([path, pkg]) => ({
  path, version: pkg.version, development: Boolean(pkg.dev), optional: Boolean(pkg.optional),
  direct: Object.hasOwn(root.dependencies || {}, path.replace('node_modules/', '')) || Object.hasOwn(root.devDependencies || {}, path.replace('node_modules/', '')),
  license: pkg.license || 'review upstream metadata',
})).sort((a, b) => a.path.localeCompare(b.path));
const go = execFileSync('go', ['list', '-m', '-json', 'all'], { encoding: 'utf8' }).trim().split(/\n}\s*\n/).map((part, i, parts) => JSON.parse(i === parts.length - 1 ? part : part + '\n}')).filter(m => !m.Main).map(m => ({ path: m.Path, version: m.Version, direct: !m.Indirect })).sort((a,b) => a.path.localeCompare(b.path));
const inventory = { tools: JSON.parse(readFileSync('scripts/tools.json', 'utf8')), go, frontend,
  counts: { go: go.length, frontendRuntime: frontend.filter(p => !p.development).length, frontendDevelopment: frontend.filter(p => p.development).length } };
const text = JSON.stringify(inventory, null, 2) + '\n';
if (process.argv.includes('--write')) writeFileSync('docs/dependencies.json', text);
else if (readFileSync('docs/dependencies.json', 'utf8') !== text) throw Error('Dependency inventory drift: review added/removed modules and justification before regenerating');
console.log(JSON.stringify(inventory.counts));
