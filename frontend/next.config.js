/**
 * `output: 'standalone'` copies the traced server graph with symlinks. On Windows,
 * creating symlinks often fails with EPERM unless Developer Mode is on or the
 * process has SeCreateSymbolicLinkPrivilege. Linux/macOS (Dockerfile, GitHub
 * Actions) are unaffected.
 *
 * Local Windows: omit standalone so `pnpm build` succeeds.
 * Opt in: set FORCE_STANDALONE=1 after enabling Developer Mode (or run build in WSL/Docker).
 */
const isWindows = process.platform === 'win32'
const forceStandalone = process.env.FORCE_STANDALONE === '1'
const useStandalone = !isWindows || forceStandalone

/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  ...(useStandalone ? { output: 'standalone' } : {}),
}

module.exports = nextConfig
