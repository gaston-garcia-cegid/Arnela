'use client';

import { useState } from 'react';
import { Navbar } from '@/components/common/Navbar';
import { Footer } from '@/components/common/Footer';
import { UnderConstruction } from '@/components/common/UnderConstruction';
import { LoginModal } from '@/components/auth/LoginModal';

export default function IntervencionPage() {
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);

  return (
    <>
      <Navbar onLoginClick={() => setIsLoginModalOpen(true)} />
      <main>
        <UnderConstruction
          title="Intervención"
          description="Estamos preparando toda la información sobre nuestros servicios de intervención terapéutica y educativa."
        />
      </main>
      <Footer />
      <LoginModal isOpen={isLoginModalOpen} onClose={() => setIsLoginModalOpen(false)} />
    </>
  );
}
