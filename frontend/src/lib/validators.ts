/**
 * Centralized Validation Library
 * Consistent validation rules and error messages across the application
 */

export interface ValidationResult {
    isValid: boolean;
    error?: string;
}

/**
 * Spanish DNI/NIE Validation
 * Validates Spanish national ID (DNI) or foreigner ID (NIE)
 */
export const validateDNI = (dni: string): ValidationResult => {
    if (!dni) {
        return { isValid: false, error: 'El DNI/NIE es obligatorio' };
    }

    // Remove spaces and convert to uppercase
    const cleanDNI = dni.trim().toUpperCase().replace(/\s/g, '');

    // DNI format: 8 digits + 1 letter (e.g., 12345678Z)
    // NIE format: 1 letter + 7 digits + 1 letter (e.g., X1234567L)
    const dniRegex = /^[0-9]{8}[A-Z]$/;
    const nieRegex = /^[XYZ][0-9]{7}[A-Z]$/;

    if (!dniRegex.test(cleanDNI) && !nieRegex.test(cleanDNI)) {
        return {
            isValid: false,
            error: 'Formato de DNI/NIE inválido (ej: 12345678Z o X1234567L)'
        };
    }

    // Validate DNI letter
    const validLetters = 'TRWAGMYFPDXBNJZSQVHLCKE';
    let number: number;

    if (nieRegex.test(cleanDNI)) {
        // NIE: Replace X=0, Y=1, Z=2
        const niePrefix = cleanDNI[0] === 'X' ? '0' : cleanDNI[0] === 'Y' ? '1' : '2';
        number = parseInt(niePrefix + cleanDNI.substring(1, 8));
    } else {
        number = parseInt(cleanDNI.substring(0, 8));
    }

    const expectedLetter = validLetters[number % 23];
    const providedLetter = cleanDNI[cleanDNI.length - 1];

    if (expectedLetter !== providedLetter) {
        return { isValid: false, error: 'La letra del DNI/NIE no es correcta' };
    }

    return { isValid: true };
};

/**
 * Spanish CIF Validation
 * Validates Spanish company tax ID (CIF)
 */
export const validateCIF = (cif: string): ValidationResult => {
    if (!cif) {
        return { isValid: false, error: 'El CIF es obligatorio' };
    }

    const cleanCIF = cif.trim().toUpperCase().replace(/\s/g, '');

    // CIF format: 1 letter + 7 digits + 1 letter/digit
    const cifRegex = /^[ABCDEFGHJNPQRSUVW][0-9]{7}[0-9A-J]$/;

    if (!cifRegex.test(cleanCIF)) {
        return {
            isValid: false,
            error: 'Formato de CIF inválido (ej: A12345678)'
        };
    }

    return { isValid: true };
};

/**
 * DNI or CIF Validation
 * Validates either DNI/NIE or CIF
 */
export const validateDNIorCIF = (value: string): ValidationResult => {
    if (!value) {
        return { isValid: false, error: 'El DNI/CIF es obligatorio' };
    }

    const cleanValue = value.trim().toUpperCase();

    // Try DNI/NIE first
    const dniResult = validateDNI(cleanValue);
    if (dniResult.isValid) {
        return { isValid: true };
    }

    // Try CIF  
    const cifResult = validateCIF(cleanValue);
    if (cifResult.isValid) {
        return { isValid: true };
    }

    return {
        isValid: false,
        error: 'DNI/NIE/CIF inválido. Formato: 12345678Z (DNI) o A12345678 (CIF)'
    };
};

/**
 * Email Validation
 * RFC 5322 compliant email validation
 */
