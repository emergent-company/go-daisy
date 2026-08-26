import { execSync, spawn } from 'child_process';
import * as fs from 'fs';
import * as path from 'path';

const PROJECT_ROOT = path.resolve(__dirname, '..', '..', '..');
const BINARY = '/tmp/go-daisy-e2e';
const PORT = process.env.GALLERY_E2E_PORT || '11001';
const PID_FILE = '/tmp/go-daisy-e2e.pid';
const LOG_FILE = '/tmp/go-daisy-e2e.log';

async function globalSetup(): Promise<void> {
  // Clean previous state
  for (const f of [PID_FILE, LOG_FILE]) {
    try { fs.unlinkSync(f); } catch { /* ignore */ }
  }

  // 1. Build gallery binary
  execSync(`go build -o ${BINARY} ./cmd/gallery`, {
    cwd: PROJECT_ROOT,
    stdio: 'inherit',
  });

  // 2. Start gallery on test port
  const logFd = fs.openSync(LOG_FILE, 'a');
  const proc = spawn(BINARY, [], {
    cwd: PROJECT_ROOT,
    env: { ...process.env, GALLERY_PORT: PORT },
    stdio: ['ignore', logFd, logFd],
    detached: false,
  });
  fs.writeFileSync(PID_FILE, String(proc.pid!));

  // 3. Wait for /gallery to respond
  const maxWait = 15_000;
  const pollMs = 300;
  const deadline = Date.now() + maxWait;
  let ready = false;
  while (Date.now() < deadline) {
    try {
      const res = await fetch(`http://localhost:${PORT}/gallery`);
      if (res.ok) { ready = true; break; }
    } catch { /* not ready yet */ }
    await new Promise(r => setTimeout(r, pollMs));
  }
  if (!ready) {
    const log = fs.readFileSync(LOG_FILE, 'utf8').slice(-2000);
    console.error('Gallery logs:\n', log);
    throw new Error(`Gallery did not become ready within ${maxWait}ms`);
  }
  console.log(`global-setup: gallery ready on port ${PORT} (PID ${proc.pid})`);
}

export default globalSetup;
