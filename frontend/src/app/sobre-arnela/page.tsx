'use client';

import { useState } from 'react';
import { Navbar } from '@/components/common/Navbar';
import { Footer } from '@/components/common/Footer';
import { UnderConstruction } from '@/components/common/UnderConstruction';
import { LoginModal } from '@/components/auth/LoginModal';

export default function SobreArnelaPage() {
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);

  return (
    <>
      <Navbar onLoginClick={() => setIsLoginModalOpen(true)} />
      <main>
        <UnderConstruction
          title="Sobre Arnela"
          description="Pronto podrás conocer más sobre nuestro equipo, nuestra historia y nuestra forma de trabajar."
        />
      </main>
      <Footer />
      <LoginModal isOpen={isLoginModalOpen} onClose={() => setIsLoginModalOpen(false)} />
    </>
  );
}
