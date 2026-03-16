<script lang="ts">
  import { onMount } from 'svelte';
  import { currentLang } from '../lib/stores';
  import { switchLocale } from '../lib/i18n';
  import { GetAvailableLocales } from '../../wailsjs/go/main/App';

  let locales: { code: string; name: string }[] = [];

  onMount(async () => {
    try {
      locales = await GetAvailableLocales();
    } catch (e) {
      console.error('Failed to get locales:', e);
    }
  });

  async function handleChange(e: Event) {
    const target = e.target as HTMLSelectElement;
    await switchLocale(target.value);
  }
</script>

<select
  value={$currentLang}
  on:change={handleChange}
  class="block w-full rounded-lg border border-gray-300 dark:border-slate-600 bg-white dark:bg-slate-700
         text-sm text-gray-700 dark:text-gray-200 px-3 py-2
         focus:ring-2 focus:ring-primary-500 focus:border-primary-500 outline-none transition"
>
  {#each locales as locale}
    <option value={locale.code}>{locale.name}</option>
  {/each}
</select>
