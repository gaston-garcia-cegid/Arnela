import { describe, it, expect } from 'vitest';
import { cn } from '../utils';

describe('utils', () => {
    describe('cn', () => {
        it('merges class names correctly', () => {
            expect(cn('foo', 'bar')).toBe('foo bar');
        });

        it('handles conditional classes', () => {
            expect(cn('foo', true && 'bar', false && 'baz')).toBe('foo bar');
        });

        it('merges tailwind classes (overrides)', () => {
            // p-4 should override px-2 py-1
            expect(cn('px-2 py-1', 'p-4')).toBe('p-4');

            // text-blue-500 should override text-red-500
            expect(cn('text-red-500', 'text-blue-500')).toBe('text-blue-500');
        });

        it('handles arrays and objects', () => {
            expect(cn(['a', 'b'], { c: true, d: false })).toBe('a b c');
        });
    });
});
