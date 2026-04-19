import { test, expect } from "@playwright/test";
import { loginBackoffice } from "./helpers/loginBackoffice";

test.describe("appointments backoffice", () => {
  test("página de citas carga tras login", async ({ page }) => {
    const email = process.env.E2E_ADMIN_EMAIL;
    const password = process.env.E2E_ADMIN_PASSWORD;
    test.skip(!email || !password, "Definir E2E_ADMIN_EMAIL y E2E_ADMIN_PASSWORD (ver e2e/.env.example)");

    await loginBackoffice(page, email!, password!);
    await page.goto("/dashboard/backoffice/appointments");
    await expect(
      page.getByRole("heading", { name: /^(mis )?citas$/i })
    ).toBeVisible();
  });
});
