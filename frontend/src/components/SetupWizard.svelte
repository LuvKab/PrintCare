<script lang="ts">
  import { onMount } from 'svelte';
  import { t, currentLang } from '../lib/stores';
  import { switchLocale } from '../lib/i18n';
  import { GetPrinters, GetAvailableLocales, SaveConfig, SelectImage } from '../../wailsjs/go/main/App';

  export let onComplete: () => void;

  let language = 'en';
  let printerName = '';
  let imagePath = '';
  let intervalDays = 7;
  let enabled = true;

  let printers: string[] = [];
  let locales: { code: string; name: string }[] = [];
  let error = '';
  let saving = false;

  onMount(async () => {
    try {
      const [printerList, localeList] = await Promise.all([
        GetPrinters(),
        GetAvailableLocales(),
      ]);
      printers = printerList;
      locales = localeList;
      language = $currentLang;
    } catch (e) {
      console.error('Setup init failed:', e);
    }
  });

  async function handleLangChange(e: Event) {
    const target = e.target as HTMLSelectElement;
    language = target.value;
    await switchLocale(language);
  }

  async function refreshPrinters() {
    printers = await GetPrinters();
  }

  async function handleSelectImage() {
    try {
      const path = await SelectImage();
      if (path) imagePath = path;
    } catch (e) {
      console.error('SelectImage failed:', e);
    }
  }

  async function handleStart() {
    error = '';
    if (!printerName) {
      error = $t('setup.printerRequired');
      return;
    }
    if (!imagePath) {
      error = $t('setup.imageRequired');
      return;
    }

    saving = true;
    try {
      await SaveConfig({
        enabled,
        printerName,
        paperSource: 0,
        intervalDays,
        imagePath,
        language,
      });
      onComplete();
    } catch (e) {
      error = String(e);
    }
    saving = false;
  }

  $: fileName = imagePath ? imagePath.split(/[/\\]/).pop() : '';
</script>

<div class="flex items-center justify-center min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50 dark:from-slate-900 dark:via-slate-900 dark:to-slate-800 p-6">
  <div class="w-full max-w-lg">
    <!-- Header -->
    <div class="text-center mb-8">
      <div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-primary-500 shadow-lg shadow-primary-500/30 mb-4">
        <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4H7v4a2 2 0 002 2zm0-16h6a2 2 0 012 2v2H7V7a2 2 0 012-2z" />
        </svg>
      </div>
      <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">{$t('setup.title')}</h1>
      <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">{$t('setup.subtitle')}</p>
    </div>

    <!-- Setup Card -->
    <div class="bg-white dark:bg-slate-800 rounded-2xl shadow-xl border border-gray-100 dark:border-slate-700 overflow-hidden">
      <div class="divide-y divide-gray-100 dark:divide-slate-700">

        <!-- Language -->
        <div class="px-6 py-5">
          <label class="block text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-2">{$t('setup.language')}</label>
          <select
            value={language}
            on:change={handleLangChange}
            class="w-full rounded-lg border border-gray-300 dark:border-slate-600 bg-white dark:bg-slate-700
                   text-sm text-gray-700 dark:text-gray-200 px-3 py-2.5 outline-none
                   focus:ring-2 focus:ring-primary-500 focus:border-primary-500 transition"
          >
            {#each locales as loc}
              <option value={loc.code}>{loc.name}</option>
            {/each}
          </select>
        </div>

        <!-- Printer -->
        <div class="px-6 py-5">
          <label class="block text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-2">{$t('setup.printer')}</label>
          <div class="flex gap-2">
            <select
              bind:value={printerName}
              class="flex-1 rounded-lg border border-gray-300 dark:border-slate-600 bg-white dark:bg-slate-700
                     text-sm text-gray-700 dark:text-gray-200 px-3 py-2.5 outline-none
                     focus:ring-2 focus:ring-primary-500 focus:border-primary-500 transition"
            >
              <option value="">{$t('setup.selectPrinter')}</option>
              {#each printers as p}
                <option value={p}>{p}</option>
              {/each}
            </select>
            <button
              on:click={refreshPrinters}
              class="px-3 py-2.5 rounded-lg border border-gray-300 dark:border-slate-600 text-gray-500 dark:text-gray-400
                     hover:bg-gray-50 dark:hover:bg-slate-700 transition"
              title={$t('setup.refreshPrinters')}
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Image -->
        <div class="px-6 py-5">
          <label class="block text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-2">{$t('setup.image')}</label>
          <div class="flex items-center gap-3">
            <button
              on:click={handleSelectImage}
              class="px-4 py-2.5 rounded-lg bg-gray-100 dark:bg-slate-700 text-sm font-medium text-gray-700 dark:text-gray-200
                     hover:bg-gray-200 dark:hover:bg-slate-600 transition-colors"
            >
              {$t('setup.selectImage')}
            </button>
            <span class="text-sm text-gray-500 dark:text-gray-400 truncate">{fileName || ''}</span>
          </div>
        </div>

        <!-- Interval -->
        <div class="px-6 py-5">
          <label class="block text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-2">{$t('setup.interval')}</label>
          <input
            type="number"
            bind:value={intervalDays}
            min="1"
            max="365"
            class="w-28 rounded-lg border border-gray-300 dark:border-slate-600 bg-white dark:bg-slate-700
                   text-sm text-gray-700 dark:text-gray-200 px-3 py-2.5 outline-none
                   focus:ring-2 focus:ring-primary-500 focus:border-primary-500 transition"
          />
        </div>

        <!-- Enable -->
        <div class="flex items-center justify-between px-6 py-5">
          <span class="text-sm font-medium text-gray-700 dark:text-gray-200">{$t('setup.enable')}</span>
          <button
            class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors
                   {enabled ? 'bg-primary-500' : 'bg-gray-300 dark:bg-slate-600'}"
            on:click={() => enabled = !enabled}
          >
            <span class="inline-block h-4 w-4 rounded-full bg-white transition-transform shadow-sm
                         {enabled ? 'translate-x-6' : 'translate-x-1'}"></span>
          </button>
        </div>
      </div>

      <!-- Error -->
      {#if error}
        <div class="mx-6 mb-4 p-3 rounded-lg bg-red-50 dark:bg-red-900/20 text-sm text-red-600 dark:text-red-400">
          {error}
        </div>
      {/if}

      <!-- Start Button -->
      <div class="px-6 pb-6">
        <button
          on:click={handleStart}
          disabled={saving}
          class="w-full py-3 rounded-xl bg-primary-500 text-white text-sm font-semibold
                 hover:bg-primary-600 active:bg-primary-700 disabled:opacity-50
                 transition-all shadow-md shadow-primary-500/25"
        >
          {saving ? '...' : $t('setup.start')}
        </button>
      </div>
    </div>
  </div>
</div>
