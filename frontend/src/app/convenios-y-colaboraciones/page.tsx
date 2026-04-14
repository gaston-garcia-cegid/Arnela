'use client';

import { useState } from 'react';
import { Navbar } from '@/components/common/Navbar';
import { Footer } from '@/components/common/Footer';
import { UnderConstruction } from '@/components/common/UnderConstruction';
import { LoginModal } from '@/components/auth/LoginModal';

export default function ConveniosPage() {
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);

  return (
    <>
      <Navbar onLoginClick={() => setIsLoginModalOpen(true)} />
      <main>
        <UnderConstruction
          title="Convenios y Colaboraciones"
          description="Pronto publicaremos información sobre nuestros convenios y entidades colaboradoras."
        />
      </main>
      <Footer />
      <LoginModal isOpen={isLoginModalOpen} onClose={() => setIsLoginModalOpen(false)} />
    </>
  );
}