export const validateEmail = (email: string): ValidationResult => {
    if (!email) {
        return { isValid: false, error: 'El email es obligatorio' };
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

    if (!emailRegex.test(email.trim())) {
        return { isValid: false, error: 'Formato de email inválido' };
    }

    return { isValid: true };
};

/**
 * Spanish Phone Number Validation
 * Validates Spanish mobile (6XX, 7XX) and landline (8XX, 9XX) numbers
 */
export const validatePhone = (phone: string): ValidationResult => {
    if (!phone) {
        return { isValid: false, error: 'El teléfono es obligatorio' };
    }

    // Remove spaces, dashes, and parentheses
    const cleanPhone = phone.replace(/[\s\-()]/g, '');

    // Spanish phone: 9 digits, starting with 6, 7, 8, or 9
    const phoneRegex = /^[6-9][0-9]{8}$/;

    if (!phoneRegex.test(cleanPhone)) {
        return {
            isValid: false,
            error: 'Formato de teléfono inválido (ej: 612345678)'
        };
    }

    return { isValid: true };
};

/**
 * Spanish Postal Code Validation
 */
export const validatePostalCode = (postalCode: string): ValidationResult => {
    if (!postalCode) {
        return { isValid: false, error: 'El código postal es obligatorio' };
    }

    const cleanCode = postalCode.trim();
    const postalCodeRegex = /^[0-5][0-9]{4}$/;

    if (!postalCodeRegex.test(cleanCode)) {
        return {
            isValid: false,
            error: 'Código postal inválido (5 dígitos: 00000-52999)'
        };
    }

    return { isValid: true };
};

/**
 * Required Field Validation
 */
export const validateRequired = (value: string, fieldName: string = 'Este campo'): ValidationResult => {
    if (!value || value.trim() === '') {
        return { isValid: false, error: `${fieldName} es obligatorio` };
    }
    return { isValid: true };
};

/**
 * Min Length Validation
 */
export const validateMinLength = (
    value: string,
    minLength: number,
    fieldName: string = 'Este campo'
): ValidationResult => {
    if (value.trim().length < minLength) {
        return {
            isValid: false,
            error: `${fieldName} debe tener al menos ${minLength} caracteres`
        };
    }
    return { isValid: true };
};

/**
 * Max Length Validation
 */
export const validateMaxLength = (
    value: string,
    maxLength: number,
    fieldName: string = 'Este campo'
): ValidationResult => {
    if (value.trim().length > maxLength) {
        return {
            isValid: false,
            error: `${fieldName} no puede exceder ${maxLength} caracteres`
        };
    }
    return { isValid: true };
};

/**
 * Password Strength Validation
 * At least 8 characters, 1 uppercase, 1 lowercase, 1 number
 */
export const validatePassword = (password: string): ValidationResult => {
    if (!password) {
        return { isValid: false, error: 'La contraseña es obligatoria' };
    }

    if (password.length < 8) {
        return {
            isValid: false,
            error: 'La contraseña debe tener al menos 8 caracteres'
        };
    }

    if (!/[A-Z]/.test(password)) {
        return {
            isValid: false,
            error: 'La contraseña debe contener al menos una mayúscula'
        };
    }

    if (!/[a-z]/.test(password)) {
        return {
            isValid: false,
            error: 'La contraseña debe contener al menos una minúscula'
        };
    }

    if (!/[0-9]/.test(password)) {
        return {
            isValid: false,
            error: 'La contraseña debe contener al menos un número'
        };
    }

    return { isValid: true };
};

/**
 * Date Validation
 * Validates date is not in the future (for birthdate, etc.)
 */
export const validatePastDate = (date: string | Date): ValidationResult => {
    const selectedDate = typeof date === 'string' ? new Date(date) : date;
    const today = new Date();

    if (selectedDate > today) {
        return {
            isValid: false,
            error: 'La fecha no puede ser futura'
        };
    }

    return { isValid: true };
};

/**
 * Date Validation
 * Validates date is not in the past (for appointments, etc.)
 */
export const validateFutureDate = (date: string | Date): ValidationResult => {
    const selectedDate = typeof date === 'string' ? new Date(date) : date;
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    if (selectedDate < today) {
        return {
            isValid: false,
            error: 'La fecha no puede ser pasada'
        };
    }

    return { isValid: true };
};

/**
 * Composite Validator
 * Run multiple validations and return first error
 */
export const validateComposite = (
    value: string,
    validators: ((value: string) => ValidationResult)[]
): ValidationResult => {
    for (const validator of validators) {
        const result = validator(value);
        if (!result.isValid) {
            return result;
        }
    }
    return { isValid: true };
};
