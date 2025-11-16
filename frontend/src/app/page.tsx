'use client';

import { useState } from 'react';
import { Navbar } from '@/components/common/Navbar';
import { Footer } from '@/components/common/Footer';
import { Hero } from '@/components/landing/Hero';
import { About } from '@/components/landing/About';
import { Services } from '@/components/landing/Services';
import { Testimonial } from '@/components/landing/Testimonial';
import { Reviews } from '@/components/landing/Reviews';
import { LoginModal } from '@/components/auth/LoginModal';

export default function HomePage() {
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);

  return (
    <>
      <Navbar onLoginClick={() => setIsLoginModalOpen(true)} />
      <main>
        <Hero />
        <About />
        <Services />
        <Testimonial />
        <Reviews />
      </main>
      <Footer />
      <LoginModal isOpen={isLoginModalOpen} onClose={() => setIsLoginModalOpen(false)} />
    </>
  );
}
