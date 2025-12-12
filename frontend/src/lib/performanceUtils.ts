/**
 * Performance Utilities
 * 
 * Collection of optimized helper functions using efficient data structures
 * to reduce computational complexity from O(n²) to O(n) or O(n log n).
 * 
 * @module performanceUtils
 */

/**
 * Count items by status in O(n) instead of O(4n)
 * 
 * BEFORE (O(4n) - 4 iterations):
 * ```ts
 * const pending = items.filter(i => i.status === 'pending').length;
 * const confirmed = items.filter(i => i.status === 'confirmed').length;
 * const completed = items.filter(i => i.status === 'completed').length;
 * const cancelled = items.filter(i => i.status === 'cancelled').length;
 * ```
 * 
 * AFTER (O(n) - single iteration):
 * ```ts
 * const stats = countByStatus(items, ['pending', 'confirmed', 'completed', 'cancelled']);
 * ```
 * 
 * @param items Array of items with status property
 * @param statusList Array of status values to count
 * @returns Object with status counts
 * 
 * @example
 * const appointments = [
 *   { id: '1', status: 'pending' },
 *   { id: '2', status: 'confirmed' },
 *   { id: '3', status: 'pending' },
 * ];
 * const stats = countByStatus(appointments, ['pending', 'confirmed', 'cancelled']);
 * // { pending: 2, confirmed: 1, cancelled: 0 }
 */
export function countByStatus<T extends { status: string }>(
  items: T[],
  statusList: string[]
): Record<string, number> {
  // Initialize all counts to 0
  const counts: Record<string, number> = {};
  statusList.forEach(status => counts[status] = 0);
  
  // Single pass O(n)
  for (const item of items) {
    if (item.status in counts) {
      counts[item.status]++;
    }
  }
  
  return counts;
}

/**
 * Create lookup Map for O(1) searches instead of O(n)
 * 
 * BEFORE (O(n) per lookup):
 * ```ts
 * const therapist = therapists.find(t => t.id === id)?.name;
 * ```
 * 
 * AFTER (O(1) per lookup after O(n) initial build):
 * ```ts
 * const therapistMap = createLookupMap(therapists, 'id');
 * const therapist = therapistMap.get(id);
 * ```
 * 
 * @param items Array of items
 * @param keyField Field to use as Map key
 * @returns Map for O(1) lookups
 * 
 * @example
 * const users = [
 *   { id: '1', name: 'Alice' },
 *   { id: '2', name: 'Bob' }
 * ];
 * const userMap = createLookupMap(users, 'id');
 * const alice = userMap.get('1'); // O(1) lookup
 */
export function createLookupMap<T, K extends keyof T>(
  items: T[],
  keyField: K
): Map<T[K], T> {
  const map = new Map<T[K], T>();
  
  for (const item of items) {
    map.set(item[keyField], item);
  }
  
  return map;
}

/**
 * Extract unique values efficiently using Set O(n)
 * 
 * BEFORE (O(n²) - Array.from + filter uses Set internally but verbose):
 * ```ts
 * const uniqueCities = Array.from(
 *   new Set(clients.map(c => c.city).filter(Boolean))
 * ).sort();
 * ```
 * 
 * AFTER (O(n + m log m) where m is unique count):
 * ```ts
 * const uniqueCities = extractUniqueSorted(clients, 'city');
 * ```
 * 
 * @param items Array of items
 * @param field Field to extract unique values from
 * @param filterEmpty Remove null/undefined/empty strings (default: true)
 * @returns Sorted array of unique values
 * 
 * @example
 * const clients = [
 *   { city: 'Madrid' },
 *   { city: 'Barcelona' },
 *   { city: 'Madrid' },
 *   { city: null }
 * ];
 * const cities = extractUniqueSorted(clients, 'city');
 * // ['Barcelona', 'Madrid']
 */
export function extractUniqueSorted<T, K extends keyof T>(
  items: T[],
  field: K,
  filterEmpty = true
): Array<Exclude<T[K], null | undefined | ''>> {
  const uniqueSet = new Set<T[K]>();
  
  for (const item of items) {
    const value = item[field];
    
    if (filterEmpty) {
      // Filter out falsy values
      if (value) {
        uniqueSet.add(value);
      }
    } else {
      uniqueSet.add(value);
    }
  }
  
  // Convert to array and sort
  const uniqueArray = Array.from(uniqueSet) as Array<Exclude<T[K], null | undefined | ''>>;
  
  return uniqueArray.sort((a, b) => {
    // Type-safe comparison
    if (typeof a === 'string' && typeof b === 'string') {
      return a.localeCompare(b);
    }
    return 0;
  });
}

