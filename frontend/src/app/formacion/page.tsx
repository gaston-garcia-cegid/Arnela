'use client';

import { useState } from 'react';
import { Navbar } from '@/components/common/Navbar';
import { Footer } from '@/components/common/Footer';
import { UnderConstruction } from '@/components/common/UnderConstruction';
import { LoginModal } from '@/components/auth/LoginModal';

export default function FormacionPage() {
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);

  return (
    <>
      <Navbar onLoginClick={() => setIsLoginModalOpen(true)} />
      <main>
        <UnderConstruction
          title="Formación"
          description="En breve encontrarás aquí toda la oferta formativa de Arnela Gabinete."
        />
      </main>
      <Footer />
      <LoginModal isOpen={isLoginModalOpen} onClose={() => setIsLoginModalOpen(false)} />
    </>
  );
}
