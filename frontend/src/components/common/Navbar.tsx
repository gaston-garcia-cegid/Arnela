'use client';

import { Button } from '@/components/ui/button';
import Link from 'next/link';
import { useState } from 'react';

interface NavbarProps {
  readonly onLoginClick?: () => void;
}

export function Navbar({ onLoginClick }: Readonly<NavbarProps>) {
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const navLinks = [
    { href: '/sobre-arnela', label: 'SOBRE ARNELA' },
    { href: '/intervencion', label: 'INTERVENCIÓN' },
    { href: '/formacion', label: 'FORMACIÓN' },
    { href: '/convenios-y-colaboraciones', label: 'CONVENIOS Y COLABORACIONES' },
    { href: '/contacto', label: 'CONTACTO' },
  ];

  return (
    <nav className="sticky top-0 z-50 border-b border-primary bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/80 shadow-sm">
      <div className="container mx-auto flex h-16 items-center justify-between px-4">
        <Link href="/" className="text-xl font-bold text-primary hover:text-accent transition-colors">
          Arnela Gabinete
        </Link>

        {/* Desktop Navigation */}
        <div className="hidden items-center gap-6 md:flex">
          {navLinks.map((link) => (
            <Link
              key={link.href}
              href={link.href}
              className="text-sm font-medium text-primary/80 transition-colors hover:text-accent"
            >
              {link.label}
            </Link>
          ))}
          {onLoginClick && (
            <Button onClick={onLoginClick} variant="default" size="sm" className="ml-2 bg-primary text-background hover:bg-accent hover:text-primary">
              Iniciar sesión
            </Button>
          )}
        </div>

        {/* Mobile Menu Button */}
        <button
          className="md:hidden text-primary"
          onClick={() => setIsMenuOpen(!isMenuOpen)}
          aria-label="Toggle menu"
        >
          <svg
            className="h-6 w-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            {isMenuOpen ? (
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />
            ) : (
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4 6h16M4 12h16M4 18h16"
              />
            )}
          </svg>
        </button>
      </div>

      {/* Mobile Menu */}
      {isMenuOpen && (
        <div className="border-t bg-background md:hidden">
          <div className="container mx-auto flex flex-col space-y-4 px-4 py-4">
            {navLinks.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                className="text-sm font-medium transition-colors hover:text-primary"
                onClick={() => setIsMenuOpen(false)}
              >
                {link.label}
              </Link>
            ))}
            {onLoginClick && (
              <Button onClick={onLoginClick} variant="default" className="w-full">
                Iniciar sesión
              </Button>
            )}
          </div>
        </div>
      )}
    </nav>
  );
}
