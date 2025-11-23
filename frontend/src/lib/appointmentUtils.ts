// Utility functions for appointment date/time handling

import { format, parseISO, isAfter, isBefore, isToday, isTomorrow, addDays, startOfDay } from 'date-fns';
import { es } from 'date-fns/locale';
import type { AppointmentStatus } from '@/types/appointment';

/**
 * Format appointment date to readable Spanish format
 * Example: "Lunes, 22 de noviembre de 2025"
 */
export function formatAppointmentDate(dateString: string): string {
  const date = parseISO(dateString);
  return format(date, "EEEE, d 'de' MMMM 'de' yyyy", { locale: es });
}

/**
 * Format appointment time to 12-hour format
 * Example: "10:00 AM"
 */
export function formatAppointmentTime(dateString: string): string {
  const date = parseISO(dateString);
  return format(date, 'h:mm a');
}

/**
 * Format appointment date and time combined
 * Example: "Lun, 22 nov - 10:00 AM"
 */
export function formatAppointmentDateTime(dateString: string): string {
  const date = parseISO(dateString);
  return format(date, "EEE, d MMM - h:mm a", { locale: es });
}

/**
 * Get relative date label (Hoy, Mañana, or date)
 */
export function getRelativeDateLabel(dateString: string): string {
  const date = parseISO(dateString);
  
  if (isToday(date)) {
    return 'Hoy';
  }
  
  if (isTomorrow(date)) {
    return 'Mañana';
  }
  
  return format(date, "EEE, d MMM", { locale: es });
}

/**
 * Format date to YYYY-MM-DD for API requests
 */
export function formatDateForAPI(date: Date): string {
  return format(date, 'yyyy-MM-dd');
}

/**
 * Format date to ISO 8601 string for API requests
 */
export function formatDateTimeForAPI(date: Date): string {
  return date.toISOString();
}

/**
 * Check if appointment is in the past
 */
export function isAppointmentPast(dateString: string): boolean {
  const date = parseISO(dateString);
  return isBefore(date, new Date());
}

/**
 * Check if appointment is upcoming (in the future)
 */
export function isAppointmentUpcoming(dateString: string): boolean {
  const date = parseISO(dateString);
  return isAfter(date, new Date());
}

/**
 * Get appointment status badge color
 */
export function getStatusColor(status: AppointmentStatus): string {
  const colors: Record<AppointmentStatus, string> = {
    pending: 'bg-yellow-100 text-yellow-800 border-yellow-200',
    confirmed: 'bg-green-100 text-green-800 border-green-200',
    cancelled: 'bg-red-100 text-red-800 border-red-200',
    completed: 'bg-blue-100 text-blue-800 border-blue-200',
    rescheduled: 'bg-purple-100 text-purple-800 border-purple-200',
  };
  
  return colors[status] || 'bg-gray-100 text-gray-800 border-gray-200';
}

/**
 * Get appointment status label in Spanish
 */
export function getStatusLabel(status: AppointmentStatus): string {
  const labels: Record<AppointmentStatus, string> = {
    pending: 'Pendiente',
    confirmed: 'Confirmada',
    cancelled: 'Cancelada',
    completed: 'Completada',
    rescheduled: 'Reprogramada',
  };
  
  return labels[status] || status;
}

/**
 * Generate time slots for a day (9:00 AM - 6:00 PM, 15min intervals)
 */
export function generateTimeSlots(date: Date): Date[] {
  const slots: Date[] = [];
  const baseDate = startOfDay(date);
  
  // Start at 9:00 AM, end at 6:00 PM
  for (let hour = 9; hour < 18; hour++) {
    for (let minute = 0; minute < 60; minute += 15) {
      const slotDate = new Date(baseDate);
      slotDate.setHours(hour, minute, 0, 0);
      slots.push(slotDate);
    }
  }
  
  return slots;
}

/**
 * Get next 7 days starting from today
 */
export function getNext7Days(): Date[] {
  const days: Date[] = [];
  const today = startOfDay(new Date());
  
  for (let i = 0; i < 7; i++) {
    days.push(addDays(today, i));
  }
  
  return days;
}

/**
 * Check if date is a weekday (Monday-Friday)
 */
export function isWeekday(date: Date): boolean {
  const day = date.getDay();
  return day >= 1 && day <= 5; // 0 = Sunday, 6 = Saturday
}

/**
 * Format duration in minutes to readable format
 */
export function formatDuration(minutes: number): string {
  if (minutes < 60) {
    return `${minutes} min`;
  }
  
  const hours = Math.floor(minutes / 60);
  const remainingMinutes = minutes % 60;
  
  if (remainingMinutes === 0) {
    return `${hours}h`;
  }
  
  return `${hours}h ${remainingMinutes}min`;
}

/**
 * Parse time slot string (ISO 8601) to Date object
 */
export function parseTimeSlot(slotString: string): Date {
  return parseISO(slotString);
}

/**
 * Group appointments by date
 */
export function groupAppointmentsByDate<T extends { startTime: string }>(
  appointments: T[]
): Record<string, T[]> {
  return appointments.reduce((groups, appointment) => {
    const date = formatDateForAPI(parseISO(appointment.startTime));
    
    if (!groups[date]) {
      groups[date] = [];
    }
    
    groups[date].push(appointment);
    return groups;
  }, {} as Record<string, T[]>);
}
