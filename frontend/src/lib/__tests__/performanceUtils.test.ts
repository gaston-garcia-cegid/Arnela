import { describe, it, expect } from 'vitest';
import {
  countByStatus,
  createLookupMap,
  extractUniqueSorted,
  createPropertyMap,
  batchFilter,
  memoize,
  WeakMapCache
} from '../performanceUtils';

describe('performanceUtils', () => {
  describe('countByStatus', () => {
    it('counts items by status in single pass', () => {
      const items = [
        { id: '1', status: 'pending' },
        { id: '2', status: 'confirmed' },
        { id: '3', status: 'pending' },
        { id: '4', status: 'cancelled' },
        { id: '5', status: 'pending' },
      ];

      const result = countByStatus(items, ['pending', 'confirmed', 'cancelled', 'completed']);

      expect(result).toEqual({
        pending: 3,
        confirmed: 1,
        cancelled: 1,
        completed: 0,
      });
    });

    it('handles empty array', () => {
      const result = countByStatus([], ['pending', 'confirmed']);
      
      expect(result).toEqual({
        pending: 0,
        confirmed: 0,
      });
    });

    it('ignores unknown statuses', () => {
      const items = [
        { id: '1', status: 'pending' },
        { id: '2', status: 'unknown' as any },
      ];

      const result = countByStatus(items, ['pending', 'confirmed']);

      expect(result).toEqual({
        pending: 1,
        confirmed: 0,
      });
    });
  });

  describe('createLookupMap', () => {
    it('creates Map for O(1) lookups', () => {
      const items = [
        { id: '1', name: 'Alice' },
        { id: '2', name: 'Bob' },
        { id: '3', name: 'Charlie' },
      ];

      const map = createLookupMap(items, 'id');

      expect(map.get('1')).toEqual({ id: '1', name: 'Alice' });
      expect(map.get('2')).toEqual({ id: '2', name: 'Bob' });
      expect(map.get('3')).toEqual({ id: '3', name: 'Charlie' });
      expect(map.get('4')).toBeUndefined();
    });

    it('handles duplicate keys by keeping last', () => {
      const items = [
        { id: '1', name: 'Alice' },
        { id: '1', name: 'Alice Updated' },
      ];

      const map = createLookupMap(items, 'id');

      expect(map.get('1')?.name).toBe('Alice Updated');
    });
  });

  describe('extractUniqueSorted', () => {
    it('extracts unique values and sorts them', () => {
      const items = [
        { city: 'Madrid' },
        { city: 'Barcelona' },
        { city: 'Madrid' },
        { city: 'Zaragoza' },
        { city: 'Barcelona' },
      ];

      const result = extractUniqueSorted(items, 'city');

      expect(result).toEqual(['Barcelona', 'Madrid', 'Zaragoza']);
    });

    it('filters out null/undefined by default', () => {
      const items = [
        { city: 'Madrid' },
        { city: null },
        { city: undefined },
        { city: '' },
        { city: 'Barcelona' },
      ];

      const result = extractUniqueSorted(items, 'city' as any);

      expect(result).toEqual(['Barcelona', 'Madrid']);
    });

    it('keeps null/undefined when filterEmpty is false', () => {
      const items = [
        { status: 'active' },
        { status: null },
        { status: 'inactive' },
      ];

      const result = extractUniqueSorted(items, 'status' as any, false);

      expect(result).toContain('active');
      expect(result).toContain('inactive');
      expect(result).toContain(null);
    });
  });

  describe('createPropertyMap', () => {
    it('maps key field to value field', () => {
      const categories = [
        { id: '1', name: 'Office Supplies', code: 'OS' },
        { id: '2', name: 'Travel', code: 'TR' },
        { id: '3', name: 'Equipment', code: 'EQ' },
      ];

      const nameMap = createPropertyMap(categories, 'id', 'name');

      expect(nameMap.get('1')).toBe('Office Supplies');
      expect(nameMap.get('2')).toBe('Travel');
      expect(nameMap.get('3')).toBe('Equipment');
    });

    it('can map to any field', () => {
      const items = [
        { id: '1', code: 'A', value: 100 },
        { id: '2', code: 'B', value: 200 },
      ];

      const valueMap = createPropertyMap(items, 'id', 'value');

      expect(valueMap.get('1')).toBe(100);
      expect(valueMap.get('2')).toBe(200);
    });
  });

  describe('batchFilter', () => {
    const expenses = [
      { id: '1', category: 'Office', amount: 100, status: 'paid' },
      { id: '2', category: 'Travel', amount: 500, status: 'pending' },
      { id: '3', category: 'Office', amount: 50, status: 'paid' },
      { id: '4', category: 'Equipment', amount: 2000, status: 'paid' },
    ];

    it('filters by equals condition', () => {
      const result = batchFilter(expenses, [
        { field: 'category', value: 'Office', condition: 'equals' }
      ]);

      expect(result).toHaveLength(2);
      expect(result[0].id).toBe('1');
      expect(result[1].id).toBe('3');
    });

    it('filters by multiple conditions (AND)', () => {
      const result = batchFilter(expenses, [
        { field: 'category', value: 'Office', condition: 'equals' },
        { field: 'status', value: 'paid', condition: 'equals' }
      ]);

      expect(result).toHaveLength(2);
    });

    it('filters by gte condition', () => {
      const result = batchFilter(expenses, [
        { field: 'amount', value: 500, condition: 'gte' }
      ]);

      expect(result).toHaveLength(2);
      expect(result.map(e => e.id)).toEqual(['2', '4']);
    });

    it('filters by lt condition', () => {
      const result = batchFilter(expenses, [
        { field: 'amount', value: 100, condition: 'lt' }
      ]);

      expect(result).toHaveLength(1);
      expect(result[0].id).toBe('3');
    });

    it('skips empty filter values', () => {
      const result = batchFilter(expenses, [
        { field: 'category', value: undefined },
        { field: 'status', value: 'paid' }
      ]);

      // Only status filter should apply
      expect(result).toHaveLength(3);
    });

    it('returns all items when no active filters', () => {
      const result = batchFilter(expenses, []);

      expect(result).toHaveLength(4);
    });
  });

  describe('memoize', () => {
    it('caches function results', () => {
      let callCount = 0;
      const expensiveCalc = memoize((n: number) => {
        callCount++;
        return n ** 2;
      });

      expect(expensiveCalc(5)).toBe(25);
      expect(callCount).toBe(1);

      // Second call should use cache
      expect(expensiveCalc(5)).toBe(25);
      expect(callCount).toBe(1);

      // Different argument should compute
      expect(expensiveCalc(10)).toBe(100);
      expect(callCount).toBe(2);
    });

    it('works with custom key generator', () => {
      let callCount = 0;
      const fn = memoize(
        (user: { id: string; name: string }) => {
          callCount++;
          return user.name.toUpperCase();
        },
        (user) => user.id // Use only ID as cache key
      );

      const user1 = { id: '1', name: 'Alice' };
      const user1Updated = { id: '1', name: 'Alice Smith' };

      expect(fn(user1)).toBe('ALICE');
      expect(callCount).toBe(1);

      // Same ID should use cache even though name changed
      expect(fn(user1Updated)).toBe('ALICE');
      expect(callCount).toBe(1);
    });
  });

  describe('WeakMapCache', () => {
    it('caches values by object reference', () => {
      const cache = new WeakMapCache<object, string>();
      let computeCount = 0;

      const obj1 = { id: 1 };
      const obj2 = { id: 2 };

      const result1 = cache.getOrCompute(obj1, () => {
        computeCount++;
        return 'computed-1';
      });

      expect(result1).toBe('computed-1');
      expect(computeCount).toBe(1);

      // Second call should use cache
      const result2 = cache.getOrCompute(obj1, () => {
        computeCount++;
        return 'should-not-call';
      });

      expect(result2).toBe('computed-1');
      expect(computeCount).toBe(1);

      // Different object should compute
      const result3 = cache.getOrCompute(obj2, () => {
        computeCount++;
        return 'computed-2';
      });

      expect(result3).toBe('computed-2');
      expect(computeCount).toBe(2);
    });

    it('has/get/set methods work correctly', () => {
      const cache = new WeakMapCache<object, number>();
      const obj = { id: 1 };

      expect(cache.has(obj)).toBe(false);
      expect(cache.get(obj)).toBeUndefined();

      cache.set(obj, 42);

      expect(cache.has(obj)).toBe(true);
      expect(cache.get(obj)).toBe(42);
    });
  });

  // Performance comparison tests (informational)
  describe('Performance Benchmarks', () => {
    it('countByStatus is faster than multiple filters', () => {
      const items = Array.from({ length: 1000 }, (_, i) => ({
        id: `${i}`,
        status: ['pending', 'confirmed', 'cancelled', 'completed'][i % 4]
      }));

      // SLOW METHOD (4 passes - O(4n))
      const startSlow = performance.now();
      const pendingSlow = items.filter(i => i.status === 'pending').length;
      const confirmedSlow = items.filter(i => i.status === 'confirmed').length;
      const cancelledSlow = items.filter(i => i.status === 'cancelled').length;
      const completedSlow = items.filter(i => i.status === 'completed').length;
      const slowTime = performance.now() - startSlow;

      // FAST METHOD (1 pass - O(n))
      const startFast = performance.now();
      const stats = countByStatus(items, ['pending', 'confirmed', 'cancelled', 'completed']);
      const fastTime = performance.now() - startFast;

      // Results should match
      expect(stats.pending).toBe(pendingSlow);
      expect(stats.confirmed).toBe(confirmedSlow);
      expect(stats.cancelled).toBe(cancelledSlow);
      expect(stats.completed).toBe(completedSlow);

      // Fast method should be faster (usually 2-3x)
      console.log(`Performance: Slow=${slowTime.toFixed(2)}ms, Fast=${fastTime.toFixed(2)}ms, Speedup=${(slowTime/fastTime).toFixed(1)}x`);
    });

    it('Map lookup is faster than Array.find', () => {
      const items = Array.from({ length: 1000 }, (_, i) => ({
        id: `${i}`,
        name: `Item ${i}`
      }));

      const lookupIds = ['500', '750', '999', '250', '100'];

      // SLOW METHOD (O(n) per lookup)
      const startSlow = performance.now();
      const resultsSlow = lookupIds.map(id => items.find(item => item.id === id));
      const slowTime = performance.now() - startSlow;

      // FAST METHOD (O(1) per lookup after O(n) build)
      const startFast = performance.now();
      const map = createLookupMap(items, 'id');
      const resultsFast = lookupIds.map(id => map.get(id));
      const fastTime = performance.now() - startFast;

      // Results should match
      expect(resultsFast).toEqual(resultsSlow);

      console.log(`Map Lookup: Slow=${slowTime.toFixed(2)}ms, Fast=${fastTime.toFixed(2)}ms, Speedup=${(slowTime/fastTime).toFixed(1)}x`);
    });

    it('batchFilter is faster than multiple filter passes', () => {
      const items = Array.from({ length: 1000 }, (_, i) => ({
        id: `${i}`,
        category: ['A', 'B', 'C'][i % 3],
        amount: Math.random() * 1000,
        status: ['active', 'inactive'][i % 2]
      }));

      // SLOW METHOD (3 passes - O(3n))
      const startSlow = performance.now();
      let filteredSlow = items;
      filteredSlow = filteredSlow.filter(i => i.category === 'A');
      filteredSlow = filteredSlow.filter(i => i.amount >= 500);
      filteredSlow = filteredSlow.filter(i => i.status === 'active');
      const slowTime = performance.now() - startSlow;

      // FAST METHOD (1 pass - O(n))
      const startFast = performance.now();
      const filteredFast = batchFilter(items, [
        { field: 'category', value: 'A', condition: 'equals' },
        { field: 'amount', value: 500, condition: 'gte' },
        { field: 'status', value: 'active', condition: 'equals' }
      ]);
      const fastTime = performance.now() - startFast;

      // Results should match
      expect(filteredFast.length).toBe(filteredSlow.length);

      console.log(`Batch Filter: Slow=${slowTime.toFixed(2)}ms, Fast=${fastTime.toFixed(2)}ms, Speedup=${(slowTime/fastTime).toFixed(1)}x`);
    });
  });
});
