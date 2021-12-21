<div class="flex flex-col">
    <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
        <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
            <div
                class="shadow overflow-hidden border-b border-gray-200 dark:border-gray-700 sm:rounded-lg"
            >
                <table
                    class="min-w-full divide-y divide-gray-200 dark:divide-gray-700"
                >
                    <thead class="bg-gray-50 dark:bg-gray-800">
                        <slot name="header" />
                    </thead>
                    <slot
                        class="bg-white dark:bg-gray-700 divide-y divide-gray-200 dark:divide-gray-800 dark:text-white"
                        name="body"
                    />
                </table>
            </div>
        </div>
    </div>
</div>
