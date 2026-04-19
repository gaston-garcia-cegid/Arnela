import { test, expect } from "@playwright/test";
import { loginBackoffice } from "./helpers/loginBackoffice";

test.describe("auth backoffice", () => {
  test("login admin llega al panel backoffice", async ({ page }) => {
    const email = process.env.E2E_ADMIN_EMAIL;
    const password = process.env.E2E_ADMIN_PASSWORD;
    test.skip(!email || !password, "Definir E2E_ADMIN_EMAIL y E2E_ADMIN_PASSWORD (ver e2e/.env.example)");

    await loginBackoffice(page, email!, password!);
    await expect(page).toHaveURL(/dashboard\/backoffice/);
  });
});
