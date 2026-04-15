'use client';

import { useState } from 'react';
import { Navbar } from '@/components/common/Navbar';
import { Footer } from '@/components/common/Footer';
import { LoginModal } from '@/components/auth/LoginModal';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import { Card, CardContent } from '@/components/ui/card';
import { toast } from 'sonner';

export default function ContactoPage() {
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    phone: '',
    subject: '',
    message: '',
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.name || !formData.email || !formData.message) {
      toast.error('Por favor, rellena los campos obligatorios.');
      return;
    }

    setIsSubmitting(true);
    try {
      const mailtoLink = `mailto:arnelagabinete@gmail.com?subject=${encodeURIComponent(
        formData.subject || `Contacto de ${formData.name}`
      )}&body=${encodeURIComponent(
        `Nombre: ${formData.name}\nEmail: ${formData.email}\nTeléfono: ${formData.phone || 'No proporcionado'}\n\n${formData.message}`
      )}`;
      globalThis.location.href = mailtoLink;

      toast.success('Abriendo tu cliente de correo...');
      setFormData({ name: '', email: '', phone: '', subject: '', message: '' });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <>
      <Navbar onLoginClick={() => setIsLoginModalOpen(true)} />
      <main className="bg-background">
        {/* Header */}
        <section className="bg-muted/30 py-16 px-4 md:py-20">
          <div className="container mx-auto text-center">
            <h1 className="mb-4 text-3xl font-bold tracking-tight text-primary sm:text-4xl md:text-5xl">
              Contacto
            </h1>
            <p className="mx-auto max-w-2xl text-lg text-foreground/70">
              ¿Tienes alguna consulta? Estamos aquí para ayudarte. Escríbenos y nos pondremos en
              contacto contigo lo antes posible.
            </p>
          </div>
        </section>

        {/* Content */}
        <section className="py-16 px-4">
          <div className="container mx-auto grid gap-12 lg:grid-cols-2">
            {/* Contact Form */}
            <div>
              <h2 className="mb-6 text-2xl font-bold text-foreground">Envíanos un mensaje</h2>
              <form onSubmit={handleSubmit} className="space-y-5">
                <div className="grid gap-5 sm:grid-cols-2">
                  <div className="space-y-2">
                    <Label htmlFor="name">Nombre *</Label>
                    <Input
                      id="name"
                      name="name"
                      value={formData.name}
                      onChange={handleChange}
                      placeholder="Tu nombre completo"
                      required
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="email">Email *</Label>
                    <Input
                      id="email"
                      name="email"
                      type="email"
                      value={formData.email}
                      onChange={handleChange}
                      placeholder="tu@email.com"
                      required
                    />
                  </div>
                </div>

                <div className="grid gap-5 sm:grid-cols-2">
                  <div className="space-y-2">
                    <Label htmlFor="phone">Teléfono</Label>
                    <Input
                      id="phone"
                      name="phone"
                      type="tel"
                      value={formData.phone}
                      onChange={handleChange}
                      placeholder="611 749 043"
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="subject">Asunto</Label>
                    <Input
                      id="subject"
                      name="subject"
                      value={formData.subject}
                      onChange={handleChange}
                      placeholder="¿Sobre qué nos escribes?"
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="message">Mensaje *</Label>
                  <Textarea
                    id="message"
                    name="message"
                    value={formData.message}
                    onChange={handleChange}
                    placeholder="Cuéntanos en qué podemos ayudarte..."
                    rows={6}
                    required
                  />
                </div>

                <Button type="submit" size="lg" disabled={isSubmitting} className="w-full sm:w-auto">
                  {isSubmitting ? 'Enviando...' : 'Enviar mensaje'}
                </Button>
              </form>
            </div>

            {/* Contact Info */}
            <div className="space-y-6">
              <h2 className="mb-6 text-2xl font-bold text-foreground">Información de contacto</h2>

              <Card className="border-l-4 border-l-primary">
                <CardContent className="p-6 space-y-4">
                  <div className="flex items-start gap-4">
                    <span className="text-2xl mt-0.5">📍</span>
                    <div>
                      <p className="font-semibold text-foreground">Dirección</p>
                      <p className="text-foreground/70">Calle García Barbón 30</p>
                      <p className="text-foreground/70">36202 Vigo, Pontevedra</p>
                    </div>
                  </div>

                  <div className="flex items-start gap-4">
                    <span className="text-2xl mt-0.5">📞</span>
                    <div>
                      <p className="font-semibold text-foreground">Teléfono</p>
                      <a href="tel:611749043" className="text-primary hover:underline">
                        611 749 043
                      </a>
                    </div>
                  </div>

                  <div className="flex items-start gap-4">
                    <span className="text-2xl mt-0.5">✉️</span>
                    <div>
                      <p className="font-semibold text-foreground">Email</p>
                      <a
                        href="mailto:arnelagabinete@gmail.com"
                        className="text-primary hover:underline"
                      >
                        arnelagabinete@gmail.com
                      </a>
                    </div>
                  </div>

                  <div className="flex items-start gap-4">
                    <span className="text-2xl mt-0.5">🕐</span>
                    <div>
                      <p className="font-semibold text-foreground">Horario</p>
                      <p className="text-foreground/70">Lunes a viernes</p>
                      <p className="text-foreground/70">10:00 - 14:00 y 16:00 - 20:00</p>
                    </div>
                  </div>

                  <div className="flex items-start gap-4">
                    <span className="text-2xl mt-0.5">📷</span>
                    <div>
                      <p className="font-semibold text-foreground">Instagram</p>
                      <a
                        href="https://www.instagram.com/arnelagabinete/"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-primary hover:underline"
                      >
                        @arnelagabinete
                      </a>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Map */}
              <Card>
                <CardContent className="p-0 overflow-hidden rounded-lg">
                  <iframe
                    title="Ubicación de Arnela Gabinete"
                    src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d2956.8!2d-8.7263!3d42.2329!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0xd2f621b7ef6df43%3A0x9f9f9f9f9f9f9f9f!2sCalle%20Garc%C3%ADa%20Barb%C3%B3n%2030%2C%2036202%20Vigo%2C%20Pontevedra!5e0!3m2!1ses!2ses!4v1700000000000"
                    width="100%"
                    height="300"
                    style={{ border: 0 }}
                    allowFullScreen
                    loading="lazy"
                    referrerPolicy="no-referrer-when-downgrade"
                  />
                </CardContent>
              </Card>
            </div>
          </div>
        </section>
      </main>
      <Footer />
      <LoginModal isOpen={isLoginModalOpen} onClose={() => setIsLoginModalOpen(false)} />
    </>
  );
}
