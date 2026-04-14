'use client';

import { useState } from 'react';
import { Navbar } from '@/components/common/Navbar';
import { Footer } from '@/components/common/Footer';
import { UnderConstruction } from '@/components/common/UnderConstruction';
import { LoginModal } from '@/components/auth/LoginModal';

export default function ContactoPage() {
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);

  return (
    <>
      <Navbar onLoginClick={() => setIsLoginModalOpen(true)} />
      <main>
        <UnderConstruction
          title="Contacto"
          description="Estamos preparando un formulario de contacto. Mientras tanto, puedes llamarnos o escribirnos directamente."
        />
      </main>
      <Footer />
      <LoginModal isOpen={isLoginModalOpen} onClose={() => setIsLoginModalOpen(false)} />
    </>
  );
}
