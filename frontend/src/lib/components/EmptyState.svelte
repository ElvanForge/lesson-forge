<script lang="ts">
    interface EmptyStateProps {
        message: string;
    }
    let { message }: EmptyStateProps = $props();
</script>

<div class="text-center py-20 animate-fade-in">
    <div
        class="inline-flex items-center justify-center w-20 h-20 rounded-3xl bg-slate-50 mb-6 text-slate-300"
    >
        <svg
            class="h-10 w-10"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
        >
            <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="1.5"
                d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
            />
        </svg>
    </div>
    <h3 class="text-lg font-bold text-slate-800">No Generations Yet</h3>
    <p class="mt-2 text-slate-500 max-w-xs mx-auto">{message}</p>
</div>
