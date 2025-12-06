import { describe, it, expect } from 'vitest';
import {
    validateDNI,
    validateCIF,
    validateDNIorCIF,
    validateEmail,
    validatePhone,
    validatePostalCode,
    validatePassword
} from '../validators';

describe('Centralized Validators', () => {
    describe('validateDNI', () => {
        it('returns valid for correct DNI', () => {
            // 12345678Z (Z is 78 mod 23 = 14 => Z)
            expect(validateDNI('12345678Z').isValid).toBe(true);
        });

        it('returns invalid for wrong letter', () => {
            expect(validateDNI('12345678A').isValid).toBe(false);
            expect(validateDNI('12345678A').error).toContain('La letra del DNI/NIE no es correcta');
        });

        it('returns valid for correct NIE (X)', () => {
            // X1234567L (X=0 => 01234567L)
            expect(validateDNI('X1234567L').isValid).toBe(true);
        });

        it('returns invalid format for weird chars', () => {
            expect(validateDNI('ABC').isValid).toBe(false);
        });

        it('is case insensitive', () => {
            expect(validateDNI('12345678z').isValid).toBe(true);
        });
    });

    describe('validateCIF', () => {
        it('returns valid for correct CIF format', () => {
            expect(validateCIF('A12345678').isValid).toBe(true);
        });

        it('returns invalid format for bad CIF', () => {
            expect(validateCIF('12345678A').isValid).toBe(false);
        });
    });

    describe('validateDNIorCIF', () => {
        it('validates DNI', () => {
            expect(validateDNIorCIF('12345678Z').isValid).toBe(true);
        });

        it('validates CIF', () => {
            expect(validateDNIorCIF('A12345678').isValid).toBe(true);
        });

        it('rejects invalid', () => {
            expect(validateDNIorCIF('INVALID').isValid).toBe(false);
        });
    });

    describe('validateEmail', () => {
        it('validates simple email', () => {
            expect(validateEmail('test@test.com').isValid).toBe(true);
        });

        it('rejects no @', () => {
            expect(validateEmail('testtest.com').isValid).toBe(false);
        });
    });

    describe('validatePhone', () => {
        it('validates mobile', () => {
            expect(validatePhone('600123456').isValid).toBe(true);
        });

        it('validates landline (9XX)', () => {
            expect(validatePhone('912345678').isValid).toBe(true);
        });

        it('rejects text', () => {
            expect(validatePhone('hola').isValid).toBe(false);
        });

        it('accepts spaces and dashes', () => {
            expect(validatePhone('600 123 456').isValid).toBe(true);
            expect(validatePhone('600-123-456').isValid).toBe(true);
        });
    });

    describe('validatePostalCode', () => {
        it('validates 28001', () => {
            expect(validatePostalCode('28001').isValid).toBe(true);
        });

        it('rejects 60000', () => {
            expect(validatePostalCode('60000').isValid).toBe(false);
        });
    });

    describe('validatePassword', () => {
        it('validates strong password', () => {
            expect(validatePassword('Pass1234').isValid).toBe(true);
        });

        it('rejects weak password (no number)', () => {
            expect(validatePassword('Password').isValid).toBe(false);
        });

        it('rejects short password', () => {
            expect(validatePassword('Pass1').isValid).toBe(false);
        });
    });
});
