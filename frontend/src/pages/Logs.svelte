<script lang="ts">
  import { onMount } from 'svelte';
  import { t } from '../lib/stores';
  import { GetLogs, ClearLogs } from '../../wailsjs/go/main/App';

  let logs = '';
  let loading = false;

  async function refresh() {
    loading = true;
    try {
      logs = await GetLogs();
    } catch (e) {
      console.error('Failed to get logs:', e);
    }
    loading = false;
  }

  async function clear() {
    try {
      await ClearLogs();
      logs = '';
    } catch (e) {
      console.error('Failed to clear logs:', e);
    }
  }

  onMount(refresh);
</script>

<div>
  <div class="flex items-center justify-between mb-6">
    <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">{$t('logs.title')}</h1>
    <div class="flex gap-2">
      <button
        on:click={refresh}
        disabled={loading}
        class="px-4 py-2 rounded-lg border border-gray-300 dark:border-slate-600 text-sm font-medium
               text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-slate-700
               disabled:opacity-50 transition-colors"
      >
        {$t('logs.refresh')}
      </button>
      <button
        on:click={clear}
        class="px-4 py-2 rounded-lg border border-red-300 dark:border-red-800 text-sm font-medium
               text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
      >
        {$t('logs.clear')}
      </button>
    </div>
  </div>

  <div class="bg-white dark:bg-slate-800 rounded-xl shadow-sm border border-gray-100 dark:border-slate-700 overflow-hidden">
    {#if logs}
      <pre class="p-4 text-xs font-mono text-gray-700 dark:text-gray-300 whitespace-pre-wrap break-words
                  max-h-[calc(100vh-220px)] overflow-y-auto leading-relaxed">{logs}</pre>
    {:else}
      <div class="p-12 text-center">
        <svg class="w-12 h-12 mx-auto text-gray-300 dark:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <p class="mt-3 text-sm text-gray-400 dark:text-gray-500">{$t('logs.empty')}</p>
      </div>
    {/if}
  </div>
</div>
