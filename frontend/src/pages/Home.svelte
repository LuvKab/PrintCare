<script lang="ts">
  import { onMount } from 'svelte';
  import { t } from '../lib/stores';
  import StatusCard from '../components/StatusCard.svelte';
  import { GetStatus } from '../../wailsjs/go/main/App';

  let status = {
    enabled: false,
    printerReady: false,
    lastPrint: '',
    nextPrint: '',
    statusText: '',
    statusLevel: 'gray',
  };
  let interval: ReturnType<typeof setInterval>;

  async function refresh() {
    try {
      status = await GetStatus();
    } catch (e) {
      console.error('GetStatus failed:', e);
    }
  }

  onMount(() => {
    refresh();
    interval = setInterval(refresh, 15000);
    return () => clearInterval(interval);
  });

  $: statusColor = ({
    ok: 'green' as const,
    warning: 'yellow' as const,
    error: 'red' as const,
  })[status.statusLevel] || 'gray' as const;

  $: indicatorClasses = ({
    ok: 'bg-green-500',
    warning: 'bg-yellow-500',
    error: 'bg-red-500',
  })[status.statusLevel] || 'bg-gray-400';
</script>

<div>
  <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-6">{$t('home.title')}</h1>

  <!-- Status Banner -->
  <div class="bg-white dark:bg-slate-800 rounded-xl p-6 shadow-sm border border-gray-100 dark:border-slate-700 mb-6">
    <div class="flex items-center gap-4">
      <div class="relative">
        <div class="w-4 h-4 rounded-full {indicatorClasses}"></div>
        {#if status.statusLevel === 'ok'}
          <div class="absolute inset-0 w-4 h-4 rounded-full {indicatorClasses} animate-ping opacity-75"></div>
        {/if}
      </div>
      <div>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-gray-100">{$t('home.status')}</h2>
        <p class="text-sm text-gray-500 dark:text-gray-400">{$t('status.' + status.statusText) || status.statusText}</p>
      </div>
      <div class="ml-auto">
        <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium
                     {status.enabled
                       ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
                       : 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-400'}">
          {$t('home.enabled')}: {status.enabled ? $t('common.on') : $t('common.off')}
        </span>
      </div>
    </div>
  </div>

  <!-- Info Cards -->
  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <StatusCard
      title={$t('home.lastPrint')}
      value={status.lastPrint || $t('home.never')}
      icon="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
      color={status.lastPrint ? 'blue' : 'gray'}
    />
    <StatusCard
      title={$t('home.nextPrint')}
      value={status.nextPrint || $t('home.notScheduled')}
      icon="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
      color={status.nextPrint ? 'green' : 'gray'}
    />
  </div>
</div>
