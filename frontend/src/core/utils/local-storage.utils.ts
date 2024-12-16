export function parseLocalStorageItem<T>(key: string): T | undefined {
  if (typeof window === 'undefined') {
    return undefined
  }

  try {
    const item = window.localStorage.getItem(key)
    return item ? (JSON.parse(item) as T) : undefined
  } catch (error) {
    console.error(
      `[useCurrentStoreStore] Failed to parse local storage item for key "${key}":`,
      error,
    )
    return undefined
  }
}
