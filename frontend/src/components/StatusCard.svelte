<script lang="ts">
  export let title: string = '';
  export let value: string = '';
  export let icon: string = '';
  export let color: 'blue' | 'green' | 'yellow' | 'red' | 'gray' = 'blue';

  const colorMap = {
    blue:   'bg-blue-50 text-blue-600 dark:bg-blue-900/20 dark:text-blue-400',
    green:  'bg-green-50 text-green-600 dark:bg-green-900/20 dark:text-green-400',
    yellow: 'bg-yellow-50 text-yellow-600 dark:bg-yellow-900/20 dark:text-yellow-400',
    red:    'bg-red-50 text-red-600 dark:bg-red-900/20 dark:text-red-400',
    gray:   'bg-gray-50 text-gray-600 dark:bg-gray-800 dark:text-gray-400',
  };
</script>

<div class="bg-white dark:bg-slate-800 rounded-xl p-5 shadow-sm border border-gray-100 dark:border-slate-700 transition-all hover:shadow-md">
  <div class="flex items-start gap-4">
    <div class="w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0 {colorMap[color]}">
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={icon} />
      </svg>
    </div>
    <div class="min-w-0">
      <p class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide">{title}</p>
      <p class="mt-1 text-lg font-semibold text-gray-900 dark:text-gray-100 truncate">{value}</p>
    </div>
  </div>
</div>
