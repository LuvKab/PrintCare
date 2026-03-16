<script lang="ts">
  import { onMount } from 'svelte';
  import { t, currentLang } from '../lib/stores';
  import { switchLocale } from '../lib/i18n';
  import { GetConfig, SaveConfig, GetPrinters, SelectImage, GetAvailableLocales } from '../../wailsjs/go/main/App';

  let enabled = false;
  let printerName = '';
  let intervalDays = 7;
  let imagePath = '';
  let language = 'en';
  let printers: string[] = [];
  let locales: { code: string; name: string }[] = [];
  let saving = false;
  let message = '';
  let messageType: 'success' | 'error' = 'success';

  onMount(async () => {
    try {
      const [cfg, printerList, localeList] = await Promise.all([
        GetConfig(),
        GetPrinters(),
        GetAvailableLocales(),
      ]);
      enabled = cfg.enabled;
      printerName = cfg.printerName;
      intervalDays = cfg.intervalDays;
      imagePath = cfg.imagePath;
      language = cfg.language;
      printers = printerList;
      locales = localeList;
    } catch (e) {
      console.error('Failed to load settings:', e);
    }
  });

  async function refreshPrinters() {
    printers = await GetPrinters();
  }

  async function selectImage() {
    try {
      const path = await SelectImage();
      if (path) {
        imagePath = path;
      }
    } catch (e) {
      console.error('SelectImage failed:', e);
    }
  }

  async function save() {
    saving = true;
    message = '';
    try {
      await SaveConfig({
        enabled,
        printerName,
        paperSource: 0,
        intervalDays,
        imagePath,
        language,
      });
      if (language !== $currentLang) {
        await switchLocale(language);
      }
      message = $t('settings.saved');
      messageType = 'success';
    } catch (e) {
      message = $t('settings.saveFailed');
      messageType = 'error';
    }
    saving = false;
    setTimeout(() => { message = ''; }, 3000);
  }
</script>

<div>
  <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-6">{$t('settings.title')}</h1>

  <div class="bg-white dark:bg-slate-800 rounded-xl shadow-sm border border-gray-100 dark:border-slate-700 divide-y divide-gray-100 dark:divide-slate-700">

    <!-- Enable toggle -->
    <div class="flex items-center justify-between px-6 py-5">
      <div>
        <label class="text-sm font-medium text-gray-900 dark:text-gray-100">{$t('settings.enable')}</label>
      </div>
      <button
        class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors
               {enabled ? 'bg-primary-500' : 'bg-gray-300 dark:bg-slate-600'}"
        on:click={() => enabled = !enabled}
      >
        <span class="inline-block h-4 w-4 rounded-full bg-white transition-transform shadow-sm
                     {enabled ? 'translate-x-6' : 'translate-x-1'}"></span>
      </button>
    </div>

    <!-- Printer -->
    <div class="px-6 py-5">
      <label class="block text-sm font-medium text-gray-900 dark:text-gray-100 mb-2">{$t('settings.printer')}</label>
      <div class="flex gap-2">
        <select
          bind:value={printerName}
          class="flex-1 rounded-lg border border-gray-300 dark:border-slate-600 bg-white dark:bg-slate-700
                 text-sm text-gray-700 dark:text-gray-200 px-3 py-2 outline-none
                 focus:ring-2 focus:ring-primary-500 focus:border-primary-500 transition"
        >
          <option value="">{$t('settings.selectPrinter')}</option>
          {#each printers as p}
            <option value={p}>{p}</option>
          {/each}
        </select>
        <button
          on:click={refreshPrinters}
          class="px-3 py-2 rounded-lg border border-gray-300 dark:border-slate-600 text-gray-600 dark:text-gray-300
                 hover:bg-gray-50 dark:hover:bg-slate-700 transition"
          title="Refresh"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Interval -->
    <div class="px-6 py-5">
      <label class="block text-sm font-medium text-gray-900 dark:text-gray-100 mb-2">{$t('settings.interval')}</label>
      <input
        type="number"
        bind:value={intervalDays}
        min="1"
        max="365"
        class="w-32 rounded-lg border border-gray-300 dark:border-slate-600 bg-white dark:bg-slate-700
               text-sm text-gray-700 dark:text-gray-200 px-3 py-2 outline-none
               focus:ring-2 focus:ring-primary-500 focus:border-primary-500 transition"
      />
    </div>

    <!-- Image -->
    <div class="px-6 py-5">
      <label class="block text-sm font-medium text-gray-900 dark:text-gray-100 mb-2">{$t('settings.image')}</label>
      <div class="flex items-center gap-3">
        <button
          on:click={selectImage}
          class="px-4 py-2 rounded-lg bg-primary-500 text-white text-sm font-medium
                 hover:bg-primary-600 active:bg-primary-700 transition-colors"
        >
          {$t('settings.selectImage')}
        </button>
        <span class="text-sm text-gray-500 dark:text-gray-400 truncate max-w-xs">
          {imagePath || $t('settings.noImage')}
        </span>
      </div>
    </div>

    <!-- Language -->
    <div class="px-6 py-5">
      <label class="block text-sm font-medium text-gray-900 dark:text-gray-100 mb-2">{$t('settings.language')}</label>
      <select
        bind:value={language}
        class="w-48 rounded-lg border border-gray-300 dark:border-slate-600 bg-white dark:bg-slate-700
               text-sm text-gray-700 dark:text-gray-200 px-3 py-2 outline-none
               focus:ring-2 focus:ring-primary-500 focus:border-primary-500 transition"
      >
        {#each locales as loc}
          <option value={loc.code}>{loc.name}</option>
        {/each}
      </select>
    </div>
  </div>

  <!-- Save -->
  <div class="mt-6 flex items-center gap-4">
    <button
      on:click={save}
      disabled={saving}
      class="px-6 py-2.5 rounded-lg bg-primary-500 text-white text-sm font-medium
             hover:bg-primary-600 active:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed
             transition-colors shadow-sm"
    >
      {saving ? '...' : $t('settings.save')}
    </button>
    {#if message}
      <span class="text-sm font-medium {messageType === 'success' ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'}">
        {message}
      </span>
    {/if}
  </div>
</div>
