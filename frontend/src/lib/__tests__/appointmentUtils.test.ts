import { describe, it, expect, vi, afterEach, beforeEach } from 'vitest';
import {
    formatAppointmentDate,
    formatAppointmentTime,
    formatAppointmentDateTime,
    getRelativeDateLabel,
    formatDateForAPI,
    formatDateTimeForAPI,
    isAppointmentPast,
    isAppointmentUpcoming,
    getStatusColor,
    getStatusLabel,
    formatDuration,
    groupAppointmentsByDate,
    generateTimeSlots,
    getNext7Days,
    isWeekday,
    parseTimeSlot
} from '../appointmentUtils';

describe('appointmentUtils - Core Logic', () => {
    beforeEach(() => {
        vi.useFakeTimers();
    });

    afterEach(() => {
        vi.useRealTimers();
    });

    // --- Formatting Tests ---
    describe('Date Formatting', () => {
        it('formatAppointmentDate returns readable Spanish date', () => {
            // 2025-11-22 is Saturday
            const iso = '2025-11-22T10:00:00';
            expect(formatAppointmentDate(iso)).toBe('sábado, 22 de noviembre de 2025');
        });

        it('formatAppointmentTime returns 12-hour format', () => {
            expect(formatAppointmentTime('2025-11-22T13:30:00')).toBe('1:30 PM');
            expect(formatAppointmentTime('2025-11-22T09:00:00')).toBe('9:00 AM');
        });

        it('formatAppointmentDateTime returns combined format', () => {
            const result = formatAppointmentDateTime('2025-11-22T13:30:00');
            // Check parts because specific output might vary slightly by locale/env
            expect(result).toContain('sáb');
            expect(result).toContain('22 nov');
            expect(result).toContain('1:30 PM');
        });

        it('formatDateForAPI returns YYYY-MM-DD', () => {
            const date = new Date('2025-12-31T10:00:00');
            expect(formatDateForAPI(date)).toBe('2025-12-31');
        });

        it('formatDateTimeForAPI returns ISO string', () => {
            const date = new Date('2025-12-31T10:00:00Z'); // UTC
            expect(formatDateTimeForAPI(date)).toBe(date.toISOString());
        });
    });

    // --- Relative Date Logic ---
    describe('getRelativeDateLabel', () => {
        it('returns "Hoy" for current date', () => {
            vi.setSystemTime(new Date('2023-01-01T12:00:00'));
            // Same day
            expect(getRelativeDateLabel('2023-01-01T15:00:00')).toBe('Hoy');
        });

        it('returns "Mañana" for next day', () => {
            vi.setSystemTime(new Date('2023-01-01T12:00:00'));
            // Next day
            expect(getRelativeDateLabel('2023-01-02T10:00:00')).toBe('Mañana');
        });

        it('returns formatted date for other days', () => {
            vi.setSystemTime(new Date('2023-01-01T12:00:00'));
            const label = getRelativeDateLabel('2023-01-05T10:00:00');
            expect(label).toContain('ene'); // Jan in Spanish
        });
    });

    // --- Time Check Logic ---
    describe('Past/Upcoming Checks', () => {
        it('identifies past appointments', () => {
            vi.setSystemTime(new Date('2023-06-15T12:00:00'));
            expect(isAppointmentPast('2023-06-15T10:00:00')).toBe(true);
            expect(isAppointmentPast('2023-06-15T14:00:00')).toBe(false);
        });

        it('identifies upcoming appointments', () => {
            vi.setSystemTime(new Date('2023-06-15T12:00:00'));
            expect(isAppointmentUpcoming('2023-06-15T14:00:00')).toBe(true);
            expect(isAppointmentUpcoming('2023-06-15T10:00:00')).toBe(false);
        });
    });

    // --- Status Logic ---
    describe('Status Helpers', () => {
        it('returns correct colors', () => {
            expect(getStatusColor('confirmed')).toContain('green');
            expect(getStatusColor('cancelled')).toContain('red');
            expect(getStatusColor('pending')).toContain('yellow');
            expect(getStatusColor('completed')).toContain('blue');
            expect(getStatusColor('rescheduled')).toContain('purple');

            // Edge case: unknown status
            expect(getStatusColor('unknown' as any)).toContain('gray');
        });

        it('returns correct labels', () => {
            expect(getStatusLabel('confirmed')).toBe('Confirmada');
            expect(getStatusLabel('cancelled')).toBe('Cancelada');

            // Edge case: return status if not mapped
            expect(getStatusLabel('unknown' as any)).toBe('unknown');
        });
    });

    // --- Duration Logic ---
    describe('formatDuration', () => {
        it('formats minutes < 60', () => {
            expect(formatDuration(45)).toBe('45 min');
        });
        it('formats exact hours', () => {
            expect(formatDuration(120)).toBe('2h');
        });
        it('formats hours and minutes', () => {
            expect(formatDuration(130)).toBe('2h 10min');
        });
        it('handles 0', () => {
            expect(formatDuration(0)).toBe('0 min');
        });
    });

    // --- Slots & Days Logic ---
    describe('Scheduling Logic', () => {
        it('generateTimeSlots generates 9AM to 6PM in 15min intervals', () => {
            const date = new Date('2023-01-01T00:00:00');
            const slots = generateTimeSlots(date);

            // 9h * 4 slots/h = 36 slots
            expect(slots).toHaveLength(36);

            // Should start at 9:00
            expect(slots[0].getHours()).toBe(9);
            expect(slots[0].getMinutes()).toBe(0);

            // Should end at 17:45 (last slot start)
            const last = slots[slots.length - 1];
            expect(last.getHours()).toBe(17);
            expect(last.getMinutes()).toBe(45);
        });

        it('getNext7Days returns 7 consecutive days', () => {
            vi.setSystemTime(new Date('2023-01-01T12:00:00'));
            const days = getNext7Days();

            expect(days).toHaveLength(7);
            expect(days[0].getDate()).toBe(1); // Today
            expect(days[1].getDate()).toBe(2);
            expect(days[6].getDate()).toBe(7);
        });

        it('isWeekday identifies valid weekdays', () => {
            // 2023-01-01 is Sunday
            expect(isWeekday(new Date('2023-01-01'))).toBe(false);
            // 2023-01-02 is Monday
            expect(isWeekday(new Date('2023-01-02'))).toBe(true);
            // 2023-01-06 is Friday
            expect(isWeekday(new Date('2023-01-06'))).toBe(true);
            // 2023-01-07 is Saturday
            expect(isWeekday(new Date('2023-01-07'))).toBe(false);
        });

        it('parseTimeSlot returns Date object', () => {
            const str = '2023-01-01T10:00:00';
            const date = parseTimeSlot(str);
            expect(date).toBeInstanceOf(Date);
            expect(date.toISOString()).toContain('2023-01-01');
        });
    });

    // --- Array Helpers ---
    describe('groupAppointmentsByDate', () => {
        it('groups appointments by YYYY-MM-DD', () => {
            const apps = [
                { id: 1, startTime: '2023-01-01T10:00:00' },
                { id: 2, startTime: '2023-01-01T14:00:00' },
                { id: 3, startTime: '2023-01-02T09:00:00' }
            ];

            const grouped = groupAppointmentsByDate(apps);

            expect(Object.keys(grouped)).toHaveLength(2);
            expect(grouped['2023-01-01']).toHaveLength(2);
            expect(grouped['2023-01-02']).toHaveLength(1);
            expect(grouped['2023-01-01'][0].id).toBe(1);
        });

        it('handles empty array', () => {
            expect(groupAppointmentsByDate([])).toEqual({});
        });
    });
});
