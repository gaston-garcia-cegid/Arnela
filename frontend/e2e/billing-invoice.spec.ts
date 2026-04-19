import { test, expect } from "@playwright/test";
import { loginBackoffice } from "./helpers/loginBackoffice";

test.describe("billing invoices", () => {
  test("listado de facturas visible para admin", async ({ page }) => {
    const email = process.env.E2E_ADMIN_EMAIL;
    const password = process.env.E2E_ADMIN_PASSWORD;
    test.skip(!email || !password, "Definir E2E_ADMIN_EMAIL y E2E_ADMIN_PASSWORD (ver e2e/.env.example)");

    await loginBackoffice(page, email!, password!);
    await page.goto("/dashboard/backoffice/billing/invoices");
    await expect(page.getByRole("heading", { name: /^facturas$/i })).toBeVisible();
  });
});
