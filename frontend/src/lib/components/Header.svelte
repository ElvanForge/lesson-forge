<script lang="ts">
    interface HeaderProps {
        title: string;
        email: string;
        credits: number;
        isLoggedIn: boolean;
        onSignOut: () => void;
    }
    let { title, email, credits, isLoggedIn, onSignOut }: HeaderProps =
        $props();
</script>

<header class="glass sticky top-0 z-50 border-b border-white/40">
    <div
        class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-20 flex items-center justify-between"
    >
        <a
            href="/"
            class="flex items-center gap-3 hover:opacity-80 transition-opacity"
        >
            <div
                class="w-10 h-10 bg-primary rounded-xl flex items-center justify-center shadow-lg shadow-primary/20"
            >
                <svg
                    class="w-6 h-6 text-white"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M13 10V3L4 14h7v7l9-11h-7z"
                    />
                </svg>
            </div>
            <h1 class="text-2xl font-extrabold tracking-tight text-primary">
                {title}
            </h1>
        </a>
        <div class="flex items-center gap-6">
            {#if isLoggedIn}
                <div class="hidden md:flex flex-col items-end">
                    <span
                        class="text-xs font-bold uppercase tracking-wider text-slate-400"
                        >Account</span
                    >
                    <span class="text-sm font-semibold text-slate-600"
                        >{email}</span
                    >
                </div>
                <div class="h-10 w-px bg-slate-200 hidden md:block"></div>
            {/if}
            <div class="flex items-center gap-3">
                <span
                    class="flex items-center gap-1.5 text-sm font-bold bg-primary/10 text-primary px-4 py-2 rounded-full border border-primary/10"
                >
                    <span class="w-2 h-2 rounded-full bg-primary animate-pulse"
                    ></span>
                    {credits} Credits
                </span>
                {#if isLoggedIn}
                    <button
                        class="text-sm font-bold text-slate-400 hover:text-red-500 transition-all"
                        onclick={onSignOut}
                    >
                        Sign Out
                    </button>
                {:else}
                    <a
                        href="/login"
                        class="text-sm font-bold text-primary hover:text-primary/80 transition-all"
                    >
                        Sign In
                    </a>
                {/if}
            </div>
        </div>
    </div>
</header>
