/**
 * Centralized Logging Service
 * Development: console. Production: optional Sentry (browser) when NEXT_PUBLIC_SENTRY_DSN is set.
 */

type LogLevel = 'info' | 'warn' | 'error' | 'debug';

interface LogContext {
    component?: string;
    action?: string;
    userId?: string;
    [key: string]: any;
}

class Logger {
    private isDevelopment: boolean;

    constructor() {
        this.isDevelopment = process.env.NODE_ENV === 'development';
    }

    /**
     * Log information message (development only)
     */
    info(message: string, context?: LogContext) {
        if (this.isDevelopment) {
            console.log(`[INFO] ${message}`, context || '');
        }
    }

    /**
     * Log warning message (development only) 
     */
    warn(message: string, context?: LogContext) {
        if (this.isDevelopment) {
            console.warn(`[WARN] ${message}`, context || '');
        }
    }

    /**
     * Log error message
     * In development: console.error
     * In production: Sentry (client) when NEXT_PUBLIC_SENTRY_DSN is configured
     */
    error(message: string, error?: Error | unknown, context?: LogContext) {
        if (this.isDevelopment) {
            console.error(`[ERROR] ${message}`, {
                error,
                context,
                stack: error instanceof Error ? error.stack : undefined,
            });
            return;
        }
        if (!process.env.NEXT_PUBLIC_SENTRY_DSN || typeof window === 'undefined') {
            return;
        }
        const err = error instanceof Error ? error : new Error(message);
        void import('@sentry/browser').then((Sentry) => {
            Sentry.captureException(err, {
                tags: { source: 'logger' },
                extra: { logMessage: message, ...context },
            });
        });
    }

    /**
     * Log debug message (development only)
     */
    debug(message: string, data?: any) {
        if (this.isDevelopment) {
            console.debug(`[DEBUG] ${message}`, data || '');
        }
    }

    /**
     * Log API request (development only)
     */
    apiRequest(method: string, url: string, data?: any) {
        if (this.isDevelopment) {
            console.log(`[API ${method}] ${url}`, data || '');
        }
    }

    /**
     * Log API response (development only)
     */
    apiResponse(method: string, url: string, status: number, data?: any) {
        if (this.isDevelopment) {
            const level = status >= 400 ? 'error' : 'log';
            console[level](`[API ${method}] ${url} - ${status}`, data || '');
        }
    }
}

// Export singleton instance
export const logger = new Logger();

// Export convenience functions
export const logInfo = (message: string, context?: LogContext) => logger.info(message, context);
export const logWarn = (message: string, context?: LogContext) => logger.warn(message, context);
export const logError = (message: string, error?: Error | unknown, context?: LogContext) =>
    logger.error(message, error, context);
export const logDebug = (message: string, data?: any) => logger.debug(message, data);
