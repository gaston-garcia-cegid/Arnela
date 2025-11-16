import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
  const path = request.nextUrl.pathname;

  // Public paths that don't require authentication
  const isPublicPath = path === '/' || path.startsWith('/sobre-') || path.startsWith('/intervencion') || path.startsWith('/formacion') || path.startsWith('/convenios') || path.startsWith('/contacto') || path.startsWith('/aviso-') || path.startsWith('/politica-') || path.startsWith('/accesibilidad') || path.startsWith('/mapa-');

  // Get auth token from cookie or local storage (via header)
  // Note: Since we're using localStorage in Zustand, we need to check on client side
  // This is a basic example - in production, use httpOnly cookies
  const token = request.cookies.get('auth-token')?.value;

  // If trying to access dashboard without token, redirect to home
  if (path.startsWith('/dashboard') && !token) {
    return NextResponse.redirect(new URL('/', request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    /*
     * Match all request paths except:
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - public folder
     */
    '/((?!_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)',
  ],
};
