import { readFileSync, readdirSync, statSync, mkdirSync, writeFileSync } from 'node:fs';
import { gzipSync } from 'node:zlib';
import { execFileSync } from 'node:child_process';
import os from 'node:os';

const manifest = JSON.parse(readFileSync('web/dist/.vite/manifest.json', 'utf8'));
const initial = new Set();
const visited = new Set();
function visit(key) {
  if (visited.has(key)) return;
  visited.add(key);
  const chunk = manifest[key];
  if (!chunk) throw Error(`Missing manifest chunk ${key}`);
  initial.add(chunk.file);
  for (const css of chunk.css || []) initial.add(css);
  for (const child of chunk.imports || []) visit(child);
}
visit('index.html');
function walk(path) {
  return readdirSync(path).flatMap(name => {
    const child = `${path}/${name}`;
    return statSync(child).isDirectory() ? walk(child) : [child.replace('web/dist/', '')];
  });
}
const all = walk('web/dist');
if (all.some(path => path.endsWith('.map'))) throw Error('Public source map');
const html = readFileSync('web/dist/index.html', 'utf8');
if (/<style\b|<script\b(?![^>]*\bsrc=)|(?:src|href)=["']https?:/i.test(html)) throw Error('Inline or remote code needs explicit accounting');
const bytes = files => Object.fromEntries(['js', 'css'].map(ext => [ext, [...files].filter(path => path.endsWith(`.${ext}`)).reduce((sum, path) => sum + gzipSync(readFileSync(`web/dist/${path}`), { level: 9 }).length, 0)]));
const sizes = { initial: bytes(initial), total: bytes(all) };
for (const [kind, ext, max] of [['initial', 'js', 200000], ['initial', 'css', 50000], ['total', 'js', 600000], ['total', 'css', 100000]]) {
  if (sizes[kind][ext] > max) throw Error(`${kind} ${ext} ${sizes[kind][ext]} exceeds ${max}`);
}
let image = null;
let services = null;
let readiness = null;
if (process.argv.includes('--image')) {
  const [details] = JSON.parse(execFileSync('docker', ['image', 'inspect', 'logistics-os:local'], { encoding: 'utf8' }));
  if (details.Os !== 'linux' || details.Architecture !== 'amd64') throw Error('Image measurement requires linux/amd64');
  if (!Number.isFinite(details.Size) || details.Size <= 0) throw Error('Missing runtime image size');
  if (details.Size > 100000000) throw Error(`Runtime image exceeds 100 MB: ${details.Size}`);
  image = { bytes: details.Size, id: details.Id, platform: 'linux/amd64', compressedTransferBytes: null };
  const config = JSON.parse(execFileSync('docker', ['compose', 'config', '--format', 'json'], { encoding: 'utf8' }));
  services = Object.keys(config.services).sort();
  if (JSON.stringify(services) !== '["app","postgres"]') throw Error('Exactly app + postgres required');
  for (const service of Object.values(config.services)) if (service.profiles) throw Error('Default service may not be optional');
  execFileSync('docker', ['run', '--rm', '--entrypoint', '/bin/sh', 'logistics-os:local', '-ec', 'test "$(id -u)" != 0; for x in node npm go gcc; do if command -v "$x"; then exit 1; fi; done; test ! -d /src; test -f /etc/ssl/certs/ca-certificates.crt; test -d /usr/share/zoneinfo; test -z "$(find /var/cache/apk -type f)"']);
  const samplesMS = [];
  for (let i = 0; i < 1000; i++) {
    const start = performance.now();
    const response = await fetch('http://127.0.0.1:8080/ready', { signal: AbortSignal.timeout(2000) });
    const body = await response.json();
    if (!response.ok || body.status !== 'ready') throw Error('Readiness failed during measurement');
    samplesMS.push(performance.now() - start);
  }
  const sorted = [...samplesMS].sort((a, b) => a - b);
  readiness = { samplesMS, p50MS: sorted[499], p95MS: sorted[949], p99MS: sorted[989],
    conditions: '1000 sequential local-network probes, shared CI runner, no warmup; diagnostic only, not a reference deployment latency baseline' };
}
const report = {
  commit: process.env.GITHUB_SHA || execFileSync('git', ['rev-parse', 'HEAD'], { encoding: 'utf8' }).trim(),
  node: process.version, go: execFileSync('go', ['version'], { encoding: 'utf8' }).trim(),
  runner: { platform: os.platform(), arch: os.arch(), cpus: os.cpus().length, memoryBytes: os.totalmem() },
  dataset: 'S0 schema history only; no business fixture',
  commands: ['npm --prefix web ci', 'npm --prefix web run build', `node scripts/budgets.mjs${image ? ' --image' : ''}`],
  gzipLevel: 9, frontend: sizes, initialAssets: [...initial].sort(), image, services, readiness,
  controlledStartupAndRSS: { status: 'pending controlled runner', startupP95Seconds: null, idleMaxRSSBytes: null },
  domainBenchmarks: 'not applicable in S0',
};
mkdirSync('reports', { recursive: true });
writeFileSync('reports/budgets.json', JSON.stringify(report, null, 2) + '\n');
console.log(JSON.stringify(report, null, 2));
