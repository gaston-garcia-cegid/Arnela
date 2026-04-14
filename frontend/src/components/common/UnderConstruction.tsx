import Link from 'next/link';
import { Button } from '@/components/ui/button';

interface UnderConstructionProps {
  readonly title: string;
  readonly description?: string;
}

export function UnderConstruction({
  title,
  description = 'Estamos trabajando para traerte este contenido muy pronto.',
}: UnderConstructionProps) {
  return (
    <section className="flex min-h-[65vh] flex-col items-center justify-center px-4 py-16 md:py-24">
      <div className="mx-auto max-w-lg text-center">
        {/* Illustration */}
        <div className="mb-8 flex justify-center">
          <svg
            width="180"
            height="180"
            viewBox="0 0 180 180"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            className="drop-shadow-md"
            aria-hidden="true"
          >
            {/* Circular background */}
            <circle cx="90" cy="90" r="88" fill="#f5ede4" stroke="#e8ddd0" strokeWidth="2" />

            {/* Easel body */}
            <line x1="60" y1="140" x2="90" y2="60" stroke="#c89f7b" strokeWidth="3" strokeLinecap="round" />
            <line x1="120" y1="140" x2="90" y2="60" stroke="#c89f7b" strokeWidth="3" strokeLinecap="round" />
            <line x1="90" y1="80" x2="90" y2="150" stroke="#c89f7b" strokeWidth="3" strokeLinecap="round" />

            {/* Canvas on easel */}
            <rect x="62" y="48" width="56" height="48" rx="3" fill="#ffffff" stroke="#d4936d" strokeWidth="2" />

            {/* Decorative strokes on canvas */}
            <line x1="70" y1="62" x2="88" y2="62" stroke="#e8c4a1" strokeWidth="3" strokeLinecap="round" />
            <line x1="70" y1="72" x2="110" y2="72" stroke="#d4936d" strokeWidth="2" strokeLinecap="round" />
            <line x1="70" y1="80" x2="100" y2="80" stroke="#e8c4a1" strokeWidth="2" strokeLinecap="round" />

            {/* Paint brush */}
            <g transform="translate(118, 38) rotate(25)">
              <rect x="0" y="0" width="6" height="32" rx="2" fill="#c89f7b" />
              <rect x="-1" y="28" width="8" height="14" rx="1" fill="#d4936d" />
              <rect x="0" y="38" width="6" height="4" rx="1" fill="#e8c4a1" />
            </g>

            {/* Small decorative dots */}
            <circle cx="42" cy="52" r="4" fill="#e8c4a1" opacity="0.6" />
            <circle cx="148" cy="72" r="3" fill="#d4936d" opacity="0.4" />
            <circle cx="38" cy="120" r="3" fill="#c89f7b" opacity="0.5" />
          </svg>
        </div>

        <h1 className="mb-3 text-2xl font-bold tracking-tight text-primary sm:text-3xl md:text-4xl">
          {title}
        </h1>

        <div className="mb-2 inline-block rounded-full bg-accent/40 px-4 py-1.5 text-sm font-medium text-primary">
          En construcción
        </div>

        <p className="mt-4 text-base text-foreground/70 md:text-lg">
          {description}
        </p>

        <p className="mt-2 text-sm text-muted-foreground">
          Mientras tanto, no dudes en contactarnos si necesitas información.
        </p>

        <div className="mt-8 flex flex-col items-center gap-3 sm:flex-row sm:justify-center">
          <Button asChild size="lg">
            <Link href="/">Volver al inicio</Link>
          </Button>
          <Button asChild variant="outline" size="lg">
            <Link href="tel:611749043">Llamar: 611 749 043</Link>
          </Button>
        </div>
      </div>
    </section>
  );
}
