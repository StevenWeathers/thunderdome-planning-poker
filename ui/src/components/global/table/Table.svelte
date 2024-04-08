<div class="overflow-x-auto">
  <table class="w-full text-sm text-left text-gray-500 dark:text-gray-400">
    <thead
      class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400"
    >
      <slot name="header" />
    </thead>
    <slot
      class="bg-slate-100 dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-600 dark:text-white"
      name="body"
    />
  </table>
</div>
