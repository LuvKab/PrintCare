<script lang="ts">
  import { onMount } from 'svelte';
  import { t } from '../lib/stores';
  import { GetStatus, GetConfig, PrintNow as DoPrint } from '../../wailsjs/go/main/App';

  let printerName = '';
  let imagePath = '';
  let printing = false;
  let message = '';
  let messageType: 'success' | 'error' = 'success';

  onMount(async () => {
    try {
      const cfg = await GetConfig();
      printerName = cfg.printerName;
      imagePath = cfg.imagePath;
    } catch (e) {
      console.error('Failed to load config:', e);
    }
  });

  async function handlePrint() {
    printing = true;
    message = '';
    try {
      await DoPrint();
      message = $t('print.success');
      messageType = 'success';
    } catch (e: any) {
      message = $t('print.failed') + ': ' + (e?.message || e);
      messageType = 'error';
    }
    printing = false;
  }

  $: fileName = imagePath ? imagePath.split(/[/\\]/).pop() : '';
</script>

<div>
  <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-6">{$t('print.title')}</h1>

  <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
    <!-- Info Panel -->
    <div class="bg-white dark:bg-slate-800 rounded-xl p-6 shadow-sm border border-gray-100 dark:border-slate-700">
      <div class="space-y-4">
        <div>
          <label class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide">{$t('print.printer')}</label>
          <p class="mt-1 text-sm font-medium text-gray-900 dark:text-gray-100">
            {printerName || '—'}
          </p>
        </div>
        <div>
          <label class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide">{$t('print.image')}</label>
          <p class="mt-1 text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
            {fileName || '—'}
          </p>
        </div>
      </div>

      <div class="mt-8">
        <button
          on:click={handlePrint}
          disabled={printing || !printerName || !imagePath}
          class="w-full py-3 rounded-lg text-white text-sm font-semibold transition-all shadow-sm
                 {printing
                   ? 'bg-gray-400 cursor-not-allowed'
                   : 'bg-primary-500 hover:bg-primary-600 active:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed'}"
        >
          {#if printing}
            <span class="inline-flex items-center gap-2">
              <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
              </svg>
              {$t('print.printing')}
            </span>
          {:else}
            {$t('print.btn')}
          {/if}
        </button>
      </div>

      {#if message}
        <div class="mt-4 p-3 rounded-lg text-sm font-medium
                    {messageType === 'success'
                      ? 'bg-green-50 text-green-700 dark:bg-green-900/20 dark:text-green-400'
                      : 'bg-red-50 text-red-700 dark:bg-red-900/20 dark:text-red-400'}">
          {message}
        </div>
      {/if}
    </div>

    <!-- Preview Panel -->
    <div class="bg-white dark:bg-slate-800 rounded-xl p-6 shadow-sm border border-gray-100 dark:border-slate-700">
      <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-4">{$t('print.preview')}</h3>
      {#if imagePath}
        <div class="rounded-lg bg-gray-100 dark:bg-slate-700 p-4 flex items-center justify-center min-h-[200px]">
          <div class="text-center">
            <svg class="w-12 h-12 mx-auto text-gray-400 dark:text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                    d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">{fileName}</p>
          </div>
        </div>
      {:else}
        <div class="rounded-lg bg-gray-50 dark:bg-slate-700/50 p-8 flex items-center justify-center min-h-[200px]">
          <p class="text-sm text-gray-400 dark:text-gray-500">{$t('print.noPreview')}</p>
        </div>
      {/if}
    </div>
  </div>
</div>
