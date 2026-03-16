import { writable, derived } from 'svelte/store';

export type Page = 'home' | 'settings' | 'print' | 'logs';

export const currentPage = writable<Page>('home');
export const darkMode = writable<boolean>(false);
export const translations = writable<Record<string, string>>({});
export const currentLang = writable<string>('en');

export const t = derived(translations, ($tr) => {
  return (key: string): string => {
    return $tr[key] || key;
  };
});
