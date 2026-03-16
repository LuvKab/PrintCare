<script lang="ts">
  import { currentPage, t, darkMode } from '../lib/stores';
  import type { Page } from '../lib/stores';

  const navItems: { page: Page; icon: string; key: string }[] = [
    { page: 'home',     icon: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-4 0a1 1 0 01-1-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 01-1 1', key: 'nav.home' },
    { page: 'settings', icon: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.066 2.573c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.573 1.066c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.066-2.573c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z M15 12a3 3 0 11-6 0 3 3 0 016 0z', key: 'nav.settings' },
    { page: 'print',    icon: 'M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4H7v4a2 2 0 002 2zm0-16h6a2 2 0 012 2v2H7V7a2 2 0 012-2z', key: 'nav.print' },
    { page: 'logs',     icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z', key: 'nav.logs' },
  ];

  function navigate(page: Page) {
    $currentPage = page;
  }

  function toggleDark() {
    $darkMode = !$darkMode;
    if ($darkMode) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }
</script>

<aside class="w-56 h-screen flex flex-col bg-white dark:bg-slate-800 border-r border-gray-200 dark:border-slate-700 shadow-sm">
  <div class="px-5 py-5 flex items-center gap-3 border-b border-gray-100 dark:border-slate-700">
    <div class="w-8 h-8 rounded-lg bg-primary-500 flex items-center justify-center">
      <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4H7v4a2 2 0 002 2zm0-16h6a2 2 0 012 2v2H7V7a2 2 0 012-2z" />
      </svg>
    </div>
    <span class="font-semibold text-gray-800 dark:text-gray-100 text-sm">{$t('app.title')}</span>
  </div>

  <nav class="flex-1 px-3 py-4 space-y-1">
    {#each navItems as item}
      <button
        class="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-150
               {$currentPage === item.page
                 ? 'bg-primary-50 text-primary-700 dark:bg-primary-900/30 dark:text-primary-300'
                 : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-slate-700 dark:hover:text-gray-200'}"
        on:click={() => navigate(item.page)}
      >
        <svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d={item.icon} />
        </svg>
        <span>{$t(item.key)}</span>
      </button>
    {/each}
  </nav>

  <div class="px-3 py-4 border-t border-gray-100 dark:border-slate-700">
    <button
      class="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm text-gray-600 hover:bg-gray-50 dark:text-gray-400 dark:hover:bg-slate-700 transition-colors"
      on:click={toggleDark}
    >
      {#if $darkMode}
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
        </svg>
        <span>{$t('sidebar.lightMode')}</span>
      {:else}
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
        </svg>
        <span>{$t('sidebar.darkMode')}</span>
      {/if}
    </button>
  </div>
</aside>