/**
 * Create index mapping for fast lookups O(1)
 * Useful when you need to find items by ID repeatedly
 * 
 * BEFORE (O(n) per .find() call):
 * ```ts
 * const categoryName = categories.find(c => c.id === expense.categoryId)?.name;
 * ```
 * 
 * AFTER (O(1) after initial O(n) build):
 * ```ts
 * const categoryMap = createPropertyMap(categories, 'id', 'name');
 * const categoryName = categoryMap.get(expense.categoryId);
 * ```
 * 
 * @param items Array of items
 * @param keyField Field to use as key
 * @param valueField Field to extract as value
 * @returns Map for O(1) lookups
 * 
 * @example
 * const categories = [
 *   { id: '1', name: 'Office Supplies' },
 *   { id: '2', name: 'Travel' }
 * ];
 * const nameMap = createPropertyMap(categories, 'id', 'name');
 * const name = nameMap.get('1'); // 'Office Supplies' in O(1)
 */
export function createPropertyMap<T, K extends keyof T, V extends keyof T>(
  items: T[],
  keyField: K,
  valueField: V
): Map<T[K], T[V]> {
  const map = new Map<T[K], T[V]>();
  
  for (const item of items) {
    map.set(item[keyField], item[valueField]);
  }
  
  return map;
}

/**
 * Batch filter by multiple conditions in O(n) instead of O(kn)
 * 
 * BEFORE (Multiple passes):
 * ```ts
 * let filtered = items;
 * if (status !== 'all') filtered = filtered.filter(i => i.status === status);
 * if (category) filtered = filtered.filter(i => i.category === category);
 * if (minAmount) filtered = filtered.filter(i => i.amount >= minAmount);
 * ```
 * 
 * AFTER (Single pass):
 * ```ts
 * const filtered = batchFilter(items, [
 *   { field: 'status', value: status, condition: 'equals' },
 *   { field: 'category', value: category, condition: 'equals' },
 *   { field: 'amount', value: minAmount, condition: 'gte' }
 * ]);
 * ```
 * 
 * @param items Array to filter
 * @param filters Array of filter conditions
 * @returns Filtered array
 */
export function batchFilter<T>(
  items: T[],
  filters: Array<{
    field: keyof T;
    value: any;
    condition?: 'equals' | 'not-equals' | 'gt' | 'gte' | 'lt' | 'lte' | 'includes';
  }>
): T[] {
  // Skip if no filters
  if (filters.length === 0) return items;
  
  // Filter out empty/undefined values
  const activeFilters = filters.filter(f => f.value !== undefined && f.value !== null && f.value !== '');
  if (activeFilters.length === 0) return items;
  
  // Single pass through items
  return items.filter(item => {
    // Check all conditions (AND logic)
    return activeFilters.every(filter => {
      const itemValue = item[filter.field];
      const filterValue = filter.value;
      
      switch (filter.condition || 'equals') {
        case 'equals':
          return itemValue === filterValue;
        case 'not-equals':
          return itemValue !== filterValue;
        case 'gt':
          return (itemValue as any) > filterValue;
        case 'gte':
          return (itemValue as any) >= filterValue;
        case 'lt':
          return (itemValue as any) < filterValue;
        case 'lte':
          return (itemValue as any) <= filterValue;
        case 'includes':
          return String(itemValue).toLowerCase().includes(String(filterValue).toLowerCase());
        default:
          return true;
      }
    });
  });
}

/**
 * Memoize expensive computations
 * 
 * @param fn Function to memoize
 * @param keyGenerator Optional custom key generator
 * @returns Memoized function
 * 
 * @example
 * const expensiveCalc = memoize((n: number) => {
 *   // Heavy computation
 *   return n ** 2;
 * });
 * 
 * expensiveCalc(5); // Computed
 * expensiveCalc(5); // Cached O(1)
 */
export function memoize<Args extends any[], Return>(
  fn: (...args: Args) => Return,
  keyGenerator?: (...args: Args) => string
): (...args: Args) => Return {
  const cache = new Map<string, Return>();
  
  return (...args: Args): Return => {
    const key = keyGenerator ? keyGenerator(...args) : JSON.stringify(args);
    
    if (cache.has(key)) {
      return cache.get(key)!;
    }
    
    const result = fn(...args);
    cache.set(key, result);
    
    return result;
  };
}

/**
 * WeakMap-based cache for object-keyed memoization
 * Automatically garbage-collected when objects are no longer referenced
 * 
 * @example
 * const cache = new WeakMapCache<User, UserStats>();
 * 
 * function getStats(user: User) {
 *   return cache.getOrCompute(user, () => calculateStats(user));
 * }
 */
export class WeakMapCache<K extends object, V> {
  private cache = new WeakMap<K, V>();
  
  /**
   * Get cached value or compute and cache it
   */
  getOrCompute(key: K, compute: () => V): V {
    if (this.cache.has(key)) {
      return this.cache.get(key)!;
    }
    
    const value = compute();
    this.cache.set(key, value);
    
    return value;
  }
  
  /**
   * Check if key exists in cache
   */
  has(key: K): boolean {
    return this.cache.has(key);
  }
  
  /**
   * Get cached value
   */
  get(key: K): V | undefined {
    return this.cache.get(key);
  }
  
  /**
   * Set value in cache
   */
  set(key: K, value: V): void {
    this.cache.set(key, value);
  }
}
