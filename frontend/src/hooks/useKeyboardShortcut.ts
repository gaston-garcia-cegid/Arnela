import { useEffect } from 'react';

export function useKeyboardShortcut(
  key: string,
  callback: () => void,
  options: {
    ctrl?: boolean;
    alt?: boolean;
    shift?: boolean;
    meta?: boolean;
  } = {}
) {
  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      const { ctrl = false, alt = false, shift = false, meta = false } = options;

      const ctrlMatch = ctrl ? event.ctrlKey || event.metaKey : true;
      const altMatch = alt ? event.altKey : true;
      const shiftMatch = shift ? event.shiftKey : true;
      const metaMatch = meta ? event.metaKey : true;

      if (
        event.key.toLowerCase() === key.toLowerCase() &&
        ctrlMatch &&
        altMatch &&
        shiftMatch &&
        metaMatch &&
        !event.repeat
      ) {
        event.preventDefault();
        callback();
      }
    };

    window.addEventListener('keydown', handleKeyDown);

    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [key, callback, options]);
}
