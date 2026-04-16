import { describe, it, expect } from 'vitest';
import { arnelaSiteImages } from '../arnelaSiteAssets';

describe('arnelaSiteAssets', () => {
  it('serves all marketing images from /images (no remote URLs)', () => {
    for (const [key, path] of Object.entries(arnelaSiteImages)) {
      expect(path, key).toMatch(/^\/images\//);
      expect(path, key).not.toMatch(/^https?:\/\//);
    }
  });

  it('includes expected stable keys for public pages', () => {
    expect(Object.keys(arnelaSiteImages).sort()).toEqual(
      [
        'acompanamientoAdultos',
        'charlasCentros',
        'conveniosHero',
        'formacionProfesionales',
        'fotoPrincipal',
        'infanciaAdolescencia',
        'intervencionFamiliar',
        'logo',
        'talleresFamilias',
        'teamBuilding',
      ].sort(),
    );
  });
});
