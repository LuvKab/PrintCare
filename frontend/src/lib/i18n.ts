import { get } from 'svelte/store';
import { translations, currentLang } from './stores';
import { GetLocale } from '../../wailsjs/go/main/App';

export async function loadLocale(lang: string): Promise<void> {
  try {
    const locale = await GetLocale(lang);
    const flat: Record<string, string> = {};
    for (const [k, v] of Object.entries(locale)) {
      if (typeof v === 'string') {
        flat[k] = v;
      }
    }
    translations.set(flat);
    currentLang.set(lang);
  } catch (e) {
    console.error('Failed to load locale:', e);
  }
}

export async function switchLocale(lang: string): Promise<void> {
  await loadLocale(lang);
}
