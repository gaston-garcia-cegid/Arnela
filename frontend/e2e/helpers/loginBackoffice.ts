import { expect, type Page } from "@playwright/test";

/**
 * Abre la home, usa el modal de login del Navbar y espera una ruta bajo `/dashboard/backoffice`.
 * Requiere usuario admin o employee (ambos caen bajo backoffice tras login).
 */
export async function loginBackoffice(
  page: Page,
  email: string,
  password: string
): Promise<void> {
  await page.goto("/");
  await page.getByRole("button", { name: /iniciar sesión/i }).click();
  await expect(page.getByRole("dialog")).toBeVisible();
  await page.getByLabel("Email").fill(email);
  await page.getByLabel("Contraseña").fill(password);
  await page.getByRole("button", { name: /^ingresar$/i }).click();
  await expect(page).toHaveURL(/\/dashboard\/backoffice/);
}
