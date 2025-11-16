import Link from 'next/link';

export function Footer() {
  return (
    <footer className="border-t border-border/40 bg-muted/30">
      <div className="container mx-auto px-4 py-12">
        {/* EU Funding Logos */}
        <div className="mb-8 flex flex-wrap items-center justify-center gap-4 border-b border-border/40 pb-8">
          <div className="text-xs text-muted-foreground md:text-sm">
            🇪🇺 Financiado por la Unión Europea NextGenerationEU
          </div>
          <div className="text-xs text-muted-foreground md:text-sm">Gobierno de España</div>
          <div className="text-xs text-muted-foreground md:text-sm">
            Plan de Recuperación, Transformación y Resilicencia
          </div>
        </div>

        <div className="grid gap-8 md:grid-cols-3">
          {/* Contact Info */}
          <div>
            <h3 className="mb-4 font-bold text-primary">ARNELA GABINETE</h3>
            <div className="space-y-2 text-sm text-foreground/70">
              <p>Calle García Barbón 30</p>
              <p>36202 Vigo</p>
              <p className="flex items-center gap-2">
                <span>📞</span>
                <Link href="tel:611749043" className="hover:text-primary transition-colors">
                  611749043
                </Link>
              </p>
              <p className="flex items-center gap-2">
                <span>✉️</span>
                <Link
                  href="mailto:arnelagabinete@gmail.com"
                  className="hover:text-primary transition-colors"
                >
                  arnelagabinete@gmail.com
                </Link>
              </p>
              <p className="pt-3 flex items-center gap-2">
                <span>🕐</span>
                <span>De lunes a viernes</span>
              </p>
              <p className="pl-6">10:00 a 14:00 y de 16:00 a 20:00</p>
            </div>
          </div>

          {/* Main Links */}
          <div>
            <h3 className="mb-4 font-bold text-primary">NAVEGACIÓN</h3>
            <ul className="space-y-2 text-sm">
              <li>
                <Link href="/sobre-arnela" className="text-foreground/70 hover:text-primary transition-colors">
                  SOBRE ARNELA
                </Link>
              </li>
              <li>
                <Link href="/intervencion" className="text-foreground/70 hover:text-primary transition-colors">
                  INTERVENCIÓN
                </Link>
              </li>
              <li>
                <Link href="/formacion" className="text-foreground/70 hover:text-primary transition-colors">
                  FORMACIÓN
                </Link>
              </li>
              <li>
                <Link
                  href="/convenios-y-colaboraciones"
                  className="text-foreground/70 hover:text-primary transition-colors"
                >
                  CONVENIOS Y COLABORACIONES
                </Link>
              </li>
              <li>
                <Link href="/contacto" className="text-foreground/70 hover:text-primary transition-colors">
                  CONTACTO
                </Link>
              </li>
            </ul>
          </div>

          {/* Legal Links */}
          <div>
            <h3 className="mb-4 font-bold text-primary">LEGAL</h3>
            <ul className="space-y-2 text-sm">
              <li>
                <Link href="/aviso-legal" className="text-muted-foreground hover:text-primary">
                  Aviso legal
                </Link>
              </li>
              <li>
                <Link
                  href="/politica-de-privacidad"
                  className="text-muted-foreground hover:text-primary"
                >
                  Política de privacidad
                </Link>
              </li>
              <li>
                <Link
                  href="/politica-de-cookies"
                  className="text-muted-foreground hover:text-primary"
                >
                  Política de cookies
                </Link>
              </li>
              <li>
                <Link href="/accesibilidad" className="text-muted-foreground hover:text-primary">
                  Accesibilidad
                </Link>
              </li>
              <li>
                <Link href="/mapa-web" className="text-muted-foreground hover:text-primary">
                  Mapa web
                </Link>
              </li>
            </ul>
          </div>
        </div>

        {/* Bottom Bar */}
        <div className="mt-8 flex flex-col items-center justify-between gap-4 border-t border-border/40 pt-8 text-sm text-foreground/60 md:flex-row">
          <p className="text-xs md:text-sm">&copy; {new Date().getFullYear()} Arnela Gabinete. Todos los derechos reservados.</p>
          <div className="flex items-center gap-4">
            <Link
              href="https://www.instagram.com/arnelagabinete/"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 hover:text-primary transition-colors text-xs md:text-sm"
            >
              <span className="text-lg">📷</span> Instagram
            </Link>
          </div>
        </div>
      </div>
    </footer>
  );
}
