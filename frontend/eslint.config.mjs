import nextConfig from "eslint-config-next/core-web-vitals";

const eslintConfig = [
  ...nextConfig,
  {
    ignores: [".next/", "node_modules/", "coverage/", "e2e/", "playwright-report/", "test-results/"],
  },
];

export default eslintConfig;
