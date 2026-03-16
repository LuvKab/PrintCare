<script lang="ts">
  import { onMount } from 'svelte';
  import Sidebar from './components/Sidebar.svelte';
  import SetupWizard from './components/SetupWizard.svelte';
  import Home from './pages/Home.svelte';
  import Settings from './pages/Settings.svelte';
  import PrintNow from './pages/PrintNow.svelte';
  import Logs from './pages/Logs.svelte';
  import { currentPage, darkMode } from './lib/stores';
  import { loadLocale } from './lib/i18n';
  import { GetConfig, IsFirstRun } from '../wailsjs/go/main/App';

  let ready = false;
  let showSetup = false;

  onMount(async () => {
    try {
      const cfg = await GetConfig();
      await loadLocale(cfg.language || 'en');
      const firstRun = await IsFirstRun();
      showSetup = firstRun;
    } catch (e) {
      console.error('Init failed:', e);
      await loadLocale('en');
      showSetup = true;
    }
    ready = true;
  });

  function onSetupComplete() {
    showSetup = false;
  }
</script>

<div class="flex h-screen bg-gray-50 dark:bg-slate-900" class:dark={$darkMode}>
  {#if !ready}
    <div class="flex-1 flex items-center justify-center">
      <div class="animate-spin w-8 h-8 border-4 border-primary-500 border-t-transparent rounded-full"></div>
    </div>
  {:else if showSetup}
    <SetupWizard onComplete={onSetupComplete} />
  {:else}
    <Sidebar />
    <main class="flex-1 overflow-y-auto">
      <div class="p-6 max-w-4xl mx-auto">
        {#if $currentPage === 'home'}
          <Home />
        {:else if $currentPage === 'settings'}
          <Settings />
        {:else if $currentPage === 'print'}
          <PrintNow />
        {:else if $currentPage === 'logs'}
          <Logs />
        {/if}
      </div>
    </main>
  {/if}
</div>
